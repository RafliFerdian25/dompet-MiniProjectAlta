package userRepository

import (
	"database/sql/driver"
	"dompet-miniprojectalta/helper"
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

type suiteUsers struct {
	suite.Suite
	userRepository UserRepository
	mock           sqlmock.Sqlmock
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *suiteUsers) SetupSuite() {
	dbGormPalsu, mocking, err := sqlmock.New()
	s.NoError(err)

	dbGorm, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      dbGormPalsu,
	}))

	s.mock = mocking

	NewUserRepository := NewUserRepository(dbGorm)
	s.userRepository = NewUserRepository
}

func (s *suiteUsers) TestCreateUser() {
	password, _ := helper.HashPassword("123456")

	testCase := []struct {
		Name            string
		Body            dto.UserDTO
		MockReturnError error
		HasReturnError  bool
		ExpectedError   error
	}{
		{
			"success",
			dto.UserDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: password,
			},
			nil,
			false,
			nil,
		},
		{
			"fail get user",
			dto.UserDTO{
				Name:    "rafli",
				Email:    "rafli@gmail.com",
				Password: password,
			},
			errors.New("error"),
			true,
			errors.New("error"),
		},
	}

	for _, v := range testCase {
		s.mock.ExpectBegin()
		if v.MockReturnError != nil {
			s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`email`,`password`) VALUES (?,?,?,?,?,?)")).
				WithArgs(AnyTime{},AnyTime{},nil,"rafli","rafli@gmail.com",password).WillReturnError(v.MockReturnError)
				s.mock.ExpectRollback()
		} else {
			s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`email`,`password`) VALUES (?,?,?,?,?,?)")).
				WithArgs(AnyTime{},AnyTime{},nil,"rafli","rafli@gmail.com",password).WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
		}
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.userRepository.CreateUser(v.Body)
			if v.HasReturnError {
				s.Equal(v.ExpectedError, err)
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}
func (s *suiteUsers) TestLoginUser() {
	password, _ := helper.HashPassword("123456")
	passwordNotMatch, _ := helper.HashPassword("123123")

	testCase := []struct {
		Name            string
		MockReturnError error
		Body            model.User
		MockReturnBody  *sqlmock.Rows
		HasReturnBody   bool
		ExpectedBody    model.User
	}{
		{
			"success",
			nil,
			model.User{
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			sqlmock.NewRows([]string{"ID", "Name", "Email", "Password"}).
				AddRow(1, "rafli", "rafli@gmail.com", password),
			true,
			model.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: password,
			},
		},
		{
			"fail get user",
			errors.New("error"),
			model.User{
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			sqlmock.NewRows([]string{"ID", "Name", "Email", "Password"}).
				AddRow(1, "rafli", "rafli@gmail.com", password),
			false,
			model.User{},
		},
		{
			"password not match",
			nil,
			model.User{
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			sqlmock.NewRows([]string{"ID", "Name", "Email", "Password"}).
				AddRow(1, "rafli", "rafli@gmail.com", passwordNotMatch),
			false,
			model.User{},
		},
	}

	for _, v := range testCase {
		if v.MockReturnError != nil {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs("rafli@gmail.com").WillReturnError(v.MockReturnError)
		} else {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs("rafli@gmail.com").WillReturnRows(v.MockReturnBody)
		}
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.userRepository.LoginUser(v.Body)
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

func (s *suiteUsers) TearDownSuite() {
	s.mock = nil
}

func TestSuiteUsers(t *testing.T) {
	suite.Run(t, new(suiteUsers))
}
