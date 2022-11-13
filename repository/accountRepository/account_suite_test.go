package accountRepository

import (
	"database/sql/driver"
	"dompet-miniprojectalta/models/dto"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type suiteAccounts struct {
	suite.Suite
	accountRepository AccountRepository
	mock           sqlmock.Sqlmock
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *suiteAccounts) SetupSuite() {
	dbGormPalsu, mocking, err := sqlmock.New()
	s.NoError(err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbGormPalsu,
	}))

	s.mock = mocking

	NewAccountRepository := NewAccountRepository(dbGorm)
	s.accountRepository = NewAccountRepository
}


func (s *suiteAccounts) TestDeleteAccount() {
	testCase := []struct {
		Name                  string
		MockReturnError error
		ParamId               uint
		HasReturnError        bool
		ExpectedError         error
		RowAffected int64
	}{
		{
			"success",
			nil,
			1,
			false,
			nil,
			1,
		},
		{
			"failed delete account",
			errors.New("error"),
			1,
			true,
			errors.New("error"),
			1,
		},
		{
			"no row affected",
			nil,
			1,
			true,
			gorm.ErrRecordNotFound,
			0,
		},
	}

	for _, v := range testCase {
		s.mock.ExpectBegin()
		if v.MockReturnError != nil {
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `deleted_at`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{},1).WillReturnError(v.MockReturnError)
				s.mock.ExpectRollback()
		} else {
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `deleted_at`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{},1).WillReturnResult(sqlmock.NewResult(1, v.RowAffected))
				s.mock.ExpectCommit()
		}
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.accountRepository.DeleteAccount(v.ParamId)

			if v.HasReturnError {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *suiteAccounts) TestGetAccountByUser() {
	testCase := []struct {
		Name            string
		MockReturnError error
		ParamUserId            uint
		MockReturnBody  *sqlmock.Rows
		HasReturnBody   bool
		ExpectedBody    []dto.AccountDTO
	}{
		{
			"success",
			nil,
			1,
			sqlmock.NewRows([]string{"id", "name", "user_id", "balance"}).
				AddRow(1, "BRI", 1, 1000000),
			true,
			[]dto.AccountDTO{
				{
					ID:       1,
					Name:     "BRI",
					UserID:  1,
					Balance:  1000000,
				},
			},
		},
		{
			"fail get account",
			errors.New("error"),
			1,
			sqlmock.NewRows([]string{"id", "name", "user_id", "balance"}).
				AddRow(1, "BRI", 1, 1000000),
			false,
			[]dto.AccountDTO{},
		},
	}

	for _, v := range testCase {
		if v.MockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `accounts`.`id`,`accounts`.`user_id`,`accounts`.`name`,`accounts`.`balance` FROM `accounts` WHERE (user_id = ? OR user_id IS NULL) AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnError(v.MockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `accounts`.`id`,`accounts`.`user_id`,`accounts`.`name`,`accounts`.`balance` FROM `accounts` WHERE (user_id = ? OR user_id IS NULL) AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(1).WillReturnRows(v.MockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.accountRepository.GetAccountByUser(v.ParamUserId)
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

func (s *suiteAccounts) TestGetAccountById() {
	testCase := []struct {
		Name            string
		MockReturnError error
		ParamId            uint
		MockReturnBody  *sqlmock.Rows
		HasReturnBody   bool
		ExpectedBody    dto.AccountDTO
	}{
		{
			"success",
			nil,
			1,
			sqlmock.NewRows([]string{"id", "name", "user_id", "balance"}).
				AddRow(1, "BRI", 1, 1000000),
			true,
			dto.AccountDTO{
					ID:       1,
					Name:     "BRI",
					UserID:  1,
					Balance:  1000000,
			},
		},
		{
			"fail get account",
			errors.New("error"),
			1,
			sqlmock.NewRows([]string{"id", "name", "user_id", "balance"}).
				AddRow(1, "BRI", 1, 1000000),
			false,
			dto.AccountDTO{},
		},
		{
			"account not found",
			gorm.ErrRecordNotFound,
			1,
			sqlmock.NewRows([]string{"id", "name", "user_id", "balance"}).
				AddRow(1, "BRI", 1, 1000000),
			false,
			dto.AccountDTO{},
		},
	}

	for _, v := range testCase {
		if v.MockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `accounts`.`id`,`accounts`.`user_id`,`accounts`.`name`,`accounts`.`balance` FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`deleted_at` IS NULL ORDER BY `accounts`.`id` LIMIT 1")).
				WithArgs(1).WillReturnError(v.MockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT `accounts`.`id`,`accounts`.`user_id`,`accounts`.`name`,`accounts`.`balance` FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`deleted_at` IS NULL ORDER BY `accounts`.`id` LIMIT 1")).
				WithArgs(1).WillReturnRows(v.MockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.accountRepository.GetAccountById(v.ParamId)
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

func (s *suiteAccounts) TestUpdateAccount() {
	testCase := []struct {
		Name            string
		Body            dto.AccountDTO
		MockReturnError error
		HasReturnError  bool
		ExpectedError   error
		RowAffected int64
	}{
		{
			"success",
			dto.AccountDTO{
				ID:       1,
				Name:     "BRI",
			},
			nil,
			false,
			nil,
			1,
		},
		{
			"failed update user",
			dto.AccountDTO{
				ID:       1,
				Name:     "BRI",
			},
			errors.New("error"),
			true,
			errors.New("error"),
			1,
		},
		{
			"no row affected",
			dto.AccountDTO{
				ID:       1,
				Name:     "BRI",
			},
			nil,
			true,
			gorm.ErrRecordNotFound,
			0,
		},
	}

	for _, v := range testCase {
		s.mock.ExpectBegin()
		if v.MockReturnError != nil {
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `updated_at`=?,`name`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{},"BRI", 1).WillReturnError(v.MockReturnError)
				s.mock.ExpectRollback()
		} else {
			s.mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `updated_at`=?,`name`=? WHERE id = ? AND `accounts`.`deleted_at` IS NULL")).
				WithArgs(AnyTime{}, "BRI", 1).WillReturnResult(sqlmock.NewResult(1, v.RowAffected))
				s.mock.ExpectCommit()
		}
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.accountRepository.UpdateAccount(v.Body)
			if v.HasReturnError {
				s.Equal(v.ExpectedError, err)
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *suiteAccounts) TestCreateAccount() {
	testCase := []struct {
		Name            string
		Body            dto.AccountDTO
		MockReturnError error
		HasReturnError  bool
		ExpectedError   error
	}{
		{
			"success",
			dto.AccountDTO{
				UserID: 1,
				Name:     "BRI",
				Balance: 10000,
			},
			nil,
			false,
			nil,
		},
		{
			"failed create user",
			dto.AccountDTO{
				UserID: 1,
				Name:     "BRI",
				Balance: 10000,
			},
			errors.New("error"),
			true,
			errors.New("error"),
		},
	}

	for _, v := range testCase {
		s.mock.ExpectBegin()
		if v.MockReturnError != nil {
			s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`name`,`balance`) VALUES (?,?,?,?,?,?)")).
				WithArgs(AnyTime{}, AnyTime{}, nil,1,"BRI", float64(10000)).WillReturnError(v.MockReturnError)
				s.mock.ExpectRollback()
		} else {
			s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`created_at`,`updated_at`,`deleted_at`,`user_id`,`name`,`balance`) VALUES (?,?,?,?,?,?)")).
				WithArgs(AnyTime{}, AnyTime{}, nil,1, "BRI", float64(10000)).WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
		}
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.accountRepository.CreateAccount(v.Body)
			if v.HasReturnError {
				s.Equal(v.ExpectedError, err)
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *suiteAccounts) TearDownSuite() {
	s.mock = nil
}

func TestSuiteAccounts(t *testing.T) {
	suite.Run(t, new(suiteAccounts))
}
