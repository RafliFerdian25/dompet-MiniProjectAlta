package debtRepository

import (
	"database/sql/driver"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type suiteDebts struct {
	suite.Suite
	debtRepository DebtRepostory
	mock           sqlmock.Sqlmock
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *suiteDebts) SetupSuite() {
	dbGormPalsu, mocking, err := sqlmock.New()
	s.NoError(err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbGormPalsu,
	}))

	s.mock = mocking

	NewDebtRepository := NewDebtRepository(dbGorm)
	s.debtRepository = NewDebtRepository
}

func (s *suiteDebts) TestGetDebt() {
	testCase := []struct {
		Name                string
		ParamUserId         uint
		ParamSubCategory    int
		ParamDebtStatus     string
		DebtMockReturnError error
		DebtMockReturnBody  *sqlmock.Rows
		HasReturnBody       bool
		ExpectedBody        []dto.GetDebtTransactionResponse
	}{
		{
			"success",
			1,
			1,
			"unpaid",
			nil,
			sqlmock.NewRows([]string{"id", "name", "sub_category_id", "account_id", "total", "remaining", "note", "debts.created_at", "debt_status"}).
				AddRow(1, "rafli", 1, 1, 10000, 10000, "test", nil, "unpaid"),
			true,
			[]dto.GetDebtTransactionResponse{
				{
					ID:            1,
					Name:          "rafli",
					SubCategoryID: 1,
					AccountID:     1,
					Total:         10000,
					Remaining:     10000,
					Note:          "test",
					// CreatedAt:     time.Now(),
					DebtStatus: "unpaid",
				},
			},
		},
		{
			"failed get debt",
			1,
			1,
			"unpaid",
			errors.New("error"),
			sqlmock.NewRows([]string{"id", "name", "sub_category_id", "account_id", "total", "remaining", "note", "debts.created_at", "debt_status"}).
				AddRow(1, "rafli", 1, 1, 10000, 10000, "test", nil, "unpaid"),
			false,
			[]dto.GetDebtTransactionResponse{
				{
					ID:            1,
					Name:          "rafli",
					SubCategoryID: 1,
					AccountID:     1,
					Total:         10000,
					Remaining:     10000,
					Note:          "test",
					// CreatedAt:     time.Now(),
					DebtStatus: "unpaid",
				},
			},
		},
	}

	for _, v := range testCase {
		if v.DebtMockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT debts.id, debts.name, transactions.sub_category_id, transactions.account_id, debts.total, debts.remaining, debts.note, debts.created_at, debts.debt_status FROM `debts` JOIN transactions ON transactions.debt_id = debts.id WHERE (debts.debt_status = ? AND transactions.sub_category_id = ? AND transactions.user_id = ?) AND `debts`.`deleted_at` IS NULL GROUP BY `debts`.`id`")).
				WithArgs(v.ParamDebtStatus, v.ParamSubCategory, v.ParamUserId).WillReturnError(v.DebtMockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT debts.id, debts.name, transactions.sub_category_id, transactions.account_id, debts.total, debts.remaining, debts.note, debts.created_at, debts.debt_status FROM `debts` JOIN transactions ON transactions.debt_id = debts.id WHERE (debts.debt_status = ? AND transactions.sub_category_id = ? AND transactions.user_id = ?) AND `debts`.`deleted_at` IS NULL GROUP BY `debts`.`id`")).
				WithArgs(v.ParamDebtStatus, v.ParamSubCategory, v.ParamUserId).WillReturnRows(v.DebtMockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.debtRepository.GetDebt(v.ParamUserId, v.ParamSubCategory, v.ParamDebtStatus)
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody, user)
				s.NoError(err)
			} else {
				s.Error(err)
				s.Empty(user)
			}
		})
	}
}

