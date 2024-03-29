package categoryMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type CategoryMock struct {
	mock.Mock
}

func (u *CategoryMock) GetCategoryByID(id uint) (dto.Category, error) {
	args := u.Called(id)

	return args.Get(0).(dto.Category), args.Error(1)
}

func (u *CategoryMock) GetAllCategory() ([]dto.Category, error) {
	args := u.Called()

	return args.Get(0).([]dto.Category), args.Error(1)
}