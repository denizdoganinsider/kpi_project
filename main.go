package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/denizdoganinsider/kpi_project/common/app"
	"github.com/denizdoganinsider/kpi_project/common/mysql"
	"github.com/denizdoganinsider/kpi_project/controller"
	"github.com/denizdoganinsider/kpi_project/controller/response"
	"github.com/denizdoganinsider/kpi_project/persistence"
	"github.com/denizdoganinsider/kpi_project/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func main() {
	ctx := context.Background()
	e := echo.New()

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(1))))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:8080"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")

			if origin != "" && origin != "http://127.0.0.1:8080" && origin != "http://localhost:3000" {
				return c.JSON(403, response.ErrorResponse{
					ErrorDescription: "Forbidden",
				})
			}
			return next(c)
		}
	})

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${method} ${uri} ${status} ${latency_human}\n",
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
				c.Request().Header.Set(echo.HeaderXRequestID, requestID)
			}
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}
	})

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "[${time_rfc3339}] ${id} ${method} ${uri} ${status} ${latency_human}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	// Load config
	configurationManager := app.NewConfigurationManager()

	// Create database connection
	db := mysql.GetConnectionPool(ctx, configurationManager.MySqlConfig)

	// Check if db is nil
	if db == nil {
		log.Fatalf("Error: Database connection is nil")
	}

	// Create tables BEFORE closing DB
	runMigrations()

	userRepository := persistence.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(userService)

	// Register routes
	userController.RegisterRoutes(e)
	authController.RegisterRoutes(e)

	// Graceful shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Start echo server and it runs background
	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Fatalf("Error starting Echo server: %v", err)
		}
	}()

	// Wait for termination signal
	<-sigs

	// Shutdown server
	fmt.Println("Shutting down server...")
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Close database connection gracefully on exit
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	fmt.Println("Database connection closed.")
}

func loadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func runMigrations() {
	loadEnvironmentVariables()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Correct database URL format for MySQL with parseTime=true
	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Run migration command
	cmd := exec.Command("migrate", "-path", "db/migrations", "-database", "mysql://"+databaseURL, "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Migration failed: %v\nOutput: %s", err, string(output))
	}
	fmt.Println("Migrations applied successfully")
}
