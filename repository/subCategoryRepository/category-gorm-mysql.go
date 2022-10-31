package subCategoryRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"

	"gorm.io/gorm"
)

type subCategoryRepository struct {
	db *gorm.DB
}

// GetSubCategoryByCategoryID implements SubCategoryRepository
func (s *subCategoryRepository) GetSubCategoryByCategoryID(id uint) ([]dto.SubCategoryDTO, error) {
	var subCategories []dto.SubCategoryDTO
	err := s.db.Model(&model.SubCategory{}).Where("category_id = ?", id).Find(&subCategories).Error
	if err != nil {
		return nil, err
	}
	return subCategories, nil
}

// CreateSubCategory implements SubCategoryRepository
func (s *subCategoryRepository) CreateSubCategory(subCategory dto.SubCategoryDTO) error {
	err := s.db.Model(&model.SubCategory{}).Create(&model.SubCategory{
		CategoryID: subCategory.CategoryID,
		UserID:     subCategory.UserID,
		Name:       subCategory.Name,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func NewCategoryRepository(db *gorm.DB) SubCategoryRepository {
	return &subCategoryRepository{
		db: db,
	}
}
