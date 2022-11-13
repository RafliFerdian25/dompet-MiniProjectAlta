package transactionService

import (
	"dompet-miniprojectalta/constant/constantCategory"
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	accountMockRepository "dompet-miniprojectalta/repository/accountRepository/accountMock"
	subCategoryMockRepository "dompet-miniprojectalta/repository/subCategoryRepository/subCategoryMock"
	transactionMockRepository "dompet-miniprojectalta/repository/transactionRepository/transactionMock"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteTransaction struct {
	suite.Suite
	transactionService TransactionService
	transactionMock    *transactionMockRepository.TransactionMock
	accountMock        *accountMockRepository.AccountMock
	subCategoryMock    *subCategoryMockRepository.SubCategoryMock
}

func (s *suiteTransaction) SetupSuite() {
	transactionMock := &transactionMockRepository.TransactionMock{}
	s.transactionMock = transactionMock

	accountMock := &accountMockRepository.AccountMock{}
	s.accountMock = accountMock

	subCategoryMock := &subCategoryMockRepository.SubCategoryMock{}
	s.subCategoryMock = subCategoryMock

	// s.transactionService = &transactionService{
	// 	transactionRepo: s.transactionMock,
	// 	accountRepo:     s.accountMock,
	// 	subCategoryRepo: s.subCategoryMock,
	// }
	NewTransactionService := NewTransactionService(s.transactionMock, s.accountMock, s.subCategoryMock)
	s.transactionService = NewTransactionService

}

func (s *suiteTransaction) TestGetAccountByUser() {
	testCase := []struct {
		Name                      string
		GetExpenseMockReturnError error
		GetExpenseMockReturnBody  []dto.GetTransactionDTO
		GetIncomeMockReturnError  error
		GetIncomeMockReturnBody   []dto.GetTransactionDTO
		ParamUserId               uint
		ParamMonth                int
		HasReturnBody             bool
		ExpectedBody              map[string]interface{}
		ExpectError               error
	}{
		{
			"success",
			nil,
			[]dto.GetTransactionDTO{
				{
					ID:            1,
					UserID:        1,
					SubCategoryID: 8,
					CategoryID:    2,
					AccountID:     1,
					Amount:        10000,
					Note:          "test",
				},
			},
			nil,
			[]dto.GetTransactionDTO{
				{
					ID:            1,
					UserID:        1,
					SubCategoryID: 8,
					CategoryID:    2,
					AccountID:     1,
					Amount:        10000,
					Note:          "test",
				},
			},
			1,
			11,
			true,
			map[string]interface{}{
				"expense": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						AccountID:     1,
						Amount:        10000,
						Note:          "test",
					},
				},
				"income": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						AccountID:     1,
						Amount:        10000,
						Note:          "test",
					},
				},
				"month_transaction": time.Month(11).String(),
			},
			nil,
		},
		{
			"failed get expense",
			errors.New(constantError.ErrorTransactionNotFound),
			[]dto.GetTransactionDTO{},
			nil,
			[]dto.GetTransactionDTO{
				{
					ID:            1,
					UserID:        1,
					SubCategoryID: 8,
					CategoryID:    2,
					AccountID:     1,
					Amount:        10000,
					Note:          "test",
				},
			},
			1,
			11,
			false,
			map[string]interface{}{},
			errors.New(constantError.ErrorTransactionNotFound),
		},
		{
			"failed get income",
			nil,
			[]dto.GetTransactionDTO{},
			errors.New(constantError.ErrorTransactionNotFound),
			[]dto.GetTransactionDTO{
				{
					ID:            1,
					UserID:        1,
					SubCategoryID: 8,
					CategoryID:    2,
					AccountID:     1,
					Amount:        10000,
					Note:          "test",
				},
			},
			1,
			11,
			false,
			map[string]interface{}{},
			errors.New(constantError.ErrorTransactionNotFound),
		},
	}

	for _, v := range testCase {
		var mockCallGetExpense = s.transactionMock.On("GetTransaction", v.ParamUserId, uint(constantCategory.ExpenseCategory), v.ParamMonth).Return(v.GetExpenseMockReturnBody, v.GetExpenseMockReturnError)
		var mockCallGetIncome = s.transactionMock.On("GetTransaction", v.ParamUserId, uint(constantCategory.IncomeCategory), v.ParamMonth).Return(v.GetIncomeMockReturnBody, v.GetIncomeMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			transactionAccounts, err := s.transactionService.GetTransaction(v.ParamUserId, v.ParamMonth)

			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, transactionAccounts)
			} else {
				s.Error(err)
			}
		})
		// remove mock
		mockCallGetExpense.Unset()
		mockCallGetIncome.Unset()
	}
}

