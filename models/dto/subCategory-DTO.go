package dto

type SubCategoryDTO struct {
	ID         uint   `json:"id"`
	CategoryID uint   `json:"category_id" validate:"required,numeric"`
	UserID     string `json:"user_id" validate:"required,numeric"`
	Name       string `json:"name" gorm:"notNull;" validate:"required"`
}