func (s *suiteDebts) TestDeleteDebt() {
	testCase := []struct {
		Name                   string
		MockReturnError        error
		ParamId                uint
		ParamAccount           dto.AccountDTO
		HasReturnError         bool
		ExpectedError          error
		AccountRowAffected     int64
		TransactionRowAffected int64
		DebtRowAffected        int64
	}{
		{
			"success",
			nil,
			1,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			false,
			nil,
			1,
			1,
			1,
		},
		// {
		// 	"account row affected is 0",
		// 	nil,
		// 	1,
		// 	dto.AccountDTO{
		// 		ID:      1,
		// 		UserID:  1,
		// 		Name:    "BRI",
		// 		Balance: 10000,
		// 	},
		// 	true,
		// 	errors.New(constantError.ErrorAccountNotFound),
		// 	0,
		// 	1,
		// 	1,
		// },
		// {
		// 	"transaction row affected is 0",
		// 	nil,
		// 	1,
		// 	dto.AccountDTO{
		// 		ID:      1,
		// 		UserID:  1,
		// 		Name:    "BRI",
		// 		Balance: 10000,
		// 	},
		// 	true,
		// 	errors.New(constantError.ErrorTransactionNotFound),
		// 	1,
		// 	0,
		// 	1,
		// },
		// {
		// 	"failed delete account",
		// 	errors.New("error"),
		// 	1,
		// 	true,
		// 	errors.New("error"),
		// 	1,
		// },
		// {
		// 	"no row affected",
		// 	nil,
		// 	1,
		// 	true,
		// 	gorm.ErrRecordNotFound,
		// 	0,
		// },
	}

	for _, v := range testCase {
		s.mock.ExpectBegin()
		if v.MockReturnError != nil {
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `deleted_at`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{}, 1).WillReturnError(v.MockReturnError)
			s.mock.ExpectRollback()
		} else if v.AccountRowAffected == 0 || v.TransactionRowAffected == 0 || v.DebtRowAffected == 0 {
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `balance`=?,`updated_at`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(float64(10000), AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, v.AccountRowAffected))
			if v.AccountRowAffected == 0 {
				s.mock.ExpectRollback()
			}
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `transactions` SET `deleted_at`=? WHERE debt_id = ? AND `transactions`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, v.TransactionRowAffected))
			if v.AccountRowAffected == 0 {
				s.mock.ExpectRollback()
			}
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `debts` SET `deleted_at`=? WHERE id = ? AND `debts`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, v.DebtRowAffected))
		} else {
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `balance`=?,`updated_at`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(float64(10000), AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, v.AccountRowAffected))
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `transactions` SET `deleted_at`=? WHERE debt_id = ? AND `transactions`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, v.TransactionRowAffected))
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `debts` SET `deleted_at`=? WHERE id = ? AND `debts`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, v.DebtRowAffected))
			s.mock.ExpectCommit()
		}
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.debtRepository.DeleteDebt(v.ParamId, v.ParamAccount)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *suiteDebts) TestGetDebtById() {
	// userId := uint(1)
	testCase := []struct {
		Name                       string
		ParamId                    uint
		DebtMockReturnError        error
		DebtMockReturnBody         *sqlmock.Rows
		TransactionMockReturnError error
		TransactionMockReturnBody  *sqlmock.Rows
		HasReturnBody              bool
		ExpectedBody               dto.Debt
	}{
		{
			"success",
			1,
			nil,
			sqlmock.NewRows([]string{"id", "name", "total", "remaining", "note", "debt_status", "status", "transaction"}).
				AddRow(1, "rafli", 100000, 100000, "test", "test", "debt", nil),
			nil,
			sqlmock.NewRows([]string{"id", "user_id", "sub_category_id", "account_id", "debt_id", "amount", "note"}).
				AddRow(1, 1, 1, 1, 1, 100000, "test"),
			true,
			dto.Debt{
				ID:         1,
				Name:       "rafli",
				Total:      100000,
				Remaining:  100000,
				Note:       "test",
				DebtStatus: "test",
				Status:     "debt",
				Transactions: []model.Transaction{
					{
						Model: gorm.Model{
							ID: 1,
						},
						UserID:        1,
						SubCategoryID: 1,
						AccountID:     1,
						DebtID:        1,
						Amount:        100000,
						Note:          "test",
					},
				},
			},
		},
		{
			"failed get debt",
			1,
			errors.New("error"),
			sqlmock.NewRows([]string{"id", "name", "total", "remaining", "note", "debt_status", "status", "transaction"}).
				AddRow(1, "rafli", 100000, 100000, "test", "test", "debt", nil),
			errors.New("error"),
			sqlmock.NewRows([]string{"id", "user_id", "sub_category_id", "account_id", "debt_id", "amount", "note"}).
				AddRow(1, 1, 1, 1, 1, 100000, "test"),
			false,
			dto.Debt{},
		},
	}

	for _, v := range testCase {
		if v.DebtMockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `debts`.`id`,`debts`.`name`,`debts`.`total`,`debts`.`remaining`,`debts`.`note`,`debts`.`debt_status`,`debts`.`status` FROM `debts` WHERE `debts`.`id` = ? AND `debts`.`deleted_at` IS NULL ORDER BY `debts`.`id` LIMIT 1")).
				WithArgs(1).WillReturnError(v.DebtMockReturnError)
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE `transactions`.`debt_id` = ? AND `transactions`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnError(v.TransactionMockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `debts`.`id`,`debts`.`name`,`debts`.`total`,`debts`.`remaining`,`debts`.`note`,`debts`.`debt_status`,`debts`.`status` FROM `debts` WHERE `debts`.`id` = ? AND `debts`.`deleted_at` IS NULL ORDER BY `debts`.`id` LIMIT 1")).
				WithArgs(1).WillReturnRows(v.DebtMockReturnBody)
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `transactions` WHERE `transactions`.`debt_id` = ? AND `transactions`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnRows(v.TransactionMockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.debtRepository.GetDebtById(v.ParamId)
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody, user)
				s.NoError(err)
			} else {
				s.Error(err)
				s.Empty(user)
			}
		})
	}
}

