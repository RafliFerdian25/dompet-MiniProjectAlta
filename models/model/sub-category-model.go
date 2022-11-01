package model

type SubCategory struct {
	ID            string    `gorm:"primarykey;size:255"`
	TimeModel     TimeModel `gorm:"embedded"`
	CategoryID    uint      `json:"category_id" gorm:"notNull;size:255"`
	UserID        string    `json:"user_id" gorm:"size:255"`
	Name          string    `json:"name" gorm:"notNull;size:255"`
	TransactionID []Transaction
}
