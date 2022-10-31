package userRepository

import (
	"dompet-miniprojectalta/models/model"
)

type UserRepository interface {
	CreateUser(user model.User) error
	LoginUser(user model.User) (model.User, error)
}