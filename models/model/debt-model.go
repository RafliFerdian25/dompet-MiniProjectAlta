package model

type Debt struct {
	ID            string    `gorm:"primarykey;size:255"`
	TimeModel     TimeModel `gorm:"embedded"`
	Name          string    `json:"name" gorm:"notNull;size:255"`
	Total         float64   `json:"total" gorm:"notNull"`
	Remaining     float64   `json:"remaining" gorm:"notNull"`
	Note          string    `json:"note"`
	DebtStatus    string    `json:"payment_status" gorm:"notNull;size:10"`
	Status        string    `json:"status" gorm:"size:10"`
	TransactionID []Transaction
}
