package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID uint `json:"user_id"`
	Name     string `json:"name" gorm:"notNull"`
	Balance float64 `json:"balance" gorm:"notNull"`
	TransactionID []Transaction
	TransactionAccountFromID []TransactionAccount `json:"transaction_account_from_id" gorm:"foreignKey:AccountFromID"`
	TransactionAccountToID []TransactionAccount `json:"transaction_account_to_id" gorm:"foreignKey:AccountToID"`
}
