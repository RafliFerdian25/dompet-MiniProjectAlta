package accountMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type UserMock struct {
	mock.Mock
}

func (u *UserMock) DeleteAccount(id, userId uint) error {
	args := u.Called(id, userId)

	return args.Error(0)
}

func (u *UserMock) GetAccountByUser(userId uint) ([]dto.AccountDTO, error) {
	args := u.Called(userId)

	return args.Get(0).([]dto.AccountDTO), args.Error(1)
}

func (u *UserMock) UpdateAccount(account dto.AccountDTO) error {
	args := u.Called(account)

	return args.Error(0)
}

func (u *UserMock) CreateAccount(account dto.AccountDTO) error {
	args := u.Called(account)

	return args.Error(0)
}