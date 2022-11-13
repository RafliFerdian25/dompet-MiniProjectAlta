package subCategoryController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/service/subCategoryService"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SubCategoryController struct {
	SubCategoryService subCategoryService.SubCategoryService
}

// CreateSubCategory is a function to create sub category
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
	
	// Validate request body
	if err = c.Validate(subCategory); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "There is an empty field",
			"error":   err.Error(),
		})
	}

	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	subCategory.UserID = &userId

	// Call service to create sub category
	err = sc.SubCategoryService.CreateSubCategory(subCategory)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail create sub category",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail create sub category",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success create sub category",
	})
}

// GetSubCategoryByUser is a function to get sub category by user
func (sc *SubCategoryController) GetSubCategoryByUser(c echo.Context) error {
	// Get user id from jwt
	userId, _ := helper.GetJwt(c)

	// get sub category by user from service
	subCategories, err := sc.SubCategoryService.GetSubCategoryByUser(userId)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail get sub category by user",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get sub category by user",
			"error":   err.Error(),
		})
	}

	// return response success
	return c.JSON(http.StatusOK, echo.Map{
		"message":     "success get sub category by user",
		"sub_categories": subCategories,
	})
}

// DeleteSubCategory is a function to delete sub category
func (sc *SubCategoryController) DeleteSubCategory(c echo.Context) error {
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

	// Call service to delete sub category
	err = sc.SubCategoryService.DeleteSubCategory(uint(id), userId)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail delete sub category",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail delete sub category",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success delete sub category",
	})
}

// UpdateSubCategory is a function to update sub category
func (sc *SubCategoryController) UpdateSubCategory(c echo.Context) error {
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
	var subCategory dto.SubCategoryDTO
	err = c.Bind(&subCategory)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail bind data",
			"error":   err.Error(),
		})
	}

	// id
	subCategory.ID = uint(id)
	
	// Get user id from jwt
	userId, _ := helper.GetJwt(c)
	subCategory.UserID = &userId

	// Call service to update sub category
	err = sc.SubCategoryService.UpdateSubCategory(subCategory)
	if err != nil {
		if val, ok := constantError.ErrorCode[err.Error()]; ok {
			return c.JSON(val, echo.Map{
				"message": "fail update sub category",
				"error":   err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail update sub category",
			"error":   err.Error(),
		})
	}

	// Return response if success
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update sub category",
	})
}
