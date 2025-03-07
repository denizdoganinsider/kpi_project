package service

import (
	"errors"
	"log"

	"github.com/denizdoganinsider/kpi_project/domain"
	"github.com/denizdoganinsider/kpi_project/persistence"
)

type IBalanceService interface {
	GetBalanceByUserID(userID int64) (*domain.Balance, error)
	UpdateBalance(userID int64, amount float64) error
	CreateBalance(userID int64, amount float64) error
}

type BalanceService struct {
	balanceRepository persistence.IBalanceRepository
}

func NewBalanceService(balanceRepository persistence.IBalanceRepository) IBalanceService {
	return &BalanceService{
		balanceRepository: balanceRepository,
	}
}

func (s *BalanceService) GetBalanceByUserID(userID int64) (*domain.Balance, error) {
	balance, err := s.balanceRepository.GetBalanceByUserID(userID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (s *BalanceService) UpdateBalance(userID int64, amount float64) error {
	balance, err := s.balanceRepository.GetBalanceByUserID(userID)
	if err != nil {
		return err
	}

	// Eğer kullanıcı için bakiye bulunmazsa, yeni bir bakiye oluşturulur
	if balance == nil {
		err = s.balanceRepository.CreateBalance(userID, amount)
		if err != nil {
			return err
		}
		log.Printf("New balance created for user %d with amount %f", userID, amount)
		return nil
	}

	// Eğer bakiye varsa, yeni bakiyeyi güncelleme işlemi yapılır
	newAmount := balance.Amount + amount
	if newAmount < 0 {
		return errors.New("insufficient balance")
	}

	err = s.balanceRepository.UpdateBalance(userID, newAmount)
	if err != nil {
		return err
	}

	log.Printf("Balance for user %d updated to %f", userID, newAmount)
	return nil
}

func (s *BalanceService) CreateBalance(userID int64, amount float64) error {
	// Kullanıcı için yeni bir bakiye oluşturulur
	err := s.balanceRepository.CreateBalance(userID, amount)
	if err != nil {
		return err
	}

	log.Printf("New balance created for user %d with amount %f", userID, amount)
	return nil
}
