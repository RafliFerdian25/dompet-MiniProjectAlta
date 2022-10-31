package categoryRepository

import (
	"dompet-miniprojectalta/models/model"
)

type CategoryRepository interface {
	GetAllCategory() ([]model.Category, error)
}