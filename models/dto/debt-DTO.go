package dto

import (
	"dompet-miniprojectalta/models/model"
	"time"
)

type Debt struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Total        float64 `json:"total"`
	Remaining    float64 `json:"remaining"`
	Note         string  `json:"note"`
	DebtStatus   string  `json:"debt_status"`
	Status       string  `json:"status"`
	Transactions []model.Transaction
}

type DebtTransactionDTO struct {
	Name          string  `json:"name" validate:"required"`
	UserID        uint    `json:"user_id"`
	DebtID        uint    `json:"debt_id"`
	SubCategoryID uint    `json:"sub_category_id" validate:"required,numeric"`
	AccountID     uint    `json:"account_id" validate:"required,numeric"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
	Note          string  `json:"note"`
}

type GetDebtTransactionResponse struct {
	ID        uint      `json:"id"`
	Name          string    `json:"name_lender" `
	SubCategoryID uint      `json:"sub_category_id" `
	AccountID     uint      `json:"account_id" `
	Total        float64   `json:"total"`
	Remaining    float64   `json:"remaining"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
	DebtStatus    string    `json:"debt_status"`
}
