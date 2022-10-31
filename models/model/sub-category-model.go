package model

import "gorm.io/gorm"

type SubCategory struct {
	gorm.Model
	CategoryID uint
	UserID uint `json:"user_id"`
	Name     string `json:"name" gorm:"notNull"`
	TransactionID []Transaction
}
