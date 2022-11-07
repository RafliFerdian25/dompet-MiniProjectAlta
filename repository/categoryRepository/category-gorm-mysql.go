package categoryRepository

import (
	"dompet-miniprojectalta/models/model"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}
// GetCategoryByID implements categoryRepository
func (s *categoryRepository) GetCategoryByID(id uint) ([]model.Category, error) {
	var CategoriesID []model.Category
	// get data category from database by id
	err := s.db.Model(&model.Category{}).Preload("SubCategories").Where("id = ?", id).Find(&CategoriesID)
	if err.Error != nil {
		return nil, err.Error
	}
	if err.RowsAffected <= 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return CategoriesID, nil
}

// GetAllCategory implements CategoryRepository
func (b *categoryRepository) GetAllCategory() ([]model.Category, error) {
	var categories []model.Category
	// get all data category from database
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