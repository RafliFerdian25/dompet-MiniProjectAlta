package dto

type SubCategoryDTO struct {
	CategoryID uint   `json:"category_id" validate:"required,numeric"`
	UserID     uint   `json:"user_id" validate:"required,numeric"`
	Name       string `json:"name" gorm:"notNull;" validate:"required"`
}
