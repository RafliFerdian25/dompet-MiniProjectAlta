package transactionMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type TransactionMock struct {
	mock.Mock
}

func (t *TransactionMock) GetTransaction(userId uint, month int) (map[string]interface{}, error) {
	args := t.Called(userId, month)

	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (t *TransactionMock) DeleteTransaction(id uint, userID uint) error {
	args := t.Called(id, userID)

	return args.Error(0)
}

func (t *TransactionMock) UpdateTransaction(transaction dto.TransactionDTO, userId uint) error {
	args := t.Called(transaction, userId)

	return args.Error(0)
}

func (t *TransactionMock) CreateTransaction(transaction dto.TransactionDTO) error {
	args := t.Called(transaction)

	return args.Error(0)
}