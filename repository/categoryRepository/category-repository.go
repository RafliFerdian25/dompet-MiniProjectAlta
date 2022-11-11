package categoryRepository

import (
	"dompet-miniprojectalta/models/model"
)

type CategoryRepository interface {
	GetCategoryByID(id uint) (model.Category, error)
	GetAllCategory() ([]model.Category, error)
}