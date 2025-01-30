package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/denizdoganinsider/kpi_project/common/app"
	"github.com/denizdoganinsider/kpi_project/common/mysql"
	"github.com/denizdoganinsider/kpi_project/controller"
	"github.com/denizdoganinsider/kpi_project/persistence"
	"github.com/denizdoganinsider/kpi_project/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()
	e := echo.New()

	// Load config
	configurationManager := app.NewConfigurationManager()

	// Create database connection
	db := mysql.GetConnectionPool(ctx, configurationManager.MySqlConfig)

	// Check if db is nil
	if db == nil {
		log.Fatalf("Error: Database connection is nil")
	}

	// Create tables BEFORE closing DB
	createTables(db)

	userRepository := persistence.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	// Register routes
	userController.RegisterRoutes(e)

	// Graceful shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Echo sunucusunu başlat (Arka planda çalışsın)
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

func createTables(db *sql.DB) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	createTransactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		from_user_id INT NOT NULL,
		to_user_id INT NOT NULL,
		amount DECIMAL(10, 2) NOT NULL,
		type VARCHAR(50) NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	createBalancesTable := `
	CREATE TABLE IF NOT EXISTS balances (
		user_id INT PRIMARY KEY,
		amount DECIMAL(10, 2) NOT NULL,
		last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	createAuditLogsTable := `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		entity_type VARCHAR(50) NOT NULL,
		entity_id INT NOT NULL,
		action VARCHAR(50) NOT NULL,
		details TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []struct {
		name string
		sql  string
	}{
		{"Users", createUsersTable},
		{"Transactions", createTransactionsTable},
		{"Balances", createBalancesTable},
		{"AuditLogs", createAuditLogsTable},
	}

	for _, table := range tables {
		_, err := db.Exec(table.sql)
		if err != nil {
			log.Fatalf("Error creating %s table: %v", table.name, err)
		}
		fmt.Printf("%s table created successfully!\n", table.name)
	}
}
