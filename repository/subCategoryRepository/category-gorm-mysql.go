package subCategoryRepository

import (
	"gorm.io/gorm"
)

type subCategoryRepository struct {
	db *gorm.DB
}

// CreateSubCategory is used to create new sub category

// func NewCategoryRepository(db *gorm.DB) SubCategoryRepository {
// 	return &subCategoryRepository{
// 		db: db,
// 	}
// }