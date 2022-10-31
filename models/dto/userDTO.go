package dto

import (
	"golang.org/x/crypto/bcrypt"
)
type UserDTO struct {
	ID   uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}