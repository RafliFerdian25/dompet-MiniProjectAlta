package transactionController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/transactionService"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	TransactionService transactionService.TransactionService
}

func (tc *TransactionController) GetTransaction(c echo.Context) error {
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

	// call service get transaction
	data, err := tc.TransactionService.GetTransaction(month, userId)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get transaction",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get transaction",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get transaction",
		"transactions": data,
	})
}

func (tc *TransactionController) DeleteTransaction(c echo.Context) error {
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
	err = tc.TransactionService.DeleteTransaction(uint(id), userId)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
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
		"message": "success delete transaction",
	})
}

func (tc *TransactionController) UpdateTransaction(c echo.Context) error {
	var transaction dto.TransactionDTO
	// Binding request body to struct
	err := c.Bind(&transaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// Get id from url
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "fail get id",
			"error":   err.Error(),
		})
	}
	transaction.ID = uint(id)

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// Call service to update account
	err = tc.TransactionService.UpdateTransaction(transaction, userId)

	// check if there is an error update transaction
	if err != nil {
		// Check if there is an client error
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail update transaction",
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail update transaction",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
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
		// Check if there is an error client
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
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
