package persistence

import (
	"database/sql"
	"errors"

	"github.com/denizdoganinsider/kpi_project/domain"
)

type IAuthRepository interface {
	RegisterUser(user domain.User) error
	GetUserByEmail(email string) (domain.User, error)
}

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) RegisterUser(user domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, role) VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.Role)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(email string) (domain.User, error) {
	query := `SELECT id, username, email, password_hash, role FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	var user domain.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}
