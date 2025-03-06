package controller

import (
	"fmt"
	"net/http"
	"strconv"

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
	e.GET("/api/v1/balance/:userID", balanceController.GetBalance)
	e.POST("/api/v1/balance/credit", balanceController.CreditBalance)
	e.POST("/api/v1/balance/debit", balanceController.DebitBalance)
}

func (balanceController *BalanceController) GetBalance(c echo.Context) error {
	userID := c.Param("userID")
	userId, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println("UserId:", userID)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	balance, err := balanceController.balanceService.GetBalanceByUserID(int64(userId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]float64{
		"balance": balance.Amount,
	})
}

func (balanceController *BalanceController) CreditBalance(c echo.Context) error {
	var request struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Credit balance (add amount)
	err := balanceController.balanceService.UpdateBalance(request.UserID, request.Amount)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	// Fetch updated balance
	updatedBalance, _ := balanceController.balanceService.GetBalanceByUserID(request.UserID)

	return c.JSON(http.StatusOK, map[string]float64{
		"balance": updatedBalance.Amount,
	})
}

func (balanceController *BalanceController) DebitBalance(c echo.Context) error {
	var request struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	// Debit balance (subtract amount)
	err := balanceController.balanceService.UpdateBalance(request.UserID, -request.Amount)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	// Fetch updated balance
	updatedBalance, _ := balanceController.balanceService.GetBalanceByUserID(request.UserID)

	return c.JSON(http.StatusOK, map[string]float64{
		"balance": updatedBalance.Amount,
	})
}
