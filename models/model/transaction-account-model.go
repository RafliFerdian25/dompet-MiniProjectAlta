package model

type TransactionAccount struct {
	ID            string    `gorm:"primarykey;size:255"`
	TimeModel     TimeModel `gorm:"embedded"`
	AccountFromID uint      `json:"account_from_id" gorm:"notNull;size:255"`
	AccountToID   uint      `json:"account_to_id" gorm:"notNull;size:255"`
	UserID        string    `json:"user_id" gorm:"notNull;size:255"`
	Amount        float64   `json:"amount" gorm:"notNull"`
	Note          string    `json:"note"`
	AdminFee      float64   `json:"admin_fee"`
}
