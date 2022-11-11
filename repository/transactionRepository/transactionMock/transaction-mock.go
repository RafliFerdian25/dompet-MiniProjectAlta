package transactionMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type TransactionMock struct {
	mock.Mock
}

func (t *TransactionMock) GetTransaction(userId, categoryID uint, month int) ([]dto.GetTransactionDTO, error) {
	args := t.Called(userId, categoryID, month)

	return args.Get(0).([]dto.GetTransactionDTO), args.Error(1)
}

func (t *TransactionMock) DeleteTransaction(id uint,  account dto.AccountDTO) error {
	args := t.Called(id, account)

	return args.Error(0)
}

func (t *TransactionMock) UpdateTransaction(newTransaction dto.TransactionDTO, oldTransaction dto.TransactionJoin, newAccount, oldAccount dto.AccountDTO) error {
	args := t.Called(newTransaction, oldTransaction, newAccount, oldAccount)

	return args.Error(0)
}

func (t *TransactionMock) GetTransactionById(id uint) (dto.TransactionJoin, error) {
	args := t.Called(id)
	
	return args.Get(0).(dto.TransactionJoin), args.Error(1)
}

func (t *TransactionMock) CreateTransaction(transaction dto.TransactionDTO, categoryID uint, account dto.AccountDTO) error {
	args := t.Called(transaction, categoryID, account)

	return args.Error(0)
}