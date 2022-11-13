package categoryService

import (
	"dompet-miniprojectalta/models/dto"
	categoryMockRepository "dompet-miniprojectalta/repository/categoryRepository/categoryMock"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type suiteCategorys struct {
	suite.Suite
	categoryService CategoryService
	mock            *categoryMockRepository.CategoryMock
}

func (s *suiteCategorys) SetupSuite() {
	mock := &categoryMockRepository.CategoryMock{}
	s.mock = mock
	NewCategoryService := NewCategoryService(s.mock)
	s.categoryService = NewCategoryService
}

func (s *suiteCategorys) TestGetCategoryByID() {
	testCase := []struct {
		Name            string
		MockReturnError error
		MockReturnBody  dto.Category
		ParamId         uint
		HasReturnBody   bool
		ExpectedBody    dto.Category
	}{
		{
			"success",
			nil,
			dto.Category{
					ID: 1,
				Name:          "test",
				SubCategories: []dto.SubCategory{},
			},
			1,
			true,
			dto.Category{
					ID: 1,
				Name:          "test",
				SubCategories: []dto.SubCategory{},
			},
		},
		{
			"failed get account",
			errors.New("error"),
			dto.Category{},
			1,
			false,
			dto.Category{},
		},
	}

	for _, v := range testCase {
		var mockCall = s.mock.On("GetCategoryByID", v.ParamId).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			accounts, err := s.categoryService.GetCategoryByID(v.ParamId)

			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, accounts)
			} else {
				s.Error(err)
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteCategorys) TestGetAllCategory() {
	testCase := []struct {
		Name            string
		MockReturnError error
		MockReturnBody  []dto.Category
		ParamUserId     uint
		HasReturnBody   bool
		ExpectedBody    []dto.Category
	}{
		{
			"success",
			nil,
			[]dto.Category{
				{
					ID:            1,
					Name:          "test",
					SubCategories: []dto.SubCategory{},
				},
			},
			1,
			true,
			[]dto.Category{
				{
					ID:            1,
					Name:          "test",
					SubCategories: []dto.SubCategory{},
				},
			},
		},
		{
			"failed get account",
			errors.New("error"),
			[]dto.Category{},
			1,
			false,
			[]dto.Category{},
		},
	}

	for _, v := range testCase {
		var mockCall = s.mock.On("GetAllCategory").Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			accounts, err := s.categoryService.GetAllCategory()

			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, accounts)
			} else {
				s.Error(err)
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func TestSuiteCategorys(t *testing.T) {
	suite.Run(t, new(suiteCategorys))
}
