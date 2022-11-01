package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name          string `json:"name" gorm:"notNull;size:20"`
	SubCategories []SubCategory
}
