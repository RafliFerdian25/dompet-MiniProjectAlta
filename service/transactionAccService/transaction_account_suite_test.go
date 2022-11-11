package transactionAccService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	accountMockRepository "dompet-miniprojectalta/repository/accountRepository/accountMock"
	transactionAccMockRepository "dompet-miniprojectalta/repository/transactionAccRepository/transactionAccMock"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteTransactionAcc struct {
	suite.Suite
	transactionAccService *transactionAccService
	transactionAccMock    *transactionAccMockRepository.TransactionAccMock
	accountMock           *accountMockRepository.AccountMock
}

func (s *suiteTransactionAcc) SetupSuite() {
	transactionAccMock := &transactionAccMockRepository.TransactionAccMock{}
	s.transactionAccMock = transactionAccMock

	accountMock := &accountMockRepository.AccountMock{}
	s.accountMock = accountMock

	s.transactionAccService = &transactionAccService{
		transAccRepo: s.transactionAccMock,
		accountRepo:  s.accountMock,
	}
}

func (s *suiteTransactionAcc) TestGetAccountByUser() {
	testCase := []struct {
		Name            string
		MockReturnError error
		MockReturnBody  []dto.GetTransactionAccountDTO
		ParamUserId     uint
		ParamMonth      int
		HasReturnBody   bool
		ExpectedBody    []dto.GetTransactionAccountDTO
	}{
		{
			"success",
			nil,
			[]dto.GetTransactionAccountDTO{
				{
					ID:            1,
					UserID:        1,
					AccountFromID: 1,
					AccountToID:   2,
					Amount:        10000,
					Note:          "test",
					AdminFee:      0,
				},
			},
			1,
			11,
			true,
			[]dto.GetTransactionAccountDTO{
				{
					ID:            1,
					UserID:        1,
					AccountFromID: 1,
					AccountToID:   2,
					Amount:        10000,
					Note:          "test",
					AdminFee:      0,
				},
			},
		},
		{
			"failed get account",
			errors.New("error"),
			[]dto.GetTransactionAccountDTO{},
			1,
			11,
			false,
			[]dto.GetTransactionAccountDTO{},
		},
	}

	for _, v := range testCase {
		var mockCall = s.transactionAccMock.On("GetTransactionAccount", v.ParamUserId, v.ParamMonth).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			transactionAccounts, err := s.transactionAccService.GetTransactionAccount(v.ParamUserId, v.ParamMonth)

			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, transactionAccounts)
			} else {
				s.Error(err)
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteTransactionAcc) TestDeleteAccount() {
	testCase := []struct {
		Name                          string
		GetTransactionMockReturnError error
		GetTransactionMockReturnBody  dto.TransactionAccount
		GetAccountFromMockReturnError error
		GetAccountFromMockReturnBody  dto.AccountDTO
		GetAccountToMockReturnError   error
		GetAccountToMockReturnBody    dto.AccountDTO
		DeleteMockReturnError         error
		ParamId                       uint
		ParamUserId                   uint
		HasReturnError                bool
		ExpectedError                 error
	}{
		{
			"success",
			nil,
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        1000,
				Note:          "test",
				AdminFee:      0,
			},
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
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
			nil,
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee:      0,
			},
			errors.New(constantError.ErrorAccountNotFound),
			dto.AccountDTO{},
			nil,
			dto.AccountDTO{},
			nil,
			1,
			1,
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"failed get account",
			nil,
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee:      0,
			},
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
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
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee:      0,
			},
			nil,
			dto.AccountDTO{},
			nil,
			dto.AccountDTO{},
			nil,
			1,
			2,
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"account To not enough balance",
			nil,
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        100000,
				Note:          "test",
				AdminFee:      0,
			},
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			1,
			1,
			true,
			errors.New(constantError.ErrorRecipientAccountNotEnoughBalance),
		},
		{
			"failed delete transaction account",
			nil,
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee:      0,
			},
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
		var mockCallGetTransactionAcc = s.transactionAccMock.On("GetTransactionAccountById", v.ParamId).Return(v.GetTransactionMockReturnBody, v.GetTransactionMockReturnError)
		var mockCallGetAccountFrom = s.accountMock.On("GetAccountById", v.GetTransactionMockReturnBody.AccountFromID).Return(v.GetAccountFromMockReturnBody, v.GetAccountFromMockReturnError)
		var mockCallGetAccountTo = s.accountMock.On("GetAccountById", v.GetTransactionMockReturnBody.AccountToID).Return(v.GetAccountToMockReturnBody, v.GetAccountToMockReturnError)
		
		v.GetAccountFromMockReturnBody.Balance = v.GetAccountFromMockReturnBody.Balance + v.GetTransactionMockReturnBody.Amount
		v.GetAccountToMockReturnBody.Balance = v.GetAccountToMockReturnBody.Balance - v.GetTransactionMockReturnBody.Amount
		var mockCallDelete = s.transactionAccMock.On("DeleteTransactionAccount", v.ParamId, v.GetAccountFromMockReturnBody, v.GetAccountToMockReturnBody).Return(v.DeleteMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.transactionAccService.DeleteTransactionAccount(v.ParamId, v.ParamUserId)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGetTransactionAcc.Unset()
		mockCallGetAccountFrom.Unset()
		mockCallGetAccountTo.Unset()
		mockCallDelete.Unset()
	}
}

