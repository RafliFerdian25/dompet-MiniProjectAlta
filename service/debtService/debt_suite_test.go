package debtService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	accountMockRepository "dompet-miniprojectalta/repository/accountRepository/accountMock"
	debtMockRepository "dompet-miniprojectalta/repository/debtRepository/debtMock"
	subCategoryMockRepository "dompet-miniprojectalta/repository/subCategoryRepository/subCategoryMock"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteDebts struct {
	suite.Suite
	debtService     *debtService
	debtMock        *debtMockRepository.DebtMock
	accountMock     *accountMockRepository.AccountMock
	subCategoryMock *subCategoryMockRepository.SubCategoryMock
}

func (s *suiteDebts) SetupSuite() {
	debtMock := &debtMockRepository.DebtMock{}
	s.debtMock = debtMock

	accountMock := &accountMockRepository.AccountMock{}
	s.accountMock = accountMock

	subCategoryMock := &subCategoryMockRepository.SubCategoryMock{}
	s.subCategoryMock = subCategoryMock

	s.debtService = &debtService{
		debtRepo:        s.debtMock,
		accountRepo:     s.accountMock,
		subCategoryRepo: s.subCategoryMock,
	}
}

func (s *suiteDebts) TestGetDebt() {
	testCase := []struct {
		Name                string
		MockReturnErrorDebt error
		MockReturnErrorLoan error
		MockReturnBodyDebt  []dto.GetDebtTransactionResponse
		MockReturnBodyLoan  []dto.GetDebtTransactionResponse
		ParamUserId         uint
		ParamDebtStatus     string
		HasReturnBody       bool
		ExpectedBody        map[string][]dto.GetDebtTransactionResponse
	}{
		{
			"success",
			nil,
			nil,
			[]dto.GetDebtTransactionResponse{
				{
					ID:            1,
					Name:          "test debt",
					SubCategoryID: 1,
					AccountID:     1,
					Total:         10000,
					Remaining:     10000,
					Note:          "test",
					CreatedAt:     time.Now(),
					DebtStatus:    "unpaid",
				},
			},
			[]dto.GetDebtTransactionResponse{
				{
					ID:            2,
					Name:          "test loan",
					SubCategoryID: 1,
					AccountID:     1,
					Total:         -10000,
					Remaining:     -10000,
					Note:          "test",
					CreatedAt:     time.Now(),
					DebtStatus:    "unpaid",
				},
			},
			1,
			"unpaid",
			true,
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "test debt",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
				"loan": {
					{
						ID:            2,
						Name:          "test loan",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         -10000,
						Remaining:     -10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
		},
		{
			"failed get account",
			errors.New("error"),
			nil,
			[]dto.GetDebtTransactionResponse{},
			[]dto.GetDebtTransactionResponse{},
			1,
			"unpaid",
			false,
			map[string][]dto.GetDebtTransactionResponse{},
		},
		{
			"failed get account",
			nil,
			errors.New("error"),
			[]dto.GetDebtTransactionResponse{},
			[]dto.GetDebtTransactionResponse{},
			1,
			"unpaid",
			false,
			map[string][]dto.GetDebtTransactionResponse{},
		},
	}

	for _, v := range testCase {
		var mockCallDebt = s.debtMock.On("GetDebt", v.ParamUserId, 1, v.ParamDebtStatus).Return(v.MockReturnBodyDebt, v.MockReturnErrorDebt)
		var mockCallLoan = s.debtMock.On("GetDebt", v.ParamUserId, 3, v.ParamDebtStatus).Return(v.MockReturnBodyLoan, v.MockReturnErrorLoan)
		s.T().Run(v.Name, func(t *testing.T) {
			debts, err := s.debtService.GetDebt(v.ParamUserId, v.ParamDebtStatus)
			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, debts)
			} else {
				s.Error(err)
			}
		})
		// remove mock
		mockCallDebt.Unset()
		mockCallLoan.Unset()
	}
}

