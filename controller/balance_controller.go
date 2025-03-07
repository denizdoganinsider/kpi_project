package controller

import (
	"net/http"
	"strconv"

	"github.com/denizdoganinsider/kpi_project/controller/request"
	"github.com/denizdoganinsider/kpi_project/controller/response"
	"github.com/denizdoganinsider/kpi_project/service"
	"github.com/labstack/echo/v4"
)

type BalanceController struct {
	balanceService service.IBalanceService
}

func NewBalanceController(balanceService service.IBalanceService) *BalanceController {
	return &BalanceController{
		balanceService: balanceService,
	}
}

func (balanceController *BalanceController) RegisterRoutes(e *echo.Echo) {
	// Balance routes
	e.GET("/api/v1/balance/:userID", balanceController.GetBalanceByUserID)
	e.POST("/api/v1/balance/credit", balanceController.CreditBalance)
	e.POST("/api/v1/balance/debit", balanceController.DebitBalance)
}

func (balanceController *BalanceController) GetBalanceByUserID(c echo.Context) error {
	userID := c.Param("userID")
	userId, err := strconv.Atoi(userID)
	if err != nil {
		/* if data format is incorrect */
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "Invalid user ID",
		})
	}

	balance, err := balanceController.balanceService.GetBalanceByUserID(int64(userId))
	if err != nil {
		/* if user doesn't have any balance */
		return c.JSON(http.StatusNotFound, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}

	/* user have a balance so we are returning the balance amount */
	return c.JSON(http.StatusOK, response.GetBalanceResponse{
		Balance: balance.Amount,
	})
}

func (balanceController *BalanceController) CreditBalance(c echo.Context) error {
	return balanceController.UpdateBalance(c, true)
}

func (balanceController *BalanceController) DebitBalance(c echo.Context) error {
	return balanceController.UpdateBalance(c, false)
}

func (balanceController *BalanceController) UpdateBalance(c echo.Context, isCredit bool) error {
	var request request.UpdateBalanceRequest
	var amount float64

	bindError := c.Bind(&request)

	if bindError != nil {
		/* if data format is incorrect */
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "Invalid request data",
		})
	}

	if isCredit {
		amount = request.Amount
	} else {
		amount = -request.Amount
	}
	/* adding the given amount to user's previous balance */
	err := balanceController.balanceService.UpdateBalance(request.UserID, amount)
	if err != nil {
		errDescription := err.Error()

		if errDescription == "user not found" {
			return c.JSON(http.StatusNotFound, response.ErrorResponse{
				ErrorDescription: errDescription,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorDescription: errDescription,
		})
	}

	/* fetching updated balance and sending to user */
	updatedBalance, _ := balanceController.balanceService.GetBalanceByUserID(request.UserID)

	return c.JSON(http.StatusOK, response.GetBalanceResponse{
		Balance: updatedBalance.Amount,
	})
}
