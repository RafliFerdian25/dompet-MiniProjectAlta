package categoryService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/categoryRepository"
)

type CategoryService interface {
	GetCategoryByID(id uint) (dto.Category, error)
	GetAllCategory() ([]dto.Category, error)
}

type categoryService struct {
	categoryRepository categoryRepository.CategoryRepository
}

// GetCategoryByID implements CategoryService
func (cs *categoryService) GetCategoryByID(id uint) (dto.Category, error) {
	categoriesID, err := cs.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return dto.Category{}, err
	}
	return categoriesID, nil
}

// GetAllCategory implements CategoryService
func (cs *categoryService) GetAllCategory() ([]dto.Category, error) {
	categories, err := cs.categoryRepository.GetAllCategory()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func NewCategoryService(categoryRepository categoryRepository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