func (s *suiteTransaction) TestDeleteAccount() {
	testCase := []struct {
		Name                          string
		GetTransactionMockReturnError error
		GetTransactionMockReturnBody  dto.TransactionJoin
		GetAccountMockReturnError     error
		GetAccountMockReturnBody      dto.AccountDTO
		DeleteMockReturnError         error
		ParamId                       uint
		ParamUserId                   uint
		HasReturnError                bool
		ExpectedError                 error
	}{
		{
			"success",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     2,
				Amount:        1000,
			},
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
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     2,
				Amount:        1000,
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
			"failed get transaction",
			errors.New(constantError.ErrorAccountNotFound),
			dto.TransactionJoin{},
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
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     2,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{},
			nil,
			1,
			2,
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"failed delete transaction account",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     2,
				Amount:        1000,
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
		var mockCallGetTransaction = s.transactionMock.On("GetTransactionById", v.ParamId).Return(v.GetTransactionMockReturnBody, v.GetTransactionMockReturnError)
		var mockCallGetAccount = s.accountMock.On("GetAccountById", v.GetTransactionMockReturnBody.AccountID).Return(v.GetAccountMockReturnBody, v.GetAccountMockReturnError)

		v.GetAccountMockReturnBody.Balance = v.GetAccountMockReturnBody.Balance - v.GetTransactionMockReturnBody.Amount
		var mockCallDelete = s.transactionMock.On("DeleteTransaction", v.ParamId, v.GetAccountMockReturnBody).Return(v.DeleteMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.transactionService.DeleteTransaction(v.ParamId, v.ParamUserId)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGetTransaction.Unset()
		mockCallGetAccount.Unset()
		mockCallDelete.Unset()
	}
}

func (s *suiteTransaction) TestUpdateTransaction() {
	userid := uint(1)
	testCase := []struct {
		Name                             string
		GetOldTransactionMockReturnError error
		GetOldTransactionMockReturnBody  dto.TransactionJoin
		GetNewAccountMockReturnError     error
		GetNewAccountMockReturnBody      dto.AccountDTO
		ParamNewAccountMock              float64
		GetOldAccountMockReturnError     error
		GetOldAccountMockReturnBody      dto.AccountDTO
		ParamOldAccountMock              float64
		GetSubCategoryMockReturnError    error
		GetSubCategoryMockReturnBody     dto.SubCategoryDTO
		UpdateMockReturnError            error
		ParamUpdateMockBody              dto.TransactionDTO
		Body                             dto.TransactionDTO
		ParamUserId                      uint
		HasReturnError                   bool
		ExpectedError                    error
	}{
		{
			"success not change account",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			7000,
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			9000,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			1,
			false,
			nil,
		},
		{
			"failed get old transaction",
			errors.New(constantError.ErrorTransactionNotFound),
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.SubCategoryDTO{},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorTransactionNotFound),
		},
		{
			"error auth old transaction",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.SubCategoryDTO{},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			2,
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"failed get sub category ",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.AccountDTO{},
			0,
			errors.New(constantError.ErrorSubCategoryNotFound),
			dto.SubCategoryDTO{},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorSubCategoryNotFound),
		},
		{
			"error auth sub category",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        2,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			2,
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"cannot change sub category",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.AccountDTO{},
			0,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 3,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorCannotChangeSubCategory),
		},
		{
			"failed get old account",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{},
			0,
			errors.New(constantError.ErrorAccountNotFound),
			dto.AccountDTO{},
			0,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"failed get new account",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			errors.New(constantError.ErrorAccountNotFound),
			dto.AccountDTO{},
			0,
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			10000,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     2,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     2,
				Amount:        2000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth when change account",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  2,
				Name:    "Mandiri",
				Balance: 10000,
			},
			0,
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			10000,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     2,
				Amount:        -2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     2,
				Amount:        2000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"old account balance not enough",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "Mandiri",
				Balance: 10000,
			},
			0,
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			10000,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -20000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        20000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorOldAccountBalanceNotEnough),
		},
		{
			"old account balance not enough",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "Mandiri",
				Balance: 10000,
			},
			0,
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			10000,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     2,
				Amount:        -20000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     2,
				Amount:        20000,
				Note:          "test",
			},
			1,
			true,
			errors.New(constantError.ErrorNewAccountBalanceNotEnough),
		},
		{
			"failed update transaction",
			nil,
			dto.TransactionJoin{
				ID:            1,
				UserID:        1,
				SubCategoryID: 8,
				CategoryID:    2,
				AccountID:     1,
				Amount:        1000,
			},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "Mandiri",
				Balance: 10000,
			},
			8000,
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			9000,
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			errors.New("failed update transaction"),
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -1000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        1000,
				Note:          "test",
			},
			1,
			true,
			errors.New("failed update transaction"),
		},
		// 	errors.New(constantError.ErrorAccountNotFound),
		// 	errors.New(constantError.ErrorNotAuthorized),
		// 	errors.New(constantError.ErrorBalanceLessThanZero),
		// 	gorm.ErrRecordNotFound,
	}

	for _, v := range testCase {
		var mockCallGetOldTransaction = s.transactionMock.On("GetTransactionById", v.Body.ID).Return(v.GetOldTransactionMockReturnBody, v.GetOldTransactionMockReturnError)
		var mockCallGetSubCategory = s.subCategoryMock.On("GetSubCategoryById", v.Body.SubCategoryID).Return(v.GetSubCategoryMockReturnBody, v.GetSubCategoryMockReturnError)
		var mockCallGetNewAccount *mock.Call
		if v.Body.AccountID != 0 && v.Body.AccountID != v.GetOldTransactionMockReturnBody.AccountID {
			mockCallGetNewAccount = s.accountMock.On("GetAccountById", v.Body.AccountID).Return(v.GetNewAccountMockReturnBody, v.GetNewAccountMockReturnError)
		}
		var mockCallGetOldAccount = s.accountMock.On("GetAccountById", v.GetOldTransactionMockReturnBody.AccountID).Return(v.GetOldAccountMockReturnBody, v.GetOldAccountMockReturnError)
		var mockCallUpdate = s.transactionMock.On("UpdateTransaction", v.ParamUpdateMockBody, v.GetOldAccountMockReturnBody.ID, v.ParamNewAccountMock, v.ParamOldAccountMock).Return(v.UpdateMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.transactionService.UpdateTransaction(v.Body, v.ParamUserId)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGetOldTransaction.Unset()
		mockCallGetSubCategory.Unset()
		if v.Body.AccountID != 0 && v.Body.AccountID != v.GetOldTransactionMockReturnBody.AccountID {
			mockCallGetNewAccount.Unset()
		}
		mockCallGetOldAccount.Unset()
		mockCallUpdate.Unset()
	}
}