func (s *suiteDebts) TestDeleteDebt() {
	testCase := []struct {
		Name                      string
		GetDebtMockReturnError    error
		GetDebtMockReturnBody     dto.Debt
		GetAccountMockReturnError error
		GetAccountMockReturnBody  dto.AccountDTO
		DeleteMockReturnError     error
		ParamId                   uint
		ParamUserId               uint
		HasReturnError            bool
		ExpectedError             error
	}{
		{
			"success",
			nil,
			dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      10000,
				Remaining:  10000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},
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
			"failed get debt",
			errors.New("error"),
			dto.Debt{},
			nil,
			dto.AccountDTO{},
			nil,
			1,
			1,
			true,
			errors.New("error"),
		},
		{
			"error auth",
			nil,
			dto.Debt{
				ID:         1,
				Name:       "test debt",
				Total:      10000,
				Remaining:  10000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        2,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},
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
			true,
			errors.New(constantError.ErrorNotAuthorized),
		},
		{
			"failed get account",
			nil,
			dto.Debt{
				ID:         1,
				Name:       "test debt",
				Total:      10000,
				Remaining:  10000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},
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
			"account not enough balance",
			nil,
			dto.Debt{
				ID:         1,
				Name:       "test debt",
				Total:      100000,
				Remaining:  10000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},
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
			true,
			errors.New(constantError.ErrorAccountNotEnoughBalance),
		},
		{
			"failed delete debt",
			nil,
			dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      10000,
				Remaining:  10000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        10000,
					},
				},
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
		MockParamDelete := dto.AccountDTO{
			ID:      v.GetAccountMockReturnBody.ID,
			UserID:  v.GetAccountMockReturnBody.UserID,
			Name:    v.GetDebtMockReturnBody.Name,
			Balance: v.GetAccountMockReturnBody.Balance - v.GetDebtMockReturnBody.Total,
		}
		var mockCallGetDebt = s.debtMock.On("GetDebtById", v.ParamId).Return(v.GetDebtMockReturnBody, v.GetDebtMockReturnError)
		var mockCallGetAccount *mock.Call
		if v.GetDebtMockReturnError == nil {
			mockCallGetAccount = s.accountMock.On("GetAccountById", v.GetDebtMockReturnBody.Transactions[0].AccountID).Return(v.GetAccountMockReturnBody, v.GetAccountMockReturnError)
		}
		var mockCallDelete = s.debtMock.On("DeleteDebt", v.ParamId, MockParamDelete).Return(v.DeleteMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.debtService.DeleteDebt(v.ParamId, v.ParamUserId)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGetDebt.Unset()
		if v.GetDebtMockReturnError == nil {
			mockCallGetAccount.Unset()
		}
		mockCallDelete.Unset()
	}
}

