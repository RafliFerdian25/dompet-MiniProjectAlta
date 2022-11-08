package transactionAccMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type TransactionAccMock struct {
	mock.Mock
}

func (ta *TransactionAccMock) GetTransactionAccount(userId uint, month int) ([]dto.GetTransactionAccountDTO, error) {
	args := ta.Called(userId, month)

	return args.Get(0).([]dto.GetTransactionAccountDTO), args.Error(1)
}

func (ta *TransactionAccMock) DeleteTransactionAccount(id, userID uint) error {
	args := ta.Called(id, userID)

	return args.Error(0)
}

func (ta *TransactionAccMock) CreateTransactionAccount(transAcc dto.TransactionAccount) error {
	args := ta.Called(transAcc)

	return args.Error(0)
}