func (s *suiteDebts) TestCreateDebt() {
	testCase := []struct {
		Name             string
		ParamDebt        dto.Debt
		ParamTransaction dto.TransactionDTO
		ParamAccount     dto.AccountDTO
		MockReturnError  error
		HasReturnError   bool
		ExpectedError    error
	}{
		{
			"success debt id is 0",
			dto.Debt{
				Name:       "rafli",
				Total:      10000,
				Remaining:  10000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
			},
			dto.TransactionDTO{
				ID:            1,
				UserID:        1,
				SubCategoryID: 1,
				AccountID:     1,
				DebtID:        0,
				Amount:        10000,
				Note:          "test",
			},
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			nil,
			false,
			nil,
		},
		{
			"success debt id not 0",
			dto.Debt{
				Name:       "rafli",
				Total:      10000,
				Remaining:  10000,
				Note:       "test",
				DebtStatus: "unpaid",
				Status:     "debt",
			},
			dto.TransactionDTO{
				ID:            1,
				UserID:        1,
				SubCategoryID: 1,
				AccountID:     1,
				DebtID:        1,
				Amount:        10000,
				Note:          "test",
			},
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BRI",
				Balance: 10000,
			},
			nil,
			false,
			nil,
		},
	}

	for _, v := range testCase {
		s.mock.ExpectBegin()
		if v.MockReturnError != nil {
			s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`name`,`balance`) VALUES (?,?,?,?,?,?)")).
				WithArgs(AnyTime{}, AnyTime{}, nil, 1, "BRI", float64(10000)).WillReturnError(v.MockReturnError)
			s.mock.ExpectRollback()
		} else {
			if v.ParamTransaction.DebtID == 0 {
				s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `debts` (`created_at`,`updated_at`,`deleted_at`,`name`,`total`,`remaining`,`note`,`debt_status`,`status`) VALUES (?,?,?,?,?,?,?,?,?)")).
					WithArgs(AnyTime{}, AnyTime{}, nil, "rafli", float64(10000), float64(10000), "test", "unpaid", "debt").WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`sub_category_id`,`account_id`,`amount`,`note`,`debt_id`) VALUES (?,?,?,?,?,?,?,?,?)")).
					WithArgs(AnyTime{}, AnyTime{}, nil, 1, 1, 1, float64(10000), "test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `debts` SET `updated_at`=?,`total`=?,`remaining`=?,`debt_status`=? WHERE id = ? AND `debts`.`deleted_at` IS NULL")).
					WithArgs(AnyTime{}, float64(10000), float64(10000), "unpaid", 1).WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `transactions` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`sub_category_id`,`account_id`,`amount`,`note`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(AnyTime{}, AnyTime{}, nil, 1, 1, 1, float64(10000), "test").WillReturnResult(sqlmock.NewResult(1, 1))
			}
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `balance`=?,`updated_at`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(float64(10000), AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()
		}
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.debtRepository.CreateDebt(v.ParamDebt, v.ParamTransaction, v.ParamAccount)
			if v.HasReturnError {
				s.Equal(v.ExpectedError, err)
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *suiteDebts) TearDownSuite() {
	s.mock = nil
}

func TestSuiteDebts(t *testing.T) {
	suite.Run(t, new(suiteDebts))
}
