package accountController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/accountService"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AccountController struct {
	AccountService accountService.AccountService
}

// DeleteAccount is a function to delete account
func (ac *AccountController) DeleteAccount(c echo.Context) error {
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

	// Call service to delete account
	err = ac.AccountService.DeleteAccount(uint(id), userId)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail delete account",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail delete account",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success delete account",
	})
}

// GetAccountByUser is a function to get account by user
func (ac *AccountController) GetAccountByUser(c echo.Context) error {
	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// get account by user from service
	userAccounts, err := ac.AccountService.GetAccountByUser(userId)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get account by user",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get account by user",
			"error":   err.Error(),
		})
	}

	// return response success
	return c.JSON(http.StatusOK, echo.Map{
		"message":     "success get account by user",
		"Accounts": userAccounts,
	})
}

// UpdateAccount is a function to update account
func (ac *AccountController) UpdateAccount(c echo.Context) error {
	// Get id from url
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "fail get id",
			"error":   err.Error(),
		})
	}

	// Binding request body to struct
	var account dto.AccountDTO
	err = c.Bind(&account)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// add id to account struct
	account.ID = uint(id)
	
	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	account.UserID = userId

	// Call service to update account
	err = ac.AccountService.UpdateAccount(account)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail update account",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail update account",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update account",
	})
}

// CreateAccount is a function to create account
func (ac *AccountController) CreateAccount(c echo.Context) error {
	var account dto.AccountDTO
	// Binding request body to struct
	err := c.Bind(&account)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// Validate request body
	if err = c.Validate(account); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "There is an empty field",
			"error":   err.Error(),
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	account.UserID = userId

	// Call service to create account
	err = ac.AccountService.CreateAccount(account)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail create account",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail create account",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success create account",
	})
}