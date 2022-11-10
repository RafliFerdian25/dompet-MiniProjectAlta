package accountMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type AccountMock struct {
	mock.Mock
}

func (u *AccountMock) DeleteAccount(id uint) error {
	args := u.Called(id)

	return args.Error(0)
}

func (u *AccountMock) GetAccountByUser(userId uint) ([]dto.AccountDTO, error) {
	args := u.Called(userId)

	return args.Get(0).([]dto.AccountDTO), args.Error(1)
}

func (u *AccountMock) GetAccountById(id uint) (dto.AccountDTO, error) {
	args := u.Called(id)
	
	return args.Get(0).(dto.AccountDTO), args.Error(1)
}

func (u *AccountMock) UpdateAccount(account dto.AccountDTO) error {
	args := u.Called(account)

	return args.Error(0)
}

func (u *AccountMock) CreateAccount(account dto.AccountDTO) error {
	args := u.Called(account)

	return args.Error(0)
}