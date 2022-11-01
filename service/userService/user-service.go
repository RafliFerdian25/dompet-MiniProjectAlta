package userService

import (
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"dompet-miniprojectalta/repository/userRepository"
)

type UserService interface {
	CreateUser(user dto.UserDTO) error
	LoginUser(user model.User) (model.User, error)
}

type userServiceImpl struct {
	userRepo userRepository.UserRepository
}

// CreateUser implements UserService
func (u *userServiceImpl) CreateUser(user dto.UserDTO) error {
	// make uuid
	user.ID = helper.GenerateUUID()
	
	// hash password
	password, errPassword := helper.HashPassword(user.Password)
	user.Password = password
	if errPassword != nil {
		return errPassword
	}

	// call repository to save user
	err := u.userRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

// LoginUser implements UserService
func (u *userServiceImpl) LoginUser(user model.User) (model.User, error) {
	// call repository to get user
	userLogin, err := u.userRepo.LoginUser(user)
	if err != nil {
		return model.User{}, err
	}
	return userLogin, nil
}

func NewUserService(userRepository userRepository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepository,
	}
}
