package dto

type Category struct {
	ID            uint          `json:"id"`
	Name          string        `json:"name"`
	SubCategories []SubCategory `json:"sub_categories"`
}
