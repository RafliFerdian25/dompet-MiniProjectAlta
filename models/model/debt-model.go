package model

import "gorm.io/gorm"

type Debt struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"notNull"`
	Name string `json:"name" gorm:"notNull"`
	Total float64 `json:"total" gorm:"notNull"`
	Remaining float64 `json:"remaining" gorm:"notNull"`
	Note string `json:"note"`
	PaymentStatus string `json:"payment_status"`
	Status string `json:"status"`
	TransactionID []Transaction
}
