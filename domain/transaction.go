package domain

import (
	"errors"
	"time"
)

type TransactionType string

const (
	CreditTransaction   TransactionType = "credit"
	DebitTransaction    TransactionType = "debit"
	TransferTransaction TransactionType = "transfer"
)

type TransactionStatus string

const (
	Pending   TransactionStatus = "pending"
	Completed TransactionStatus = "completed"
	Failed    TransactionStatus = "failed"
)

type Transaction struct {
	ID        int64
	FromUser  int64
	ToUser    *int64
	Amount    float64
	Type      TransactionType
	Status    TransactionStatus
	CreatedAt time.Time
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if t.Type == TransferTransaction && (t.ToUser == nil || *t.ToUser == t.FromUser) {
		return errors.New("invalid transfer transaction")
	}
	return nil
}
