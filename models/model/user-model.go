package model

type User struct {
	ID                   string    `gorm:"primarykey;size:255"`
	TimeModel            TimeModel `gorm:"embedded"`
	Name                 string    `json:"name" gorm:"notNull;size:255"`
	Email                string    `json:"email" gorm:"notNull;unique;size:255"`
	Password             string    `json:"password" gorm:"notNull"`
	SubCategoryID        []SubCategory
	TransactionAccountID []TransactionAccount
	AccountID            []Account
	TransactionID        []Transaction
}
