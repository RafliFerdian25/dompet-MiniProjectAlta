package subCategoryController

import (
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/subCategoryService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SubCategoryController struct {
	SubCategoryService subCategoryService.SubCategoryService
}

func (sc *SubCategoryController) CreateSubCategory(c echo.Context) error {
	var subCategory dto.SubCategoryDTO
	// Binding request body to struct
	err := c.Bind(&subCategory)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	subCategory.UserID = userId

	// Validate request body
	if err = c.Validate(subCategory); err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
    }

	// Call service to create sub category
	err = sc.SubCategoryService.CreateSubCategory(subCategory)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail create sub category",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
