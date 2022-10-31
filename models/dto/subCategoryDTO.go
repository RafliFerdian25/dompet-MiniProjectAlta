package dto

type SubCategoryDTO struct {
	CategoryID uint
	UserID     uint   `json:"user_id"`
	Name       string `json:"name" gorm:"notNull;"`
}
