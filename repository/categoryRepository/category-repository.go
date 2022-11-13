package categoryRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
)

type CategoryRepository interface {
	GetCategoryByID(id uint) (model.Category, error)
	GetAllCategory() ([]dto.Category, error)
}