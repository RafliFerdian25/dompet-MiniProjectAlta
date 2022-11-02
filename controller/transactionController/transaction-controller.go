package transactionController

import (
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/transactionService"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	TransactionService transactionService.TransactionService
}

func (tc *TransactionController) CreateTransaction(c echo.Context) error {
	var transaction dto.TransactionDTO
	// Binding request body to struct
	err := c.Bind(&transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// Validate request body
	if err = c.Validate(transaction); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "There is an empty field",
			"error":   err.Error(),
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	transaction.UserID = userId

	// Call service to create account
	err = tc.TransactionService.CreateTransaction(transaction)

	// check if there is an error create transaction
	if err != nil {
		// Check if there is an error authorized
		errAuthorized := strings.Contains(err.Error(), "authorized")
		if errAuthorized {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "you are not authorized",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail create transaction",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
