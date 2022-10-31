package subCategoryRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type SubCategoryRepository interface {
	CreateSubCategory(subCategory dto.SubCategoryDTO) error
	// UpdateSubCategory(subCategory dto.SubCategoryDTO) error
	// DeleteSubCategory(id string) error
	// GetSubCategoryByCategoryID(id uint) ([]blogModel.SubCategoryDTO, error)
	// GetSubCategoryByTitle(title string) (blogModel.SubCategoryDTO, error)
}