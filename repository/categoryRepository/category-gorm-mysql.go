package categoryRepository

import (
	"dompet-miniprojectalta/models/model"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// GetAllCategory implements CategoryRepository
func (b *categoryRepository) GetAllCategory() ([]model.Category, error) {
	var categories []model.Category
	err := b.db.Model(&model.Category{}).Preload("SubCategories").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}