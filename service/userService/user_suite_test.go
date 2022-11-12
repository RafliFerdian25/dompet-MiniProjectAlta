package userService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	userMockRepository "dompet-miniprojectalta/repository/userRepository/userMock"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteUsers struct {
	suite.Suite
	userService UserService
	mock        *userMockRepository.UserMock
}

func (s *suiteUsers) SetupSuite() {
	mock := &userMockRepository.UserMock{}
	s.mock = mock
	NewUserService := NewUserService(s.mock)
	s.userService = NewUserService
}

func (s *suiteUsers) TestCreateUsers() {
	testCase := []struct {
		Name            string
		MockReturnError error
		Body            dto.UserDTO
		HasReturnError  bool
	}{
		{
			"success",
			nil,
			dto.UserDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			false,
		},
		{
			"fail create",
			errors.New("error"),
			dto.UserDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			true,
		},
		// {
		// 	"error hash password",
		// 	errors.New("error"),
		// 	dto.UserDTO{
		// 		Name:     "rafli",
		// 		Email:    "rafli@gmail.com",
		// 		Password: string(mock.AnythingOfType("int")),
		// 	},
		// 	true,
		// },
	}

	for _, v := range testCase {
		var mockCall = s.mock.On("CreateUser", mock.Anything).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			err := s.userService.CreateUser(v.Body)
			if v.HasReturnError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteUsers) TestLoginUsers() {
	testCase := []struct {
		Name            string
		MockReturnBody  model.User
		MockReturnError error
		Body            model.User
		HasReturnBody   bool
		ExpectedBody    model.User
	}{
		{
			"success",
			model.User{
				Model:    gorm.Model{ID: 1},
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			nil,
			model.User{
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			true,
			model.User{
				Model:    gorm.Model{ID: 1},
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
		},
		{
			"failed",
			model.User{},
			errors.New(constantError.ErrorEmailOrPasswordNotMatch),
			model.User{
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			false,
			model.User{},
		},
	}

	for _, v := range testCase {
		var mockCall = s.mock.On("LoginUser", v.Body).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			user, err := s.userService.LoginUser(v.Body)

			if v.HasReturnBody {
				s.Equal(v.ExpectedBody, user)
				s.NoError(err)
			} else {
				s.Error(err)
				s.Empty(user)
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func TestSuiteUsers(t *testing.T) {
	suite.Run(t, new(suiteUsers))
}
