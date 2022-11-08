package categoryController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/service/categoryService"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	CategoryService categoryService.CategoryService
}

func (cc *CategoryController) GetCategoryByID(c echo.Context) error  {
	// get id from url param
	paramID := c.Param("categoryId")
	// convert string to int
	categoryID, err := strconv.Atoi(paramID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "fail get category by id",
			"error":   err.Error(),
		})
	}

	// get category by id from service
	categoriesID, err := cc.CategoryService.GetCategoryByID(uint(categoryID))
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get category by id",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get category by id",
			"error":   err.Error(),
		})
	}

	// return response success
	return c.JSON(http.StatusOK, echo.Map{
		"message":   "success",
		"categories": categoriesID,
	})
}

func (cc *CategoryController) GetAllCategory(c echo.Context) error {
	categories, err := cc.CategoryService.GetAllCategory()
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get all category",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get all category",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":   "success",
		"categories": categories,
	})
}

