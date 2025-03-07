package controller

import (
	"net/http"
	"strconv"

	"github.com/denizdoganinsider/kpi_project/controller/request"
	"github.com/denizdoganinsider/kpi_project/controller/response"
	"github.com/denizdoganinsider/kpi_project/service"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionService service.ITransactionService
}

func NewTransactionController(transactionService service.ITransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

func (transactionController *TransactionController) RegisterRoutes(e *echo.Echo) {
	// Transaction routes
	e.GET("/api/v1/transactions/:id", transactionController.GetTransactionByID)
	e.GET("/api/v1/transactions/history/:userID", transactionController.GetTransactionHistory)
	e.POST("/api/v1/transactions/credit", transactionController.Credit)
	e.POST("/api/v1/transactions/debit", transactionController.Debit)
	e.POST("/api/v1/transactions/transfer", transactionController.Transfer)
}

func (transactionController *TransactionController) GetTransactionByID(c echo.Context) error {
	id := c.Param("id")
	transactionID, _ := strconv.Atoi(id)

	transaction, err := transactionController.transactionService.GetTransactionByID(int64(transactionID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, transaction)
}

func (transactionController *TransactionController) GetTransactionHistory(c echo.Context) error {
	userID := c.Param("userID")
	userId, _ := strconv.Atoi(userID)

	transactions, err := transactionController.transactionService.GetTransactionHistory(int64(userId))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, transactions)
}

func (transactionController *TransactionController) Credit(c echo.Context) error {
	var request struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	transaction, err := transactionController.transactionService.Credit(request.UserID, request.Amount)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, transaction)
}

func (transactionController *TransactionController) Debit(c echo.Context) error {
	var request struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request data",
		})
	}

	transaction, err := transactionController.transactionService.Debit(request.UserID, request.Amount)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, transaction)
}

func (transactionController *TransactionController) Transfer(c echo.Context) error {
	var request request.TransactionRequest

	err := c.Bind(&request)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			ErrorDescription: "Invalid request data",
		})
	}

	transaction, err := transactionController.transactionService.Transfer(request.FromUserID, request.ToUserID, request.Amount)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, response.ErrorResponse{
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, transaction)
}
