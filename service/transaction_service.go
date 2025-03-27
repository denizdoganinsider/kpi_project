package service

import (
	"errors"
	"time"

	"github.com/denizdoganinsider/kpi_project/domain"
	"github.com/denizdoganinsider/kpi_project/persistence"
)

type ITransactionService interface {
	Credit(userID int64, amount float64) (*domain.Transaction, error)
	Debit(userID int64, amount float64) (*domain.Transaction, error)
	Transfer(fromUserID int64, toUserID int64, amount float64) (*domain.Transaction, error)
	GetTransactionHistory(userID int64) ([]domain.Transaction, error)
	GetTransactionByID(transactionID int64) (*domain.Transaction, error)
}

type TransactionService struct {
	transactionRepository persistence.ITransactionRepository
	balanceRepo           persistence.IBalanceRepository
}

func NewTransactionService(transactionRepository persistence.ITransactionRepository, balanceRepo persistence.IBalanceRepository) ITransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		balanceRepo:           balanceRepo,
	}
}

func (transactionService *TransactionService) Credit(userID int64, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	transaction := &domain.Transaction{
		FromUser:  userID,
		Amount:    amount,
		Type:      domain.CreditTransaction,
		Status:    domain.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := transactionService.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	err = transactionService.balanceRepo.UpdateBalance(userID, amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (transactionService *TransactionService) Debit(userID int64, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	transaction := &domain.Transaction{
		FromUser:  userID,
		Amount:    -amount,
		Type:      domain.DebitTransaction,
		Status:    domain.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := transactionService.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	err = transactionService.balanceRepo.UpdateBalance(userID, -amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (transactionService *TransactionService) Transfer(fromUserID int64, toUserID int64, amount float64) (*domain.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if fromUserID == toUserID {
		return nil, errors.New("transfer cannot be to the same user")
	}

	transaction := &domain.Transaction{
		FromUser:  fromUserID,
		ToUser:    &toUserID,
		Amount:    amount,
		Type:      domain.TransferTransaction,
		Status:    domain.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := transactionService.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	err = transactionService.balanceRepo.UpdateBalance(fromUserID, -amount)
	if err != nil {
		return nil, err
	}

	err = transactionService.balanceRepo.UpdateBalance(toUserID, amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (transactionService *TransactionService) GetTransactionHistory(userID int64) ([]domain.Transaction, error) {
	return transactionService.transactionRepository.GetUserTransactions(userID)
}

func (transactionService *TransactionService) GetTransactionByID(transactionID int64) (*domain.Transaction, error) {
	return transactionService.transactionRepository.GetTransactionByID(transactionID)
}
