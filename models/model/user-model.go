package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"notNull"` 
	Email    string `json:"email" gorm:"notNull"`
	Password string `json:"password" gorm:"notNull"`
	SubCategoryID []SubCategory
	TransactionAccountID []TransactionAccount
	AccountID []Account
	DebtID []Debt
	TransactionID []Transaction
}
