package userRepository

import (
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/model"
	"errors"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// CreateUser implements UserRepository
func (u *userRepository) CreateUser(user model.User) error {
	err := u.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// LoginUser implements UserRepository
func (u *userRepository) LoginUser(user model.User) (model.User, error) {
	var userLogin model.User
	err := u.db.Model(&model.User{}).Where("email = ?", user.Email).Find(&userLogin).Error
	if err != nil {
		return model.User{}, err
	}
	match := helper.CheckPasswordHash(user.Password, userLogin.Password)
	if !match {
		return model.User{}, errors.New("email or password not match")
	}
	return userLogin, nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