func (s *suiteTransaction) TestCreateTransaction() {
	userid := uint(1)
	testCase := []struct {
		Name                          string
		GetSubCategoryMockReturnError error
		GetSubCategoryMockReturnBody  dto.SubCategoryDTO
		GetAccountMockReturnError     error
		GetAccountMockReturnBody      dto.AccountDTO
		CreateMockReturnError         error
		Body                          dto.TransactionDTO
		ParamCreateMockBody                          dto.TransactionDTO
		HasReturnError                bool
		ExpectedError                 error
	}{
		{
			"success",
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.AccountDTO{
				UserID:  1,
				Name:    "Mandiri",
				Balance: 10000,
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			false,
			nil,
		},
		{
			"failed get sub category",
			errors.New(constantError.ErrorSubCategoryNotFound),
			dto.SubCategoryDTO{},
			nil,
			dto.AccountDTO{},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			true,
			errors.New(constantError.ErrorSubCategoryNotFound),
		},
		{
			"error auth sub category",
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.AccountDTO{},
			nil,
			dto.TransactionDTO{
				UserID:        2,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        2,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"failed get account",
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			errors.New(constantError.ErrorAccountNotFound),
			dto.AccountDTO{},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			true,
			errors.New(constantError.ErrorAccountNotFound),
		},
		{
			"error auth account",
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.AccountDTO{
				UserID:  2,
				Name:    "Mandiri",
				Balance: 10000,
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"category id is 1 or category is debt",
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 1,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.AccountDTO{
				UserID:  1,
				Name:    "Mandiri",
				Balance: 10000,
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			true,
			errors.New(constantError.ErrorCannotUseCategory),
		},
		{
			"account not enough balance",
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.AccountDTO{
				UserID:  1,
				Name:    "Mandiri",
				Balance: 10000,
			},
			nil,
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        20000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -20000,
				Note:          "test",
			},
			true,
			errors.New(constantError.ErrorAccountNotEnoughBalance),
		},
		{
			"failed create account",
			nil,
			dto.SubCategoryDTO{
				ID:         7,
				CategoryID: 2,
				UserID:     &userid,
				Name:       "Entertainment",
			},
			nil,
			dto.AccountDTO{
				UserID:  1,
				Name:    "Mandiri",
				Balance: 10000,
			},
			errors.New("error"),
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        2000,
				Note:          "test",
			},
			dto.TransactionDTO{
				UserID:        1,
				SubCategoryID: 8,
				AccountID:     1,
				Amount:        -2000,
				Note:          "test",
			},
			true,
			errors.New("error"),
		},
	}

	for _, v := range testCase {
		var mockCallGetSubCategory = s.subCategoryMock.On("GetSubCategoryById", v.Body.SubCategoryID).Return(v.GetSubCategoryMockReturnBody, v.GetSubCategoryMockReturnError)
		var mockCallGetNewAccount = s.accountMock.On("GetAccountById", v.Body.AccountID).Return(v.GetAccountMockReturnBody, v.GetAccountMockReturnError)
		var mockCallCreate = s.transactionMock.On("CreateTransaction", v.ParamCreateMockBody, v.GetSubCategoryMockReturnBody.CategoryID, v.GetAccountMockReturnBody).Return(v.CreateMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.transactionService.CreateTransaction(v.Body)
			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGetSubCategory.Unset()
		mockCallGetNewAccount.Unset()
		mockCallCreate.Unset()
	}
}

func TestSuiteTransactions(t *testing.T) {
	suite.Run(t, new(suiteTransaction))
}
