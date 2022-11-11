package accountService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	accountMockRepository "dompet-miniprojectalta/repository/accountRepository/accountMock"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteAccounts struct {
	suite.Suite
	accountService *accountService
	mock           *accountMockRepository.AccountMock
}

func (s *suiteAccounts) SetupSuite() {
	mock := &accountMockRepository.AccountMock{}
	s.mock = mock
	s.accountService = &accountService{
		accRepo: s.mock,
	}
}

func (s *suiteAccounts) TestDeleteAccount() {
	testCase := []struct {
		Name                  string
		GetMockReturnError    error
		GetMockReturnBody     dto.AccountDTO
		DeleteMockReturnError error
		ParamId               uint
		ParamUserId           uint
		HasReturnError        bool
		ExpectedError         error
	}{
		{
			"success",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
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
			dto.AccountDTO{},
			nil,
			1,
			1,
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth",
			nil,
			dto.AccountDTO{},
			nil,
			1,
			2,
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"failed delete account",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			gorm.ErrRecordNotFound,
			1,
			1,
			true,
			gorm.ErrRecordNotFound,
		},
	}

	for _, v := range testCase {
		var mockCallGet = s.mock.On("GetAccountById", v.ParamId).Return(v.GetMockReturnBody, v.GetMockReturnError)
		var mockCallDelete = s.mock.On("DeleteAccount", v.ParamId).Return(v.DeleteMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.accountService.DeleteAccount(v.ParamId, v.ParamUserId)

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

func (s *suiteAccounts) TestGetAccountByUser() {
	testCase := []struct {
		Name            string
		MockReturnError error
		MockReturnBody  []dto.AccountDTO
		ParamUserId     uint
		HasReturnBody   bool
		ExpectedBody    []dto.AccountDTO
	}{
		{
			"success",
			nil,
			[]dto.AccountDTO{
				{
					ID:      1,
					UserID:  1,
					Name:    "test",
					Balance: 10000,
				},
			},
			1,
			true,
			[]dto.AccountDTO{
				{
					ID:      1,
					UserID:  1,
					Name:    "test",
					Balance: 10000,
				},
			},
		},
		{
			"failed get account",
			errors.New("error"),
			[]dto.AccountDTO{},
			1,
			false,
			[]dto.AccountDTO{},
		},
	}

	for _, v := range testCase {
		var mockCall = s.mock.On("GetAccountByUser", v.ParamUserId).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			accounts, err := s.accountService.GetAccountByUser(v.ParamUserId)

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

func (s *suiteAccounts) TestUpdateAccount() {
	testCase := []struct {
		Name                  string
		GetMockReturnError    error
		GetMockReturnBody     dto.AccountDTO
		UpdateMockReturnError error
		Body                  dto.AccountDTO
		HasReturnError        bool
		ExpectedError         error
	}{
		{
			"success",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:     1,
				UserID: 1,
				Name:   "Mandiri",
			},
			false,
			nil,
		},
		{
			"failed get account",
			errors.New(constantError.ErrorAccountNotFound),
			dto.AccountDTO{},
			nil,
			dto.AccountDTO{
				ID:     1,
				UserID: 1,
				Name:   "Mandiri",
			},
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  2,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:     1,
				UserID: 1,
				Name:   "Mandiri",
			},
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"balace less than 0",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "Mandiri",
				Balance: -10000,
			},
			true,
			errors.New(constantError.ErrorBalanceLessThanZero),
		},
		{
			"failed update account",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			gorm.ErrRecordNotFound,
			dto.AccountDTO{
				ID:     1,
				UserID: 1,
				Name:   "Mandiri",
			},
			true,
			gorm.ErrRecordNotFound,
		},
	}

	for _, v := range testCase {
		var mockCallGet = s.mock.On("GetAccountById", v.Body.ID).Return(v.GetMockReturnBody, v.GetMockReturnError)
		var mockCallUpdate = s.mock.On("UpdateAccount", v.Body).Return(v.UpdateMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.accountService.UpdateAccount(v.Body)

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

func (s *suiteAccounts) TestCreateAccount() {
	testCase := []struct {
		Name            string
		MockReturnError error
		Body            dto.AccountDTO
		HasReturnError  bool
		ExpectedError   error
	}{
		{
			"success",
			nil,
			dto.AccountDTO{
				UserID: 1,
				Name:   "Mandiri",
				Balance: 10000,
			},
			false,
			nil,
		},
		{
			"balace less than 0",
			nil,
			dto.AccountDTO{
				UserID: 1,
				Name:   "Mandiri",
				Balance: -10000,
			},
			true,
			errors.New(constantError.ErrorBalanceLessThanZero),
		},
		{
			"failed create account",
			errors.New("error"),
			dto.AccountDTO{
				UserID: 1,
				Name:   "Mandiri",
				Balance: 10000,
			},
			true,
			errors.New("error"),
		},
	}

	for _, v := range testCase {
		var mockCall = s.mock.On("CreateAccount", v.Body).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.accountService.CreateAccount(v.Body)
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

func TestSuiteAccounts(t *testing.T) {
	suite.Run(t, new(suiteAccounts))
}
