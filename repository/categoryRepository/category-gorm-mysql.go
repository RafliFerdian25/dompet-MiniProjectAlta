package categoryRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// GetAllCategory implements CategoryRepository
func (b *categoryRepository) GetAllCategory() ([]dto.CategoryDTO, error) {
	var categorys []dto.CategoryDTO
	err := b.db.Model(&model.Category{}).Find(&categorys).Error
	if err != nil {
		return nil, err
	}
	return categorys, nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}