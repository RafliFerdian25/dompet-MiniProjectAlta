package subCategoryMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type SubCategoryMock struct {
	mock.Mock
}

func (sc *SubCategoryMock)GetSubCategoryById(id uint) (dto.SubCategoryDTO, error) {
	args := sc.Called(id)

	return args.Get(0).(dto.SubCategoryDTO), args.Error(1)
}

func (sc *SubCategoryMock)GetSubCategoryByUser(userId uint) ([]dto.SubCategoryDTO, error) {
	args := sc.Called(userId)

	return args.Get(0).([]dto.SubCategoryDTO), args.Error(1)
}

func (sc *SubCategoryMock)CreateSubCategory(subCategory dto.SubCategoryDTO) error {
	args := sc.Called(subCategory)

	return args.Error(0)
}

func (sc *SubCategoryMock)DeleteSubCategory(id uint) error {
	args := sc.Called(id)

	return args.Error(0)
}

func (sc *SubCategoryMock)UpdateSubCategory(subCategory dto.SubCategoryDTO) error {
	args := sc.Called(subCategory)

	return args.Error(0)
}
