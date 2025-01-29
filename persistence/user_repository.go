package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/denizdoganinsider/kpi_project/domain"
)

type IUserRepository interface {
	GetAllUsers() []domain.User
	GetUsersByRole(role string) []domain.User
	AddUser(user domain.User) error
	GetById(id int64) (domain.User, error)
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

	log.Fatalf("Product added with %v", result)
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
		return domain.User{}, errors.New(fmt.Sprintf("While getting user with id %d", userId))
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
