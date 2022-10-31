package subCategoryRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type SubCategoryRepository interface {
	GetSubCategoryByCategoryID(id uint) ([]dto.SubCategoryDTO, error)
	CreateSubCategory(subCategory dto.SubCategoryDTO) error
	// UpdateSubCategory(subCategory dto.SubCategoryDTO) error
	// DeleteSubCategory(id string) error
}