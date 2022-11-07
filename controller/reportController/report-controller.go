package reportController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/service/reportService"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ReportController struct {
	ReportService reportService.ReportService
}

func (ac *ReportController) GetCashflow(c echo.Context) error {
	// get query parameters
	paramPeriod := c.QueryParam("period")
	if paramPeriod != "month" && paramPeriod != "week" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid period",
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// call service to get cashflow
	cashflow, err := ac.ReportService.GetCashflow(userId, paramPeriod)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get cashflow",
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get cashflow",
			"error": err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get cashflow",
		"data":    cashflow,
	})
}


func (ac *ReportController) GetReportbyCategory(c echo.Context) error {
	// get query parameters
	paramPeriod := c.QueryParam("period")
	if paramPeriod != "month" && paramPeriod != "week" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid period",
		})
	}

	paramNumberPeriod := c.QueryParam("number_period")
	var numberPeriod int
	if paramNumberPeriod == "" {
		if paramPeriod == "month" {
			numberPeriod = int(time.Now().Month())
		} else if paramPeriod == "week" {
			_, numberPeriod = time.Now().ISOWeek()
		}
	} else {
		var err error
		numberPeriod, err = strconv.Atoi(paramNumberPeriod)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid number_period",
				"error":   "number_period must be integer",
			})
		}
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// call service method to get report
	dataReport, err := ac.ReportService.GetReportbyCategory(userId, paramPeriod, numberPeriod)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get report by category",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get report by category",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get report by category",
		"data":    dataReport,
	})
}

func (ac *ReportController) GetAnalyticPeriod(c echo.Context) error {
	// get query parameters
	paramPeriod := c.QueryParam("period")
	if paramPeriod != "month" && paramPeriod != "week" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid period",
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// call service method to get report
	dataReport, err := ac.ReportService.GetAnalyticPeriod(userId, paramPeriod)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get analytic period",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get analytic period",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get report " + paramPeriod,
		"data":    dataReport,
	})
}
