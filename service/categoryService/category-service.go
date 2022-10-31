package categoryService

import (
	"dompet-miniprojectalta/models/model"
	"dompet-miniprojectalta/repository/categoryRepository"
)

type CategoryService interface {
	GetAllCategory() ([]model.Category, error)
}

type categoryService struct {
	categoryRepository categoryRepository.CategoryRepository
}

// GetAllCategory implements CategoryService
func (c *categoryService) GetAllCategory() ([]model.Category, error) {
	categories, err := c.categoryRepository.GetAllCategory()
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
