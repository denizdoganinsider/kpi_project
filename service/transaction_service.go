package service

import (
	"errors"
	"time"

	"github.com/denizdoganinsider/kpi_project/domain"
	"github.com/denizdoganinsider/kpi_project/persistence"
)

// ITransactionService, transaction service katmanının sunduğu metodları tanımlar
type ITransactionService interface {
	Credit(userID int64, amount float64) (*domain.Transaction, error)
	Debit(userID int64, amount float64) (*domain.Transaction, error)
	Transfer(fromUserID int64, toUserID int64, amount float64) (*domain.Transaction, error)
	GetTransactionHistory(userID int64) ([]domain.Transaction, error)
	GetTransactionByID(transactionID int64) (*domain.Transaction, error)
}

// TransactionService, transaction işlemlerini yöneten struct'tır
type TransactionService struct {
	transactionRepo persistence.ITransactionRepository
	balanceRepo     persistence.IBalanceRepository // Bakiye işlemleri için balance repository'si (örneğin debit, credit için)
}

// NewTransactionService, yeni bir TransactionService oluşturur
func NewTransactionService(transactionRepo persistence.ITransactionRepository, balanceRepo persistence.IBalanceRepository) ITransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		balanceRepo:     balanceRepo,
	}
}

// Credit, kullanıcının hesabına kredi ekler
func (s *TransactionService) Credit(userID int64, amount float64) (*domain.Transaction, error) {
	// Bakiye kontrolü yapabilirsiniz (örneğin negatif bakiyeleri engellemek için)
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Yeni bir transaction oluşturur
	transaction := &domain.Transaction{
		FromUser:  userID, // Krediyi veren kullanıcı
		Amount:    amount,
		Type:      domain.CreditTransaction,
		Status:    domain.Pending, // İlk olarak pending olarak gelir
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Transaction'ı veritabanına ekler
	err := s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Bakiye güncellenmesi
	err = s.balanceRepo.UpdateBalance(userID, amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Debit, kullanıcının hesabından para çeker
func (s *TransactionService) Debit(userID int64, amount float64) (*domain.Transaction, error) {
	// Bakiye kontrolü yapabilirsiniz (örneğin yeterli bakiye kontrolü)
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Yeni bir transaction oluşturur
	transaction := &domain.Transaction{
		FromUser:  userID,  // Parayı çeken kullanıcı
		Amount:    -amount, // Miktar negatif olacak çünkü çekme işlemi
		Type:      domain.DebitTransaction,
		Status:    domain.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Transaction'ı veritabanına ekler
	err := s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Bakiye güncellenmesi
	err = s.balanceRepo.UpdateBalance(userID, -amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Transfer, bir kullanıcıdan başka bir kullanıcıya para transferi yapar
func (s *TransactionService) Transfer(fromUserID int64, toUserID int64, amount float64) (*domain.Transaction, error) {
	// Yeterli bakiye kontrolü
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if fromUserID == toUserID {
		return nil, errors.New("transfer cannot be to the same user")
	}

	// Transfer işlemi için yeni bir transaction oluşturur
	transaction := &domain.Transaction{
		FromUser:  fromUserID, // Para gönderen kullanıcı
		ToUser:    &toUserID,  // Para alıcı kullanıcı
		Amount:    amount,
		Type:      domain.TransferTransaction,
		Status:    domain.Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Transaction'ı veritabanına ekler
	err := s.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Gönderenin bakiyesi azaltılacak
	err = s.balanceRepo.UpdateBalance(fromUserID, -amount)
	if err != nil {
		return nil, err
	}

	// Alıcının bakiyesi arttırılacak
	err = s.balanceRepo.UpdateBalance(toUserID, amount)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionHistory, kullanıcıya ait tüm işlemleri getirir
func (s *TransactionService) GetTransactionHistory(userID int64) ([]domain.Transaction, error) {
	return s.transactionRepo.GetUserTransactions(userID)
}

// GetTransactionByID, belirtilen ID'ye sahip bir işlemi getirir
func (s *TransactionService) GetTransactionByID(transactionID int64) (*domain.Transaction, error) {
	return s.transactionRepo.GetTransactionByID(transactionID)
}
