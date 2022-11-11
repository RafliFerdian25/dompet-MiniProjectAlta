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

	"github.com/stretchr/testify/suite"
)

type suiteTransaction struct {
	suite.Suite
	transactionService *transactionService
	transactionMock    *transactionMockRepository.TransactionMock
	accountMock           *accountMockRepository.AccountMock
	subCategoryMock    *subCategoryMockRepository.SubCategoryMock
}

func (s *suiteTransaction) SetupSuite() {
	transactionMock := &transactionMockRepository.TransactionMock{}
	s.transactionMock = transactionMock

	accountMock := &accountMockRepository.AccountMock{}
	s.accountMock = accountMock

	subCategoryMock := &subCategoryMockRepository.SubCategoryMock{}
	s.subCategoryMock = subCategoryMock

	s.transactionService = &transactionService{
		transactionRepo: s.transactionMock,
		accountRepo:  s.accountMock,
		subCategoryRepo: s.subCategoryMock,
	}
}


func (s *suiteTransaction) TestGetAccountByUser() {
	testCase := []struct {
		Name            string
		GetExpenseMockReturnError error
		GetExpenseMockReturnBody  []dto.GetTransactionDTO
		GetIncomeMockReturnError error
		GetIncomeMockReturnBody  []dto.GetTransactionDTO
		ParamUserId     uint
		ParamMonth      int
		HasReturnBody   bool
		ExpectedBody     map[string]interface{}
		ExpectError     error
	}{
		{
			"success",
			nil,
			[]dto.GetTransactionDTO{
				{
					ID: 1,
					UserID: 1,
					SubCategoryID: 8,
					CategoryID: 2,
					AccountID: 1,
					Amount: 10000,
					Note: "test",
				},
			},
			nil,
			[]dto.GetTransactionDTO{
				{
					ID: 1,
					UserID: 1,
					SubCategoryID: 8,
					CategoryID: 2,
					AccountID: 1,
					Amount: 10000,
					Note: "test",
				},
			},
			1,
			11,
			true,
			map[string]interface{}{
				"expense": []dto.GetTransactionDTO{
					{
						ID: 1,
						UserID: 1,
						SubCategoryID: 8,
						CategoryID: 2,
						AccountID: 1,
						Amount: 10000,
						Note: "test",
					},
				},
				"income": []dto.GetTransactionDTO{
					{
						ID: 1,
						UserID: 1,
						SubCategoryID: 8,
						CategoryID: 2,
						AccountID: 1,
						Amount: 10000,
						Note: "test",
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
					ID: 1,
					UserID: 1,
					SubCategoryID: 8,
					CategoryID: 2,
					AccountID: 1,
					Amount: 10000,
					Note: "test",
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
					ID: 1,
					UserID: 1,
					SubCategoryID: 8,
					CategoryID: 2,
					AccountID: 1,
					Amount: 10000,
					Note: "test",
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
		GetAccountMockReturnError error
		GetAccountMockReturnBody  dto.AccountDTO
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
				CategoryID:   2,
				AccountID:   2,
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
		// {
		// 	"failed get account",
		// 	nil,
		// 	dto.TransactionAccount{
		// 		ID:            1,
		// 		UserID:        1,
		// 		AccountFromID: 1,
		// 		AccountToID:   2,
		// 		Amount:        10000,
		// 		Note:          "test",
		// 		AdminFee:      0,
		// 	},
		// 	errors.New(constantError.ErrorAccountNotFound),
		// 	dto.AccountDTO{},
		// 	nil,
		// 	1,
		// 	1,
		// 	true,
		// 	errors.New(constantError.ErrorAccountNotFound),
		// },
		// {
		// 	"failed get account",
		// 	nil,
		// 	dto.TransactionAccount{
		// 		ID:            1,
		// 		UserID:        1,
		// 		AccountFromID: 1,
		// 		AccountToID:   2,
		// 		Amount:        10000,
		// 		Note:          "test",
		// 		AdminFee:      0,
		// 	},
		// 	nil,
		// 	dto.AccountDTO{
		// 		ID:      1,
		// 		UserID:  1,
		// 		Name:    "test",
		// 		Balance: 10000,
		// 	},
		// 	nil,
		// 	1,
		// 	1,
		// 	true,
		// 	errors.New(constantError.ErrorAccountNotFound),
		// },
		// {
		// 	"error auth",
		// 	nil,
		// 	dto.TransactionAccount{
		// 		ID:            1,
		// 		UserID:        1,
		// 		AccountFromID: 1,
		// 		AccountToID:   2,
		// 		Amount:        10000,
		// 		Note:          "test",
		// 		AdminFee:      0,
		// 	},
		// 	nil,
		// 	dto.AccountDTO{},
		// 	nil,
		// 	1,
		// 	2,
		// 	true,
		// 	errors.New(constantError.ErrorNotAuthorized),
		// },
		// {
		// 	"account To not enough balance",
		// 	nil,
		// 	dto.TransactionAccount{
		// 		ID:            1,
		// 		UserID:        1,
		// 		AccountFromID: 1,
		// 		AccountToID:   2,
		// 		Amount:        100000,
		// 		Note:          "test",
		// 		AdminFee:      0,
		// 	},
		// 	nil,
		// 	dto.AccountDTO{
		// 		ID:      1,
		// 		UserID:  1,
		// 		Name:    "test",
		// 		Balance: 10000,
		// 	},
		// 	nil,
		// 	1,
		// 	1,
		// 	true,
		// 	errors.New(constantError.ErrorRecipientAccountNotEnoughBalance),
		// },
		// {
		// 	"failed delete transaction account",
		// 	nil,
		// 	dto.TransactionAccount{
		// 		ID:            1,
		// 		UserID:        1,
		// 		AccountFromID: 1,
		// 		AccountToID:   2,
		// 		Amount:        10000,
		// 		Note:          "test",
		// 		AdminFee:      0,
		// 	},
		// 	nil,
		// 	dto.AccountDTO{
		// 		ID:      1,
		// 		UserID:  1,
		// 		Name:    "test",
		// 		Balance: 10000,
		// 	},
		// 	gorm.ErrRecordNotFound,
		// 	1,
		// 	1,
		// 	true,
		// 	gorm.ErrRecordNotFound,
		// },
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

func TestSuiteTransactions(t *testing.T) {
	suite.Run(t, new(suiteTransaction))
}
