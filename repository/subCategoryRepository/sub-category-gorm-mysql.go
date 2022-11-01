package subCategoryRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"errors"

	"gorm.io/gorm"
)

type subCategoryRepository struct {
	db *gorm.DB
}

// GetSubCategoryById implements SubCategoryRepository
func (sc *subCategoryRepository) GetSubCategoryById(id uint) (dto.SubCategoryDTO, error) {
	var subCategory dto.SubCategoryDTO
	err := sc.db.Model(&model.SubCategory{}).First(&subCategory, id).Error
	if err != nil {
		return dto.SubCategoryDTO{}, err
	}
	return subCategory, nil
}

// GetSubCategoryByUser implements SubCategoryRepository
func (sc *subCategoryRepository) GetSubCategoryByUser(userId uint) ([]dto.SubCategoryDTO, error) {
	var SubCategoryUsers []dto.SubCategoryDTO
	// get data sub category from database by user
	err := sc.db.Model(&model.SubCategory{}).Where("user_id = ?", userId).Or("user_id IS NULL").Find(&SubCategoryUsers).Error
	if err != nil {
		return nil, err
	}
	return SubCategoryUsers, nil
}

// CreateSubCategory implements SubCategoryRepository
func (s *subCategoryRepository) CreateSubCategory(subCategory dto.SubCategoryDTO) error {
	// create new subcategory
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

// DeleteSubCategory implements SubCategoryRepository
func (sc *subCategoryRepository) DeleteSubCategory(id uint) error {
	// delete data sub category from database by id
	err := sc.db.Model(&model.SubCategory{}).Where("id = ?", id).Delete(&model.SubCategory{})
	if err.RowsAffected <= 0 {
		return errors.New("subcategory not found")
	}

	return nil
}

// UpdateSubCategory implements SubCategoryRepository
func (sc *subCategoryRepository) UpdateSubCategory(subCategory dto.SubCategoryDTO) error {
	// update subcategory with new data
	err := sc.db.Model(&model.SubCategory{}).Where("id = ?", subCategory.ID).Updates(&model.SubCategory{
		CategoryID: subCategory.CategoryID,
		Name:       subCategory.Name,
	})
	if err.Error != nil {
		return err.Error
	}
	if err.RowsAffected <= 0 {
		return errors.New("subcategory not found")
	}

	return nil
}

func NewCategoryRepository(db *gorm.DB) SubCategoryRepository {
	return &subCategoryRepository{
		db: db,
	}
}
