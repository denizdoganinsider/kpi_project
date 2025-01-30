package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

func GetConnectionPool(ctx context.Context, config Config) *sql.DB {
	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.UserName, config.Password, config.Host, config.Port, config.DbName)

	// Connect to database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.MaxConnections)
	db.SetConnMaxIdleTime(time.Duration(config.MaxConnectionIdleTime) * time.Second)

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("MySQL connection established successfully!")

	return db
}
