package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name                 string `json:"name" gorm:"notNull;size:255"`
	Email                string `json:"email" gorm:"notNull;unique;size:255"`
	Password             string `json:"password" gorm:"notNull"`
	SubCategoryID        []SubCategory
	TransactionAccountID []TransactionAccount
	AccountID            []Account
	TransactionID        []Transaction
}
