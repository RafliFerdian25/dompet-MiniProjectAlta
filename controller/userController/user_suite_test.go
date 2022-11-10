package userController

import (
	"bytes"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	userMockService "dompet-miniprojectalta/service/userService/userMock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type suiteUsers struct {
	suite.Suite
	userController *UserController
	mock           *userMockService.UserMock
}

func (s *suiteUsers) SetupSuite() {
	mock := &userMockService.UserMock{}
	s.mock = mock
	s.userController = &UserController{
		UserService: s.mock,
	}
}

func (s *suiteUsers) TestCreateUsers() {
	testCase := []struct {
		Name               string
		Method             string
		Body               dto.UserDTO
		mockParam dto.UserDTO
		MockReturnError	error
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success create user",
			"POST",
			dto.UserDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			dto.UserDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			nil,
			http.StatusOK,
			"success create user",
		},
		{
			"fail bind data",
			"POST",
			dto.UserDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			dto.UserDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			nil,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"There is an empty field",
			"POST",
			dto.UserDTO{
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			dto.UserDTO{
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			nil,
			http.StatusBadRequest,
			"There is an empty field",
		},
		{
			"failCreate",
			"POST",
			dto.UserDTO{
				Name: "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			dto.UserDTO{
				Name: "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			errors.New("error"),
			http.StatusInternalServerError,
			"fail create user",
		},
	}

	for i, v := range testCase {
		var mockCall = s.mock.On("CreateUser", v.mockParam).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			r := httptest.NewRequest(v.Method, "/signup", bytes.NewBuffer(res))
			if i != 1 {
				r.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			e.Validator = &helper.CustomValidator{
				Validator: validator.New(),
			}
			ctx := e.NewContext(r, w)

			err := s.userController.CreateUser(ctx)
			s.NoError(err)

			s.Equal(v.ExpectedStatusCode, w.Result().StatusCode)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)

			s.NoError(err)
			s.Equal(v.ExpectedMesaage, resp["message"].(string))
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteUsers) TestLoginUsers() {
	testCase := []struct {
		Name               string
		ExpectedStatusCode int
		Method             string
		MockReturnBody    model.User
		MockReturnError   error
		Body               model.User
		HasReturnBody      bool
		ExpectedBody       dto.UserResponseDTO
		ExpectedMessage    string
	}{
		{
			"success",
			http.StatusOK,
			"GET",
			model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			nil,
			model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			true,
			dto.UserResponseDTO{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Token: "123456",
			},
			"success login",
		},
		{
			"fail bind data",
			http.StatusInternalServerError,
			"GET",
			model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			nil,
			model.User{},
			false,
			dto.UserResponseDTO{},
			"fail bind data",
		},
		{
			"fail login",
			http.StatusInternalServerError,
			"GET",
			model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			errors.New("error"),
			model.User{},
			false,
			dto.UserResponseDTO{},
			"fail login",
		},
	}

	for i, v := range testCase {
		var mockCall = s.mock.On("LoginUser", v.Body).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			r := httptest.NewRequest(v.Method, "/login", bytes.NewBuffer(res))
			if i != 1 {
				r.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			
			err := s.userController.LoginUser(ctx)
			s.NoError(err)

			s.Equal(v.ExpectedStatusCode, w.Result().StatusCode)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)
			s.Equal(v.ExpectedMessage, resp["message"].(string))

			if v.HasReturnBody {
				s.Equal(v.ExpectedBody.Name, resp["user"].(map[string]interface{})["name"])
				s.Equal(v.ExpectedBody.Email, resp["user"].(map[string]interface{})["email"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}


func TestSuiteUsers(t *testing.T) {
	suite.Run(t, new(suiteUsers))
}
