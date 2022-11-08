package userMockService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"

	"github.com/stretchr/testify/mock"
)

type UserMock struct {
	mock.Mock
}

func (u *UserMock) CreateUser(user dto.UserDTO) error {
	args := u.Called(user)

	return args.Error(0)
}

func (u *UserMock) LoginUser(user model.User) (model.User, error) {
	args := u.Called(user)

	return args.Get(0).(model.User), args.Error(1)
}