func (s *suiteTransactionAcc) TestCreateTransactionAccount() {
	testCase := []struct {
		Name            string
		MockReturnError error
		GetAccountFromMockReturnError error
		GetAccountFromMockReturnBody  dto.AccountDTO
		GetAccountToMockReturnError   error
		GetAccountToMockReturnBody    dto.AccountDTO
		Body            dto.TransactionAccount
		HasReturnError  bool
		ExpectedError   error
	}{
		{
			"success",
			nil,
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			dto.TransactionAccount{
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        1000,
				Note:          "test",
				AdminFee:      0,
			},
			false,
			nil,
		},
		{
			"error auth",
			nil,
			errors.New(constantError.ErrorAccountNotFound),
			dto.AccountDTO{},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			dto.TransactionAccount{
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        1000,
				Note:          "test",
				AdminFee:      0,
			},
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth",
			nil,
			nil,
			dto.AccountDTO{},
			errors.New(constantError.ErrorAccountNotFound),
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			dto.TransactionAccount{
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        1000,
				Note:          "test",
				AdminFee:      0,
			},
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth",
			nil,
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  2,
				Name:    "test",
				Balance: 10000,
			},
			dto.TransactionAccount{
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        1000,
				Note:          "test",
				AdminFee:      0,
			},
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"balace less than 0",
			nil,
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        100000,
				Note:          "test",
				AdminFee:      0,
			},
			true,
			errors.New(constantError.ErrorAccountNotEnoughBalance),
		},
		{
			"failed create account",
			errors.New("error"),
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			dto.TransactionAccount{
				ID:            1,
				UserID:        1,
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        1000,
				Note:          "test",
				AdminFee:      0,
			},
			true,
			errors.New("error"),
		},
	}

	for _, v := range testCase {
		var mockCallGetAccountFrom = s.accountMock.On("GetAccountById", v.Body.AccountFromID).Return(v.GetAccountFromMockReturnBody, v.GetAccountFromMockReturnError)
		var mockCallGetAccountTo = s.accountMock.On("GetAccountById", v.Body.AccountToID).Return(v.GetAccountToMockReturnBody, v.GetAccountToMockReturnError)
		
		v.GetAccountFromMockReturnBody.Balance = v.GetAccountFromMockReturnBody.Balance - v.Body.Amount
		v.GetAccountToMockReturnBody.Balance = v.GetAccountToMockReturnBody.Balance + v.Body.Amount
		
		var mockCall = s.transactionAccMock.On("CreateTransactionAccount", v.Body, v.GetAccountFromMockReturnBody, v.GetAccountToMockReturnBody).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.transactionAccService.CreateTransactionAccount(v.Body)
			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGetAccountFrom.Unset()
		mockCallGetAccountTo.Unset()
		mockCall.Unset()
	}
}
func TestSuiteTransactionAccs(t *testing.T) {
	suite.Run(t, new(suiteTransactionAcc))
}
