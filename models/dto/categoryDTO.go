package dto

import "dompet-miniprojectalta/models/model"

type CategoryDTO struct {
	ID            uint   `json:"id"`
	Name          string `json:"name" gorm:"notNull"`
	SubCategoryID []model.SubCategory
}