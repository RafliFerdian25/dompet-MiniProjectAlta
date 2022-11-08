package debtMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type DebtMock struct {
	mock.Mock
}

func (d *DebtMock) GetDebt(userId uint, debtStatus string) (map[string]interface{}, error) {
	args := d.Called(userId, debtStatus)

	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (d *DebtMock) DeleteDebt(id, userId uint) error {
	args := d.Called(id, userId)

	return args.Error(0)
}

func (d *DebtMock) CreateDebt(debtTransaction dto.DebtTransactionDTO) error {
	args := d.Called(debtTransaction)

	return args.Error(0)
}