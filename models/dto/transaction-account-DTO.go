package dto

type TransactionAccount struct {
	ID            uint    `json:"id"`
	AccountFromID uint    `json:"account_from_id" validate:"required,numeric"`
	AccountToID   uint    `json:"account_to_id" validate:"required,numeric"`
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount" validate:"required,numeric,gt=0"`
	Note          string  `json:"note"`
	AdminFee      float64 `json:"admin_fee"`
}
