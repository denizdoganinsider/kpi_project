package persistence

import (
	"database/sql"

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

func (repo *BalanceRepository) GetBalanceByUserID(userID int64) (*domain.Balance, error) {
	query := `SELECT user_id, amount, last_updated_at FROM balances WHERE user_id = ?`
	row := repo.db.QueryRow(query, userID)

	var balance domain.Balance
	err := row.Scan(&balance.UserID, &balance.Amount, &balance.LastUpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &balance, nil
}

func (repo *BalanceRepository) UpdateBalance(userID int64, amount float64) error {
	query := `UPDATE balances SET amount = ?, last_updated_at = NOW() WHERE user_id = ?`
	_, err := repo.db.Exec(query, amount, userID)
	return err
}

func (repo *BalanceRepository) CreateBalance(userID int64, amount float64) error {
	query := `INSERT INTO balances (user_id, amount, last_updated_at) VALUES (?, ?, NOW())`
	_, err := repo.db.Exec(query, userID, amount)
	return err
}
