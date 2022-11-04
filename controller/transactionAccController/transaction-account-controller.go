package transactionAccController

import (
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/transactionAccService"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type TransactionAccController struct {
	TransAccService transactionAccService.TransactionAccService
}

// DeleteTransactionAccount
func (tac *TransactionAccController) DeleteTransactionAccount(c echo.Context) error {
	// Get id from url
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "fail get id",
			"error":   err.Error(),
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// call service to delete the transaction
	err = tac.TransAccService.DeleteTransactionAccount(uint(id), userId)
	if err != nil {
		// check if there is an error client
		errClient := strings.Contains(err.Error(), "client :")
		if errClient {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "fail delete transaction",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail delete transaction",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}

// CreateTransactionAccount 
func (tac *TransactionAccController) CreateTransactionAccount(c echo.Context) error {
	var transAcc dto.TransactionAccount
	// Binding request body to struct
	err := c.Bind(&transAcc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// Validate request body
	if err = c.Validate(transAcc); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "There is an empty field",
			"error":   err.Error(),
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	transAcc.UserID = userId

	// Call service to create account
	err = tac.TransAccService.CreateTransactionAccount(transAcc)

	// check if there is an error create transaction
	if err != nil {
		// Check if there is an error client
		errAuthorized := strings.Contains(err.Error(), "authorized")
		errCategory := strings.Contains(err.Error(), "change category")
		errBalance := strings.Contains(err.Error(), "enough balance")
		if errAuthorized || errCategory || errBalance {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "fail create transaction",
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