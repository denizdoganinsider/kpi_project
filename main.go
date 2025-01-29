package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/denizdoganinsider/kpi_project/persistence"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Global database instance
var db *sql.DB

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file wasn't loaded.: %v", err)
	}

	// Get connection's information from .env file
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("MySQL connection completed successfully!")

	// Create Tables
	createTables()

	/*
		hashedPassword, err := hashPassword("password123")

		if err != nil {
			log.Fatal("While hashing password, an error occurred:", err)
		}

		// user adding to db
		userID, err := insertUser(db, "deniz", "deniz.dogan@useinsider.com", hashedPassword, "admin")
		if err != nil {
			log.Fatal("While adding user, an error occurred:", err)
		}

		fmt.Println("New user was added. User ID:", userID) */

	/* startServer() */

	var userRepository persistence.IUserRepository = persistence.NewUserRepository(db)

	/* GetAllUsers Test */
	/* fmt.Println(userRepository.GetAllUsers()) */

	/* GetUsersByRole Test */
	/* fmt.Println(userRepository.GetUsersByRole("admin")) */

	/* AddUser Test  */
	/*
		hashedPassword, err := hashPassword("testPassword")
		if err != nil {
			log.Fatal("While hashing password, an error occurred:", err)
		}

		newUser := domain.User{
			Username:     "test",
			Email:        "test@useinsider.com",
			PasswordHash: hashedPassword,
			Role:         "normal-user",
		}

		userRepository.AddUser(newUser)
	*/

	user, _ := userRepository.GetById(2)

	fmt.Println(user)
}

/*
	func insertUser(db *sql.DB, username, email, passwordHash, role string) (int64, error) {
		query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES (?, ?, ?, ?)
		`
		result, err := db.Exec(query, username, email, passwordHash, role)
		if err != nil {
			return 0, err
		}

		// getting added user's id
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}

		return lastInsertID, nil
	}
*/

/*
	func startServer() {
		e := echo.New()

		e.Start("localhost:8080")
	}
*/

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func createTables() {
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

	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating Users table: %v", err)
	}
	fmt.Println("Users table created successfully!")

	_, err = db.Exec(createTransactionsTable)
	if err != nil {
		log.Fatalf("Error creating Transactions table: %v", err)
	}
	fmt.Println("Transactions table created successfully!")

	_, err = db.Exec(createBalancesTable)
	if err != nil {
		log.Fatalf("Error creating Balances table: %v", err)
	}
	fmt.Println("Balances table created successfully!")

	_, err = db.Exec(createAuditLogsTable)
	if err != nil {
		log.Fatalf("Error creating Audit_Logs table: %v", err)
	}
	fmt.Println("Audit_Logs table created successfully!")
}
