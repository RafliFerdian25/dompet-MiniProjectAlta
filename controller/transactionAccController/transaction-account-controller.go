package transactionAccController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/transactionAccService"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type TransactionAccController struct {
	TransAccService transactionAccService.TransactionAccService
}

// GetTransactionAccount 
func (tac *TransactionAccController) GetTransactionAccount(c echo.Context) error {
	// Get month from url
	paramMonth := c.QueryParam("month")
	var month int
	if paramMonth == "" {
		month = int(time.Now().Month())
	} else {
		var err error
		month, err = strconv.Atoi(paramMonth)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "fail get month",
				"error":   err.Error(),
			})
		}
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// Call service to get transaction
	transAcc, err := tac.TransAccService.GetTransactionAccount(userId, month)
	if err != nil {
		// Check if there is an error client
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get transaction account",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get transaction account",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get transaction account",
		"transaction_account": map[string]interface{}{
			"month": map[string]interface{}{
				strconv.Itoa(month): transAcc,
			},
		},
	})
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
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail delete transaction account",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail delete transaction account",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success delete transaction account",
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
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail create transaction account",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail create transaction account",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success create transaction account",
	})
}