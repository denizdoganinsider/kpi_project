package service

import (
	"github.com/denizdoganinsider/kpi_project/domain"
	"github.com/denizdoganinsider/kpi_project/persistence"
)

type FakeUserRepository struct {
	users []domain.User
}

func NewFakeUserRepository(initialUsers []domain.User) persistence.IUserRepository {
	return &FakeUserRepository{
		users: initialUsers,
	}
}

func (fakeUserRepository *FakeUserRepository) GetAllUsers() []domain.User {
	return fakeUserRepository.users
}

func (fakeUserRepository *FakeUserRepository) GetUsersByRole(role string) []domain.User {
	/* Not implemented yet */
	return []domain.User{}
}

func (fakeUserRepository *FakeUserRepository) AddUser(user domain.User) error {

	fakeUserRepository.users = append(fakeUserRepository.users, domain.User{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	return nil
}

func (fakeUserRepository *FakeUserRepository) GetById(userId int64) (domain.User, error) {
	/* Not implemented yet */
	return domain.User{}, nil
}

func (fakeUserRepository *FakeUserRepository) DeleteById(id int64) error {
	/* Not implemented yet */
	return nil
}

func (fakeUserRepository *FakeUserRepository) UpdateUsername(username string, id int64) error {
	/* Not implemented yet */
	return nil
}
