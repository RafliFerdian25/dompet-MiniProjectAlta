package categoryRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type CategoryRepository interface {
	GetCategoryByID(id uint) (dto.Category, error)
	GetAllCategory() ([]dto.Category, error)
}