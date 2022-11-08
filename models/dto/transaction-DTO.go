package dto

import "time"

type TransactionDTO struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	SubCategoryID uint    `json:"sub_category_id" validate:"required"`
	AccountID     uint    `json:"account_id" validate:"required"`
	DebtID        uint    `json:"debt_id"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
	Note          string  `json:"note"`
}

type TransactionJoin struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	SubCategoryID uint    `json:"sub_category_id"`
	CategoryID    uint    `json:"category_id"`
	AccountID     uint    `json:"account_id"`
	Amount        float64 `json:"amount"`
}

type GetTransactionDTO struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	SubCategoryID uint    `json:"sub_category_id"`
	CategoryID    uint    `json:"category_id"`
	AccountID     uint    `json:"account_id"`
	Amount        float64 `json:"amount"`
	Note          string  `json:"note"`
	CreatedAt	 time.Time `json:"created_at"`
}

type TransactionReport struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	SubCategoryID uint      `json:"sub_category_id"`
	CategoryID    uint      `json:"category_id"`
	AccountID     uint      `json:"account_id"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransactionReportPeriod struct {
	Period     string `json:"period"`
	Total      int64  `json:"total"`
}

type ReportSpendingCategoryPeriod struct {
	SubCategory string  `json:"sub_category"`
	Period      string  `json:"period"`
	Total       int64   `json:"total"`
	Persentage  float64 `json:"persentage"`
}
