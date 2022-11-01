package model

type Transaction struct {
	ID            string    `gorm:"primarykey;size:255"`
	TimeModel     TimeModel `gorm:"embedded"`
	UserID        string    `json:"user_id" gorm:"notNull;size:255"`
	SubCategoryID string    `json:"category_id" gorm:"notNull;size:255"`
	AccountID     string    `json:"account_id" gorm:"notNull;size:255"`
	DebtID        string    `json:"debt_id" gorm:"size:255"`
	Amount        float64   `json:"amount" gorm:"notNull"`
	Note          string    `json:"note"`
}
