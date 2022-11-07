package debtController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/debtService"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type DebtController struct {
	DebtService debtService.DebtService
}

// delete debt controller
func (dc *DebtController) DeleteDebt(c echo.Context) error {
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

	// call service to delete the debt
	err = dc.DebtService.DeleteDebt(uint(id), userId)
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
		"message": "success",
	})

}


func (dc *DebtController) CreateDebt(c echo.Context) error {
	var debtTrans dto.DebtTransactionDTO
	// Binding request body to struct
	err := c.Bind(&debtTrans)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// Validate request body
	if err = c.Validate(debtTrans); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "There is an empty field",
			"error":   err.Error(),
		})
	}

	if debtTrans.SubCategoryID == 2 || debtTrans.SubCategoryID == 4 {
		if debtTrans.DebtID == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "fail create transaction",
				"error":   "debt id is required",
			})
		}
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	debtTrans.UserID = userId

	// Call service to create account
	err = dc.DebtService.CreateDebt(debtTrans)

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