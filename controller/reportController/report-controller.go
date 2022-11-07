package reportController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/service/reportService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReportController struct {
	ReportService reportService.ReportService
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
		"data": dataReport,
	})
}