package model

import "gorm.io/gorm"

type Debt struct {
	gorm.Model
	Name          string  `json:"name" gorm:"notNull;size:255"`
	Total         float64 `json:"total" gorm:"notNull;default:0"`
	Remaining     float64 `json:"remaining" gorm:"notNull"`
	Note          string  `json:"note"`
	DebtStatus    string  `json:"debt_status" gorm:"notNull;size:10"`
	Status        string  `json:"status" gorm:"size:10"`
	Transaction []Transaction
}
