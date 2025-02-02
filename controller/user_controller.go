package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/denizdoganinsider/kpi_project/controller/request"
	"github.com/denizdoganinsider/kpi_project/controller/response"
	"github.com/denizdoganinsider/kpi_project/service"
	"github.com/labstack/echo/v4"
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

	e.GET("/", userController.HandleRootPage)

	e.POST("/api/v1/submit", userController.SubmitNewUser)

	e.GET("/api/v1/users", userController.GetAllUsers)
	e.GET("/api/v1/users/:id", userController.GetUserById)
	e.POST("/api/v1/users", userController.AddUser)
	e.PUT("/api/v1/users/:id", userController.UpdateUsername)
	e.DELETE("/api/v1/users/:id", userController.DeleteUserById)
}

func (userController *UserController) HandleRootPage(c echo.Context) error {
	return c.HTML(http.StatusOK,
		`<form action="/api/v1/submit" method="POST">
			<label for="name">First name:</label>
			<input type="text" id="name" name="name"><br><br>
			<label for="surname">Last name:</label>
			<input type="text" id="surname" name="surname"><br><br>
			<input type="submit" value="Submit">
		</form>`)
}

func (userController *UserController) SubmitNewUser(c echo.Context) error {
	firstName := c.FormValue("name")
	lastName := c.FormValue("surname")

	fmt.Println("Received Form Data:")
	fmt.Println("First Name:", firstName)
	fmt.Println("Last Name:", lastName)

	return c.String(http.StatusOK, "Form submitted successfully!")
}

func (userController *UserController) GetAllUsers(c echo.Context) error {
	role := c.QueryParam("role")

	if len(role) == 0 {
		allUsers := userController.userService.GetAllUsers()
		return c.JSON(http.StatusOK, response.ToResponseList(allUsers))
	}

	usersWithGivenRole := userController.userService.GetUsersByRole(role)
	return c.JSON(http.StatusOK, response.ToResponseList(usersWithGivenRole))
}

func (userController *UserController) GetUserById(c echo.Context) error {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)

	userContent, err := userController.userService.GetById(int64(userId))

	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.ToResponse(userContent))
}

func (userController *UserController) AddUser(c echo.Context) error {
	var addUserRequest request.AddUserRequest

	bindError := c.Bind(&addUserRequest)

	if bindError != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: bindError.Error(),
		})
	}

	validationError := userController.userService.AddUser(addUserRequest.ToModel())

	if validationError != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorDescription: validationError.Error(),
		})
	}

	return c.NoContent(http.StatusCreated)
}

func (userController *UserController) UpdateUsername(c echo.Context) error {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)

	username := c.QueryParam("username")

	if len(username) == 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "Username is required",
		})
	}

	if len(username) < 4 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "Username should have at least 4 characters",
		})
	}

	userController.userService.UpdateUsername(username, int64(userId))

	return c.NoContent(http.StatusOK)
}

func (userController *UserController) DeleteUserById(c echo.Context) error {
	id := c.Param("id")
	userId, _ := strconv.Atoi(id)

	err := userController.userService.DeleteById(int64(userId))

	if err != nil {
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
