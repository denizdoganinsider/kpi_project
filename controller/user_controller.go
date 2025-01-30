package controller

import (
	"net/http"

	"github.com/denizdoganinsider/kpi_project/service"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (userController *UserController) RegisterRoutes(e *echo.Echo) {
	/*
		Should be implemented

		GET /api/v1/users
		GET /api/v1/users/{id}
		PUT /api/v1/users/{id}
		DELETE /api/v1/users/{id}

		Extras

		POST /api/v1/users
	*/

	e.GET("/api/v1/users", userController.GetAllUsers)
	e.GET("api/v1/users/:id", userController.GetUserById)
	e.POST("/api/v1/users", userController.AddUser)
	e.PUT("/api/v1/users/:id", userController.UpdateUsername)
	e.DELETE("/api/v1/users/{id}", userController.DeleteUserById)
}

func (userController *UserController) GetAllUsers(c echo.Context) error {
	role := c.QueryParam("role")

	if len(role) == 0 {
		return c.HTML(http.StatusOK, "<h1>Get All Users</h1>")
	}

	return c.HTML(http.StatusOK, "<h1>Get Users by Role</h1>")
}

func (userController *UserController) GetUserById(c echo.Context) error {
	return nil
}

func (userController *UserController) AddUser(c echo.Context) error {
	return nil
}

func (userController *UserController) UpdateUsername(c echo.Context) error {
	return nil
}

func (userController *UserController) DeleteUserById(c echo.Context) error {
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
