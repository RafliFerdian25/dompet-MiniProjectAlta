package model

import "gorm.io/gorm"

type SubCategory struct {
	gorm.Model
	CategoryID    uint   `json:"category_id" gorm:"notNull"`
	UserID        uint `json:"user_id"`
	Name          string `json:"name" gorm:"notNull;size:255"`
	TransactionID []Transaction
}
