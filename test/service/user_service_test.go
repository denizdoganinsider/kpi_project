package service

import (
	"os"
	"testing"

	"github.com/denizdoganinsider/kpi_project/domain"
	"github.com/denizdoganinsider/kpi_project/service"
	"github.com/denizdoganinsider/kpi_project/service/model"
	"github.com/stretchr/testify/assert"
)

var userService service.IUserService

func TestMain(m *testing.M) {
	initialUsers := []domain.User{
		{
			Username:     "TestUser1",
			Email:        "TestEmail1",
			PasswordHash: "TestPasswordHash1",
			Role:         "TestRole1",
		},
		{
			Username:     "TestUser2",
			Email:        "TestEmail2",
			PasswordHash: "TestPasswordHash2",
			Role:         "TestRole2",
		},
	}

	fakeUserRepository := NewFakeUserRepository(initialUsers)
	userService = service.NewUserService(fakeUserRepository)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func Test_GetAllUsers(t *testing.T) {
	t.Run("GetAllUsers", func(t *testing.T) {
		actualUsers := userService.GetAllUsers()
		assert.Equal(t, 2, len(actualUsers))
	})
}

func Test_WhenNoValidationErrorOccurred_ShouldAddUser(t *testing.T) {
	t.Run("WhenNoValidationErrorOccurred_ShouldAddUser", func(t *testing.T) {
		userService.AddUser(model.UserCreate{
			Username: "TestUserForAdding",
			Email:    "TestEmailForUserAdding@useinsider.com",
			Password: "TestPasswordHashForAdding",
			Role:     "TestRoleForAdding",
		})

		actualUsers := userService.GetAllUsers()

		assert.Equal(t, 3, len(actualUsers))
	})
}
