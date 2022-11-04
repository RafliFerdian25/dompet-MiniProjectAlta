package userController

import (
	"bytes"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"dompet-miniprojectalta/service/userService/userMock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type suiteUsers struct {
	suite.Suite
	userController *UserController
	mock           *userMock.UserMock
}

func (s *suiteUsers) SetupSuite() {
	mock := &userMock.UserMock{}
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
		HasReturnBody      bool
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success",
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
			true,
			http.StatusOK,
			"success",
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
			false,
			http.StatusInternalServerError,
			"Failed",
		},
	}

	for _, v := range testCase {
		var mockCall = s.mock.On("CreateUser", v.mockParam)
		switch v.Name {
		case "success":
			mockCall.Return(nil)
		case "failCreate":
			mockCall.Return(errors.New("Failed"))
		}
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			r := httptest.NewRequest(v.Method, "/signup", bytes.NewBuffer(res))
			r.Header.Set("Content-Type", "application/json")
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

			if v.HasReturnBody {
				var resp map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&resp)

				s.NoError(err)
				s.Equal(v.ExpectedMesaage, resp["message"].(string))
			}
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
		Body               model.User
		HasReturnBody      bool
		ExpectedBody       model.User
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
			false,
			model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
		},
		{
			"failLogin",
			http.StatusInternalServerError,
			"GET",
			model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
			false,
			model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			},
		},
	}

	for _, v := range testCase {
		var mockAuth model.User
		switch v.Name {
		case "success":
			mockAuth = model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			}
		case "failLogin":
			mockAuth = model.User{
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			}
		}
		var mockCall = s.mock.On("LoginUser", mockAuth)
		switch v.Name {
		case "success":
			mockCall.Return(model.User{
				Model:   gorm.Model{
					ID: 1,
				},
				Name:     "rafli",
				Email:    "rafli@gmail.com",
				Password: "123456",
			}, nil)
		case "failLogin":
			mockCall.Return(model.User{}, errors.New("error"))
		}
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			r := httptest.NewRequest(v.Method, "/login", bytes.NewBuffer(res))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			err := s.userController.LoginUser(ctx)
			s.NoError(err)

			s.Equal(v.ExpectedStatusCode, w.Result().StatusCode)

			if v.HasReturnBody {
				var resp map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&resp)

				s.NoError(err)
				s.Equal(t, v.ExpectedBody.Name, resp["user"].(map[string]interface{})["name"])
				s.Equal(t, v.ExpectedBody.Email, resp["user"].(map[string]interface{})["email"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}


func TestSuiteUsers(t *testing.T) {
	suite.Run(t, new(suiteUsers))
}
