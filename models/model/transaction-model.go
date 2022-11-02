package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID        uint    `json:"user_id" gorm:"notNull"`
	SubCategoryID uint    `json:"category_id" gorm:"notNull"`
	AccountID     uint    `json:"account_id" gorm:"notNull"`
	DebtID        uint    `json:"debt_id" gorm:"default:NULL"`
	Amount        float64 `json:"amount" gorm:"notNull"`
	Note          string  `json:"note"`
}
