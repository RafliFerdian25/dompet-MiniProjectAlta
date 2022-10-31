package categoryService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/categoryRepository"
)

type CategoryService interface {
	GetAllCategory() ([]dto.CategoryDTO, error)
}

type categoryService struct {
	categoryRepository categoryRepository.CategoryRepository
}

// GetAllCategory implements CategoryService
func (c *categoryService) GetAllCategory() ([]dto.CategoryDTO, error) {
	categorys, err := c.categoryRepository.GetAllCategory()
	if err != nil {
		return nil, err
	}
	return categorys, nil
}

func NewCategoryService(categoryRepository categoryRepository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
