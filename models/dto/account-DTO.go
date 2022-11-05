package dto

type AccountDTO struct {
	ID      uint    `json:"id"`
	UserID  uint    `json:"user_id"`
	Name    string  `json:"name" validate:"required"`
	Balance float64 `json:"balance" validate:"required,numeric,gte=0"`
}
