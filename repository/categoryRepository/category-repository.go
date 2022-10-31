package categoryRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type CategoryRepository interface {
	GetAllCategory() ([]dto.CategoryDTO, error)
}