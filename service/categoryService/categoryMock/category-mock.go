package categoryMockService

import (
	"dompet-miniprojectalta/models/model"

	"github.com/stretchr/testify/mock"
)

type CategoryMock struct {
	mock.Mock
}

func (u *CategoryMock) GetCategoryByID(id uint) ([]model.Category, error) {
	args := u.Called(id)

	return args.Get(0).([]model.Category), args.Error(1)
}

func (u *CategoryMock) GetAllCategory() ([]model.Category, error) {
	args := u.Called()

	return args.Get(0).([]model.Category), args.Error(1)
}