package debtMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type DebtMock struct {
	mock.Mock
}

func (d *DebtMock) GetDebt(userId uint, subCategory int, debtStatus string) ([]dto.GetDebtTransactionResponse, error) {
	args := d.Called(userId, subCategory ,debtStatus)

	return args.Get(0).([]dto.GetDebtTransactionResponse), args.Error(1)
}

func (d *DebtMock) DeleteDebt(id uint, account dto.AccountDTO) error {
	args := d.Called(id, account)

	return args.Error(0)
}

func (d *DebtMock) GetDebtById(id uint) (dto.Debt, error) {
	args := d.Called(id)

	return args.Get(0).(dto.Debt), args.Error(1)
}

func (d *DebtMock) CreateDebt(debt dto.Debt, transaction dto.TransactionDTO, account dto.AccountDTO) error {
	args := d.Called(debt, transaction, account)

	return args.Error(0)
}