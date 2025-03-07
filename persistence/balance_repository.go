package persistence

import (
	"database/sql"
	"fmt"

	"github.com/denizdoganinsider/kpi_project/domain"
)

type IBalanceRepository interface {
	GetBalanceByUserID(userID int64) (*domain.Balance, error)
	UpdateBalance(userID int64, amount float64) error
	CreateBalance(userID int64, amount float64) error
}

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) IBalanceRepository {
	return &BalanceRepository{db: db}
}

func (balanceRepository *BalanceRepository) GetBalanceByUserID(userID int64) (*domain.Balance, error) {
	var count int
	doesUserExistsQuery := `SELECT COUNT(*) FROM users WHERE id = ?`
	errForUserExistence := balanceRepository.db.QueryRow(doesUserExistsQuery, userID).Scan(&count)
	if errForUserExistence != nil {
		return nil, errForUserExistence
	} else if count == 0 {
		return nil, fmt.Errorf("user not found with given id %d", userID)
	}

	query := `SELECT user_id, amount, last_updated_at FROM balances WHERE user_id = ?`
	row := balanceRepository.db.QueryRow(query, userID)

	var balance domain.Balance
	err := row.Scan(&balance.UserID, &balance.Amount, &balance.LastUpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user doesn't have balance with given id %d", userID)
		}
		return nil, err
	}

	return &balance, nil
}

func (balanceRepository *BalanceRepository) UpdateBalance(userID int64, amount float64) error {
	query := `UPDATE balances SET amount = ?, last_updated_at = NOW() WHERE user_id = ?`
	_, err := balanceRepository.db.Exec(query, amount, userID)
	return err
}

func (balanceRepository *BalanceRepository) CreateBalance(userID int64, amount float64) error {
	query := `INSERT INTO balances (user_id, amount, last_updated_at) VALUES (?, ?, NOW())`
	_, err := balanceRepository.db.Exec(query, userID, amount)
	return err
}
