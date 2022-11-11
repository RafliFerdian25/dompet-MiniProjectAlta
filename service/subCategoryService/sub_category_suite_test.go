package subCategoryService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	subCategoryMockRepository "dompet-miniprojectalta/repository/subCategoryRepository/subCategoryMock"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteSubCategory struct {
	suite.Suite
	subCategoryService     *subCategoryService
	subCategoryMock *subCategoryMockRepository.SubCategoryMock
}

func (s *suiteSubCategory) SetupSuite() {
	subCategoryMock := &subCategoryMockRepository.SubCategoryMock{}
	s.subCategoryMock = subCategoryMock

	// panggil new
	s.subCategoryService = &subCategoryService{
		subCategoryRepository: s.subCategoryMock,
	}
}

func (s *suiteSubCategory) TestGetSubCategoryByUser() {
	testCase := []struct {
		Name            string
		MockReturnError error
		MockReturnBody  []dto.SubCategoryDTO
		ParamUserId     uint
		HasReturnBody   bool
		ExpectedBody    []dto.SubCategoryDTO
	}{
		{
			"success",
			nil,
			[]dto.SubCategoryDTO{
				{
					ID:      1,
					UserID:  1,
					CategoryID: 1,
					Name:    "test",
				},
			},
			1,
			true,
			[]dto.SubCategoryDTO{
				{
					ID:      1,
					UserID:  1,
					CategoryID: 1,
					Name:    "test",
				},
			},
		},
		{
			"failed get account",
			errors.New("error"),
			[]dto.SubCategoryDTO{},
			1,
			false,
			[]dto.SubCategoryDTO{},
		},
	}

	for _, v := range testCase {
		var mockCall = s.subCategoryMock.On("GetSubCategoryByUser", v.ParamUserId).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			accounts, err := s.subCategoryService.GetSubCategoryByUser(v.ParamUserId)

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

func (s *suiteSubCategory) TestCreateSubCategory() {
	testCase := []struct {
		Name            string
		MockReturnError error
		Body            dto.SubCategoryDTO
		HasReturnError  bool
		ExpectedError   error
	}{
		{
			"success",
			nil,
			dto.SubCategoryDTO{
					UserID:  1,
					CategoryID: 1,
					Name:    "test",
			},
			false,
			nil,
		},
		{
			"failed create account",
			errors.New("error"),
			dto.SubCategoryDTO{
				UserID:  1,
				CategoryID: 1,
				Name:    "test",
			},
			true,
			errors.New("error"),
		},
	}

	for _, v := range testCase {
		var mockCall = s.subCategoryMock.On("CreateSubCategory", v.Body).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.subCategoryService.CreateSubCategory(v.Body)
			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteSubCategory) TestDeleteSubCategory() {
	testCase := []struct {
		Name                  string
		GetMockReturnError    error
		GetMockReturnBody     dto.SubCategoryDTO
		DeleteMockReturnError error
		ParamId               uint
		ParamUserId           uint
		HasReturnError        bool
		ExpectedError         error
	}{
		{
			"success",
			nil,
			dto.SubCategoryDTO{
				ID: 	1,
				UserID:  1,
				CategoryID: 1,
				Name:    "test",
			},
			nil,
			1,
			1,
			false,
			nil,
		},
		{
			"failed get account",
			errors.New(constantError.ErrorAccountNotFound),
			dto.SubCategoryDTO{},
			nil,
			1,
			1,
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth",
			nil,
			dto.SubCategoryDTO{
				ID: 	1,
				UserID:  1,
				CategoryID: 1,
				Name:    "test",},
			nil,
			1,
			2,
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"failed delete account",
			nil,
			dto.SubCategoryDTO{
				ID: 	1,
				UserID:  1,
				CategoryID: 1,
				Name:    "test",
			},
			gorm.ErrRecordNotFound,
			1,
			1,
			true,
			gorm.ErrRecordNotFound,
		},
	}

	for _, v := range testCase {
		var mockCallGet = s.subCategoryMock.On("GetSubCategoryById", v.ParamId).Return(v.GetMockReturnBody, v.GetMockReturnError)
		var mockCallDelete = s.subCategoryMock.On("DeleteSubCategory", v.ParamId).Return(v.DeleteMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.subCategoryService.DeleteSubCategory(v.ParamId, v.ParamUserId)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGet.Unset()
		mockCallDelete.Unset()
	}
}

func (s *suiteSubCategory) TestUpdateSubCategory() {
	testCase := []struct {
		Name                  string
		GetMockReturnError    error
		GetMockReturnBody     dto.SubCategoryDTO
		UpdateMockReturnError error
		Body                  dto.SubCategoryDTO
		HasReturnError        bool
		ExpectedError         error
	}{
		{
			"success",
			nil,
			dto.SubCategoryDTO{
				ID:      1,
				CategoryID: 1,
				UserID:  1,
				Name:    "test",
			},
			nil,
			dto.SubCategoryDTO{
				ID:     1,
				CategoryID: 1,
				UserID: 1,
				Name:   "Salary",
			},
			false,
			nil,
		},
		{
			"failed get sub category",
			errors.New(constantError.ErrorAccountNotFound),
			dto.SubCategoryDTO{},
			nil,
			dto.SubCategoryDTO{
				ID:     1,
				CategoryID: 1,
				UserID: 1,
				Name:   "Salary",
			},
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth",
			nil,
			dto.SubCategoryDTO{
				ID:      1,
				CategoryID: 1,
				UserID:  2,
				Name:    "test",
			},
			nil,
			dto.SubCategoryDTO{
				ID:     1,
				CategoryID: 1,
				UserID: 1,
				Name:   "Salary",
			},
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"failed update sub category",
			nil,
			dto.SubCategoryDTO{
				ID:      1,
				CategoryID: 1,
				UserID:  1,
				Name:    "test",
			},
			gorm.ErrRecordNotFound,
			dto.SubCategoryDTO{
				ID:     1,
				CategoryID: 1,
				UserID: 1,
				Name:   "Salary",
			},
			true,
			gorm.ErrRecordNotFound,
		},
	}

	for _, v := range testCase {
		var mockCallGet = s.subCategoryMock.On("GetSubCategoryById", v.Body.ID).Return(v.GetMockReturnBody, v.GetMockReturnError)
		var mockCallUpdate = s.subCategoryMock.On("UpdateSubCategory", v.Body).Return(v.UpdateMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.subCategoryService.UpdateSubCategory(v.Body)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGet.Unset()
		mockCallUpdate.Unset()
	}
}

func TestSuiteSubCategory(t *testing.T) {
	suite.Run(t, new(suiteSubCategory))
}
