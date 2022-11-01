package model

type Account struct {
	ID                       string    `gorm:"primarykey;size:255"`
	TimeModel                TimeModel `gorm:"embedded"`
	UserID                   string    `json:"user_id" gorm:"size:255"`
	Name                     string    `json:"name" gorm:"notNull;size:255"`
	Balance                  float64   `json:"balance" gorm:"notNull"`
	TransactionID            []Transaction
	TransactionAccountFromID []TransactionAccount `json:"transaction_account_from_id" gorm:"foreignKey:AccountFromID"`
	TransactionAccountToID   []TransactionAccount `json:"transaction_account_to_id" gorm:"foreignKey:AccountToID"`
}
