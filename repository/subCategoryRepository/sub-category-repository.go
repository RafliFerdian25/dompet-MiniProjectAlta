package subCategoryRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type SubCategoryRepository interface {
	GetSubCategoryById(id uint) (dto.SubCategoryDTO, error)
	GetSubCategoryByUser(userId string) ([]dto.SubCategoryDTO, error)
	CreateSubCategory(subCategory dto.SubCategoryDTO) error
	DeleteSubCategory(id uint) error
	UpdateSubCategory(subCategory dto.SubCategoryDTO) error
}