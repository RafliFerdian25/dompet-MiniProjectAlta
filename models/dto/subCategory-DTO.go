package dto

type SubCategoryDTO struct {
	ID         uint   `json:"id"`
	CategoryID uint   `json:"category_id" validate:"required,numeric"`
	UserID     uint   `json:"user_id"`
	Name       string `json:"name" gorm:"notNull;" validate:"required"`
}
