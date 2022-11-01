package dto

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name" gorm:"notNull"`
}