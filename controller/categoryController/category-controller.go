package categoryController

import (
	"dompet-miniprojectalta/service/categoryService"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	CategoryService categoryService.CategoryService
}

func (b *CategoryController) GetAllCategory(c echo.Context) error {
	categorys, err := b.CategoryService.GetAllCategory()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "fail get all category",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":   "success",
		"categorys": categorys,
	})
}
