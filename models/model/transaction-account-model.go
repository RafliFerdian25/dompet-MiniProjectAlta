package model

import "gorm.io/gorm"

type TransactionAccount struct {
	gorm.Model
	AccountFromID uint `json:"account_from_id" gorm:"notNull"`
	AccountToID uint `json:"account_to_id" gorm:"notNull"`
	UserID uint `json:"user_id" gorm:"notNull"`
	Amount float64 `json:"amount" gorm:"notNull"`
	Note string `json:"note"`
	AdminFee float64 `json:"admin_fee"`
}
