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

func (balanceService *BalanceService) GetBalanceByUserID(userID int64) (*domain.Balance, error) {
	balance, err := balanceService.balanceRepository.GetBalanceByUserID(userID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (balanceService *BalanceService) UpdateBalance(userID int64, amount float64) error {
	balance, err := balanceService.balanceRepository.GetBalanceByUserID(userID)
	if err != nil && err.Error() != "user doesn't have balance" {
		return err
	}

	/* if we cannot find any balance for user, system create a new balance */
	if balance == nil {
		err = balanceService.balanceRepository.CreateBalance(userID, amount)
		if err != nil {
			return err
		}
		log.Printf("New balance created for user %d with amount %f", userID, amount)
		return nil
	}

	/* if there is a balance we update to new balance */
	newAmount := balance.Amount + amount
	if newAmount < 0 {
		return errors.New("insufficient balance")
	}

	err = balanceService.balanceRepository.UpdateBalance(userID, newAmount)
	if err != nil {
		return err
	}

	log.Printf("Balance for user %d updated to %f", userID, newAmount)
	return nil
}

func (balanceService *BalanceService) CreateBalance(userID int64, amount float64) error {
	// Kullanıcı için yeni bir bakiye oluşturulur
	err := balanceService.balanceRepository.CreateBalance(userID, amount)
	if err != nil {
		return err
	}

	log.Printf("New balance created for user %d with amount %f", userID, amount)
	return nil
}
