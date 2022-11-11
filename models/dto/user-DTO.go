package dto

type UserDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponseDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}