func (s *suiteDebts) TestCreateDebt() {
	testCase := []struct {
		Name                      string
		GetDebtMockReturnError    error
		GetDebtMockReturnBody     dto.Debt
		GetAccountMockReturnError error
		GetAccountMockReturnBody  dto.AccountDTO
		CreateMockReturnError     error
		Body                      dto.DebtTransactionDTO
		HasReturnError            bool
		ExpectedError             error
	}{
		{
			"success debt without debt id",
			nil,
			dto.Debt{},
			nil,
			dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			nil,
			dto.DebtTransactionDTO{
				Name:   "test",
				UserID: 1,
				// DebtID: 1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			false,
			nil,
		},
		{
			Name:                      "success loan without debt id",
			GetDebtMockReturnError:    nil,
			GetDebtMockReturnBody:     dto.Debt{},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:   "test",
				UserID: 1,
				// DebtID: 1,
				SubCategoryID: 3,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: false,
			ExpectedError:  nil,
		},
		{
			Name:                   "success debt with debt id",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test debt",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        10000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        5000,
				Note:          "test",
			},
			HasReturnError: false,
			ExpectedError:  nil,
		},
		{
			Name:                   "success debt repayment with debt id",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test debt",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        10000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      2,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 2,
				AccountID:     1,
				Amount:        5000,
				Note:          "test",
			},
			HasReturnError: false,
			ExpectedError:  nil,
		},
		{
			Name:                      "failed get account",
			GetDebtMockReturnError:    nil,
			GetDebtMockReturnBody:     dto.Debt{},
			GetAccountMockReturnError: errors.New(constantError.ErrorAccountNotFound),
			GetAccountMockReturnBody:  dto.AccountDTO{},
			CreateMockReturnError:     nil,
			Body: dto.DebtTransactionDTO{
				Name:   "test",
				UserID: 1,
				// DebtID: 1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(constantError.ErrorAccountNotFound),
		},
		{
			Name:                      "error auth",
			GetDebtMockReturnError:    nil,
			GetDebtMockReturnBody:     dto.Debt{},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      2,
				UserID:  2,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:   "test",
				UserID: 1,
				// DebtID: 1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(constantError.ErrorNotAuthorized),
		},
		{
			Name:                      "account balance not enough",
			GetDebtMockReturnError:    nil,
			GetDebtMockReturnBody:     dto.Debt{},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:   "test",
				UserID: 1,
				// DebtID: 1,
				SubCategoryID: 2,
				AccountID:     1,
				Amount:        100000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(constantError.ErrorAccountNotEnoughBalance),
		},
		{
			Name:                      "failed get debt",
			GetDebtMockReturnError:    errors.New("error"),
			GetDebtMockReturnBody:     dto.Debt{},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody:  dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError:     nil,
			Body: dto.DebtTransactionDTO{
				Name:   "test",
				UserID: 1,
				DebtID: 1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New("error"),
		},
		{
			Name:                      "debt cannot change category",
			GetDebtMockReturnError:    nil,
			GetDebtMockReturnBody:     dto.Debt{
				ID:         1,
				Name:       "test debt",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody:  dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError:     nil,
			Body: dto.DebtTransactionDTO{
				Name:   "test",
				UserID: 1,
				DebtID: 1,
				SubCategoryID: 3,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(constantError.ErrorCannotChangeSubCategory),
		},
		{
			Name:                   "error auth debt",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        2,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(constantError.ErrorNotAuthorized),
		},
		{
			Name:                   "loan cannot change category",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "loan",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 3,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(constantError.ErrorCannotChangeSubCategory),
		},
		{
			Name:                   "error auth debt",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "loan",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 3,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(constantError.ErrorCannotChangeSubCategory),
		},
		{
			Name:                   "error amount greater than remaining debt",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        10000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 2,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(fmt.Sprint("Input amount is more than remaining debt. Unpaid amount is ", 5000)),
		},
		{
			Name:                   "error amount greater than remaining loan",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      -10000,
				Remaining:  -5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "loan",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 3,
						AccountID:     1,
						DebtID:        1,
						Amount:        10000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: nil,
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 4,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New(fmt.Sprint("Input amount is more than remaining loan. Unpaid amount is ", 5000)),
		},
		{
			Name:                   "failed create debt",
			GetDebtMockReturnError: nil,
			GetDebtMockReturnBody: dto.Debt{
				ID:         1,
				Name:       "test",
				Total:      10000,
				Remaining:  5000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        10000,
					},
				},
			},
			GetAccountMockReturnError: nil,
			GetAccountMockReturnBody: dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "test",
				Balance: 10000,
			},
			CreateMockReturnError: errors.New("error"),
			Body: dto.DebtTransactionDTO{
				Name:          "test",
				UserID:        1,
				DebtID:        1,
				SubCategoryID: 1,
				AccountID:     1,
				Amount:        10000,
				Note:          "test",
			},
			HasReturnError: true,
			ExpectedError:  errors.New("error"),
		},
	}

	for _, v := range testCase {
		var mockCallGetAccount = s.accountMock.On("GetAccountById", v.Body.AccountID).Return(v.GetAccountMockReturnBody, v.GetAccountMockReturnError)
		var mockCallGetDebt *mock.Call
		if v.Body.DebtID != 0 {
			mockCallGetDebt = s.debtMock.On("GetDebtById", v.Body.DebtID).Return(v.GetDebtMockReturnBody, v.GetDebtMockReturnError)
		}
		var mockCallCreate = s.debtMock.On("CreateDebt", mock.Anything, mock.Anything, mock.Anything).Return(v.CreateMockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.debtService.CreateDebt(v.Body)
			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCallGetAccount.Unset()
		if v.Body.DebtID != 0 {
			mockCallGetDebt.Unset()
		}
		mockCallCreate.Unset()
	}
}
func TestSuiteDebts(t *testing.T) {
	suite.Run(t, new(suiteDebts))
}
