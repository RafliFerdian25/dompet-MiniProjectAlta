package accountController

import (
	"bytes"
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	accountMockService "dompet-miniprojectalta/service/accountService/accountMock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type suiteAccount struct {
	suite.Suite
	accountController *AccountController
	mock              *accountMockService.AccountMock
}

func (s *suiteAccount) SetupSuite() {
	mock := &accountMockService.AccountMock{}
	s.mock = mock
	s.accountController = &AccountController{
		AccountService: s.mock,
	}
}

func (s *suiteAccount) TestDeleteAccount() {
	testCase := []struct {
		Name               string
		Method             string
		ParamID            string
		userId             uint
		MockReturnError    error
		HasReturnBody      bool
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success delete",
			"DELETE",
			"1",
			1,
			nil,
			true,
			http.StatusOK,
			"success delete account",
		},
		{
			"fail get id",
			"DELETE",
			"w",
			1,
			errors.New("error"),
			true,
			http.StatusBadRequest,
			"fail get id",
		},
		{
			"fail delete",
			"DELETE",
			"1",
			1,
			errors.New("error"),
			true,
			http.StatusInternalServerError,
			"fail delete account",
		},
		{
			"fail delete",
			"DELETE",
			"1",
			1,
			errors.New(constantError.ErrorAccountNotFound),
			true,
			http.StatusNotFound,
			"fail delete account",
		},
	}
	for _, v := range testCase {
		id, _ := strconv.Atoi(v.ParamID)
		mockCall := s.mock.On("DeleteAccount", uint(id), v.userId).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/account/"+v.ParamID, nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/accounts/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamID)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.accountController.DeleteAccount(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			if v.HasReturnBody {
				var resp map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&resp)
				s.NoError(err)
				s.Equal(v.ExpectedMesaage, resp["message"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteAccount) TestGetAccountByUser() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		MockReturnBody     []dto.AccountDTO
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       []dto.AccountDTO
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get account by user",
			"GET",
			1,
			[]dto.AccountDTO{
				{
					ID:      1,
					UserID:  1,
					Name:    "BCA",
					Balance: 100000,
				},
				{
					ID:      2,
					UserID:  1,
					Name:    "Mandiri",
					Balance: 100000,
				},
			},
			nil,
			true,
			[]dto.AccountDTO{
				{
					ID:      1,
					UserID:  1,
					Name:    "BCA",
					Balance: 100000,
				},
				{
					ID:      2,
					UserID:  1,
					Name:    "Mandiri",
					Balance: 100000,
				},
			},
			http.StatusOK,
			"success get account by user",
		},
		{
			"fail get account by user",
			"GET",
			1,
			[]dto.AccountDTO{},
			errors.New("error"),
			false,
			[]dto.AccountDTO{},
			http.StatusInternalServerError,
			"fail get account by user",
		},
	}
	for _, v := range testCase {
		mockCall := s.mock.On("GetAccountByUser", v.userId).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/account", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/accounts")
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.accountController.GetAccountByUser(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			if v.HasReturnBody {
				var resp map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&resp)
				s.NoError(err)

				s.Equal(v.ExpectedMesaage, resp["message"])
				s.Equal(v.ExpectedBody[0].Name, resp["accounts"].([]interface{})[0].(map[string]interface{})["name"])
				s.Equal(v.ExpectedBody[0].Balance, resp["accounts"].([]interface{})[0].(map[string]interface{})["balance"])
				s.Equal(v.ExpectedBody[1].Name, resp["accounts"].([]interface{})[1].(map[string]interface{})["name"])
				s.Equal(v.ExpectedBody[1].Balance, resp["accounts"].([]interface{})[1].(map[string]interface{})["balance"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteAccount) TestUpdateAccount() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		Body               dto.AccountDTO
		ParamID            string
		MockReturnError    error
		MockParamBody      dto.AccountDTO
		HasReturnBody      bool
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success update account",
			"PUT",
			1,
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			"1",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			true,
			http.StatusOK,
			"success update account",
		},
		{
			"fail bind data",
			"PUT",
			1,
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			"1",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			true,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"fail get id",
			"PUT",
			1,
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			"w",
			nil,
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			true,
			http.StatusBadRequest,
			"fail get id",
		},
		{
			"fail update account",
			"PUT",
			1,
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			"1",
			errors.New(constantError.ErrorNotAuthorized),
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			true,
			http.StatusUnauthorized,
			"fail update account",
		},
		{
			"fail update account",
			"PUT",
			1,
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			"1",
			errors.New("error"),
			dto.AccountDTO{
				ID:      1,
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			true,
			http.StatusInternalServerError,
			"fail update account",
		},
	}
	for i, v := range testCase {
		mockCall := s.mock.On("UpdateAccount", v.MockParamBody).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			// Create request
			r := httptest.NewRequest(v.Method, "/account", bytes.NewBuffer(res))
			if i != 1 {
				r.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/accounts")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamID)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.accountController.UpdateAccount(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			if v.HasReturnBody {
				var resp map[string]interface{}
				err := json.NewDecoder(w.Result().Body).Decode(&resp)
				s.NoError(err)

				s.Equal(v.ExpectedMesaage, resp["message"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteAccount) TestCreateAccount() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		Body               dto.AccountDTO
		MockReturnError    error
		MockParamBody      dto.AccountDTO
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success create account",
			"POST",
			1,
			dto.AccountDTO{
				Name:    "BCA",
				Balance: 100000,
			},
			nil,
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			http.StatusOK,
			"success create account",
		},
		{
			"fail bind data",
			"POST",
			1,
			dto.AccountDTO{
				Name:    "BCA",
				Balance: 100000,
			},
			nil,
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"There is an empty field",
			"POST",
			1,
			dto.AccountDTO{
				Balance: 100000,
			},
			nil,
			dto.AccountDTO{
				UserID:  1,
				Balance: 100000,
			},
			http.StatusBadRequest,
			"There is an empty field",
		},
		{
			"fail create account",
			"POST",
			1,
			dto.AccountDTO{
				Name:    "BCA",
				Balance: 100000,
			},
			errors.New(constantError.ErrorNotAuthorized),
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			http.StatusUnauthorized,
			"fail create account",
		},
		{
			"fail create account",
			"POST",
			1,
			dto.AccountDTO{
				Name:    "BCA",
				Balance: 100000,
			},
			errors.New("error"),
			dto.AccountDTO{
				UserID:  1,
				Name:    "BCA",
				Balance: 100000,
			},
			http.StatusInternalServerError,
			"fail create account",
		},
	}
	for i, v := range testCase {
		mockCall := s.mock.On("CreateAccount", v.MockParamBody).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			// Create request
			r := httptest.NewRequest(v.Method, "/account", bytes.NewBuffer(res))
			if i != 1 {
				r.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			e.Validator = &helper.CustomValidator{
				Validator: validator.New(),
			}
			ctx := e.NewContext(r, w)
			ctx.SetPath("/accounts")
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.accountController.CreateAccount(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
		})
		// remove mock
		mockCall.Unset()
	}
}

func TestSuiteAccount(t *testing.T) {
	suite.Run(t, new(suiteAccount))
}
