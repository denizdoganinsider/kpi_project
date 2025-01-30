package service

import (
	"fmt"
	"strings"

	"github.com/denizdoganinsider/kpi_project/domain"
	"github.com/denizdoganinsider/kpi_project/persistence"
	"github.com/denizdoganinsider/kpi_project/service/common"
	"github.com/denizdoganinsider/kpi_project/service/model"
)

type IUserService interface {
	Add(UserCreate model.UserCreate) error
	DeleteById(id int64) error
	GetById(id int64) (domain.User, error)
	UpdateUsername(username string, id int64) error
	GetAllUsers() []domain.User
	GetUsersByRole(role string) []domain.User
}

type UserService struct {
	userRepository persistence.IUserRepository
}

func NewUserService(userRepository persistence.IUserRepository) IUserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (userService *UserService) Add(userCreate model.UserCreate) error {
	validateError := validateProductCreate(userCreate)
	if validateError != nil {
		return validateError
	}

	return userService.userRepository.AddUser(domain.User{
		Username:     userCreate.Username,
		Email:        userCreate.Email,
		PasswordHash: userCreate.PasswordHash,
		Role:         userCreate.Role,
	})
}

func (userService *UserService) DeleteById(id int64) error {
	return userService.userRepository.DeleteById(id)
}

func (userService *UserService) GetById(id int64) (domain.User, error) {
	return userService.userRepository.GetById(id)
}

func (userService *UserService) UpdateUsername(username string, id int64) error {
	return userService.userRepository.UpdateUsername(username, id)
}

func (userService *UserService) GetAllUsers() []domain.User {
	return userService.userRepository.GetAllUsers()
}

func (userService *UserService) GetUsersByRole(role string) []domain.User {
	return userService.userRepository.GetUsersByRole(role)
}

func validateProductCreate(userCreate model.UserCreate) error {
	if !strings.Contains(userCreate.Email, common.AT_SYMBOl) {
		return fmt.Errorf("the given email doesn't contains %s", common.AT_SYMBOl)
	}
	return nil
}
