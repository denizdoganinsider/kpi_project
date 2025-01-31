package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/denizdoganinsider/kpi_project/domain"
)

type IUserRepository interface {
	GetAllUsers() []domain.User
	GetUsersByRole(role string) []domain.User
	AddUser(user domain.User) error
	GetById(id int64) (domain.User, error)
	DeleteById(id int64) error
	UpdateUsername(username string, id int64) error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (userRepository *UserRepository) GetAllUsers() []domain.User {
	ctx := context.Background()
	query := "SELECT * FROM users"

	userRows, err := userRepository.db.QueryContext(ctx, query)

	if err != nil {
		log.Fatalf("Error getting all user's table: %v", err)
		return []domain.User{}
	}

	return extractUserFromRows(userRows)
}

func (userRepository *UserRepository) GetUsersByRole(role string) []domain.User {
	ctx := context.Background()
	query := `SELECT * FROM users WHERE role = ?`

	userRows, err := userRepository.db.QueryContext(ctx, query, role)

	if err != nil {
		log.Fatalf("Error getting all user's table by role: %v", err)
		return []domain.User{}
	}

	return extractUserFromRows(userRows)
}

func (userRepository *UserRepository) AddUser(user domain.User) error {
	insertQuery := `INSERT INTO users (username, email, password_hash, role) VALUES (?, ?, ?, ?)`

	result, err := userRepository.db.Exec(insertQuery, user.Username, user.Email, user.PasswordHash, user.Role)

	if err != nil {
		log.Fatalf("Failed to add new user %v", err)
		return err
	}

	log.Printf("User added with %v", result)
	return nil
}

func (userRepository *UserRepository) GetById(userId int64) (domain.User, error) {
	getByIdSql := `SELECT * FROM users WHERE id = ?`

	queryRow := userRepository.db.QueryRow(getByIdSql, userId)

	var id int64
	var username string
	var email string
	var password_hash string
	var role string
	var created_at string
	var updated_at string

	err := queryRow.Scan(&id, &username, &email, &password_hash, &role, &created_at, &updated_at)

	if err != nil {
		return domain.User{}, fmt.Errorf("error occurred while getting user with id %d", userId)
	}

	return domain.User{
		Id:           id,
		Username:     username,
		Email:        email,
		PasswordHash: password_hash,
		Role:         role,
		CreatedAt:    created_at,
	}, nil
}

func (userRepository *UserRepository) DeleteById(id int64) error {
	deleteSql := `DELETE FROM users WHERE id = ?`

	_, err := userRepository.GetById(id)

	if err != nil {
		return errors.New("user not found")
	}

	_, err = userRepository.db.Exec(deleteSql, id)

	if err != nil {
		return fmt.Errorf("error deleting user with id %d", id)
	}

	log.Println("User deleted successfully")
	return nil
}

func (userRepository *UserRepository) UpdateUsername(username string, id int64) error {
	updateSql := `UPDATE users SET username = ?, updated_at = ? WHERE id = ?`
	updatedAt := time.Now()
	_, err := userRepository.db.Exec(updateSql, username, updatedAt, id)

	if err != nil {
		return fmt.Errorf("error updating user with id %d: %w", id, err)
	}

	log.Println("Username updated successfully")
	return nil
}

func extractUserFromRows(userRows *sql.Rows) []domain.User {
	var users = []domain.User{}
	var id int64
	var username string
	var email string
	var password_hash string
	var roleName string
	var created_at string
	var updated_at string

	for userRows.Next() {
		userRows.Scan(&id, &username, &email, &password_hash, &roleName, &created_at, &updated_at)
		users = append(users, domain.User{
			Id:           id,
			Username:     username,
			Email:        email,
			PasswordHash: password_hash,
			Role:         roleName,
			CreatedAt:    created_at,
			UpdatedAt:    updated_at,
		})
	}

	return users
}
