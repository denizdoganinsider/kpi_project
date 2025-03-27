package persistence

import (
	"database/sql"

	"github.com/denizdoganinsider/kpi_project/domain"
)

type ITransactionRepository interface {
	CreateTransaction(transaction *domain.Transaction) error
	GetTransactionByID(id int64) (*domain.Transaction, error)
	UpdateTransactionStatus(id int64, status domain.TransactionStatus) error
	GetUserTransactions(userID int64) ([]domain.Transaction, error)
	UpdateBalance(userID int64, amount float64) error
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) ITransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(transaction *domain.Transaction) error {
	var toUser sql.NullInt64
	if transaction.ToUser != nil {
		toUser.Int64 = *transaction.ToUser
		toUser.Valid = true
	}

	query := `INSERT INTO transactions (from_user_id, to_user_id, amount, type, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := repo.db.Exec(query, transaction.FromUser, transaction.ToUser, transaction.Amount, transaction.Type, transaction.Status, transaction.CreatedAt)
	return err
}

func (repo *TransactionRepository) GetTransactionByID(id int64) (*domain.Transaction, error) {
	query := `SELECT * FROM transactions WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var transaction domain.Transaction
	var toUser sql.NullInt64
	err := row.Scan(&transaction.ID, &transaction.FromUser, &toUser, &transaction.Amount, &transaction.Type, &transaction.Status, &transaction.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if toUser.Valid {
		transaction.ToUser = &toUser.Int64
	}

	return &transaction, nil
}

func (repo *TransactionRepository) UpdateTransactionStatus(id int64, status domain.TransactionStatus) error {
	query := `UPDATE transactions SET status = ? WHERE id = ?`
	_, err := repo.db.Exec(query, status, id)
	return err
}

func (repo *TransactionRepository) GetUserTransactions(userID int64) ([]domain.Transaction, error) {
	query := `SELECT id, from_user_id, to_user_id, amount, type, status, created_at FROM transactions WHERE from_user_id = ? OR to_user_id = ?`
	rows, err := repo.db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction

	for rows.Next() {
		var transaction domain.Transaction
		var toUser sql.NullInt64
		err := rows.Scan(&transaction.ID, &transaction.FromUser, &toUser, &transaction.Amount, &transaction.Type, &transaction.Status, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}

		if toUser.Valid {
			transaction.ToUser = &toUser.Int64
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (repo *TransactionRepository) UpdateBalance(userID int64, amount float64) error {
	query := `UPDATE balances SET amount = amount + ? WHERE user_id = ?`
	_, err := repo.db.Exec(query, amount, userID)
	return err
}
