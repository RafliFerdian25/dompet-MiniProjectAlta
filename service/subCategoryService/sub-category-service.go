package subCategoryService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/subCategoryRepository"
)

type SubCategoryService interface {
	CreateSubCategory(subCategory dto.SubCategoryDTO) error
}

type subCategoryService struct {
	subCategoryRepository subCategoryRepository.SubCategoryRepository
}

// CreateSubCategory implements SubCategoryService
func (s *subCategoryService) CreateSubCategory(subCategory dto.SubCategoryDTO) error {
	err := s.subCategoryRepository.CreateSubCategory(subCategory)
	if err != nil {
		return err
	}
	return nil
}

func NewSubCategoryService(subCategoryRepository subCategoryRepository.SubCategoryRepository) SubCategoryService {
	return &subCategoryService{
		subCategoryRepository: subCategoryRepository,
	}
}