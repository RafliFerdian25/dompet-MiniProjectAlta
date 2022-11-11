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

func (ta *TransactionAccMock) DeleteTransactionAccount(id uint, accountFrom, accountTo dto.AccountDTO) error {
	args := ta.Called(id, accountFrom, accountTo)

	return args.Error(0)
}

func (ta *TransactionAccMock) GetTransactionAccountById(id uint) (dto.TransactionAccount, error) {
	args := ta.Called(id)

	return args.Get(0).(dto.TransactionAccount), args.Error(1)
}

func (ta *TransactionAccMock) CreateTransactionAccount(transAcc dto.TransactionAccount, accountFrom, accountTo dto.AccountDTO) error {
	args := ta.Called(transAcc, accountFrom, accountTo)

	return args.Error(0)
}