package categoryService

import (
	"dompet-miniprojectalta/models/model"
	"dompet-miniprojectalta/repository/categoryRepository"
)

type CategoryService interface {
	GetCategoryByCategoryID(id uint) ([]model.Category, error)
	GetAllCategory() ([]model.Category, error)
}

type categoryService struct {
	categoryRepository categoryRepository.CategoryRepository
}

// GetCategoryByCategoryID implements CategoryService
func (cs *categoryService) GetCategoryByCategoryID(id uint) ([]model.Category, error) {
	categoriesID, err := cs.categoryRepository.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	return categoriesID, nil
}

// GetAllCategory implements CategoryService
func (cs *categoryService) GetAllCategory() ([]model.Category, error) {
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
