package model

type Category struct {
	ID            string    `gorm:"primarykey;size:255"`
	TimeModel     TimeModel `gorm:"embedded"`
	Name          string    `json:"name" gorm:"notNull;size:20"`
	SubCategories []SubCategory
}
