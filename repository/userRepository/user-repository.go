package userRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
)

type UserRepository interface {
	CreateUser(user dto.UserDTO) error
	LoginUser(user model.User) (model.User, error)
}