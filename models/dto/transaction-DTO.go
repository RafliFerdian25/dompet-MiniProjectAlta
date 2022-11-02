package dto

type TransactionDTO struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	SubCategoryID uint    `json:"sub_category_id" validate:"required"`
	AccountID     uint    `json:"account_id" validate:"required"`
	DebtID        uint    `json:"debt_id"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
	Note          string  `json:"note"`
}
