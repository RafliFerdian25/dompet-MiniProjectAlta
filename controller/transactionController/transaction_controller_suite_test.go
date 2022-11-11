package transactionController

import (
	"bytes"
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	transactionMockService "dompet-miniprojectalta/service/transactionService/transactionMock"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type suiteTransaction struct {
	suite.Suite
	transactionController *TransactionController
	mock                  *transactionMockService.TransactionMock
}

func (s *suiteTransaction) SetupTest() {
	mock := &transactionMockService.TransactionMock{}
	s.mock = mock
	s.transactionController = &TransactionController{
		TransactionService: s.mock,
	}
}

func (s *suiteTransaction) TestGetTransaction() {
	month := 11
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		ParamMonth         string
		MockReturnBody     map[string]interface{}
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       map[string]interface{}
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get transaction",
			"GET",
			1,
			"11",
			map[string]interface{}{
				"expense": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"income": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"month_transaction": time.Month(month).String(),
			},
			nil,
			true,
			map[string]interface{}{
				"expense": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"income": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"month_transaction": time.Month(month).String(),
			},
			http.StatusOK,
			"success get transaction",
		},
		{
			"success get transaction but no param month",
			"GET",
			1,
			"",
			map[string]interface{}{
				"expense": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"income": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"month_transaction": time.Month(month).String(),
			},
			nil,
			true,
			map[string]interface{}{
				"expense": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"income": []dto.GetTransactionDTO{
					{
						ID:            1,
						UserID:        1,
						SubCategoryID: 8,
						CategoryID:    2,
						Amount:        100000,
						Note:          "test",
						CreatedAt:     time.Now(),
					},
				},
				"month_transaction": time.Month(month).String(),
			},
			http.StatusOK,
			"success get transaction",
		},
		{
			"fail get month",
			"GET",
			1,
			"w",
			map[string]interface{}{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			map[string]interface{}{},
			http.StatusBadRequest,
			"fail get month",
		},
		{
			"fail get transaction with error auth",
			"GET",
			1,
			"11",
			map[string]interface{}{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			map[string]interface{}{},
			http.StatusUnauthorized,
			"fail get transaction",
		},
		{
			"fail get transaction with error internal server",
			"GET",
			1,
			"11",
			map[string]interface{}{},
			errors.New("error"),
			false,
			map[string]interface{}{},
			http.StatusInternalServerError,
			"fail get transaction",
		},
	}
	for _, v := range testCase {
		var month int
		if v.ParamMonth == "" {
			month = int(time.Now().Month())
		} else {
			month, _ = strconv.Atoi(v.ParamMonth)
		}
		mockCall := s.mock.On("GetTransaction", v.userId, month).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/transaction-accounts", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/transaction-accounts")
			ctx.QueryParams().Add("month", v.ParamMonth)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.transactionController.GetTransaction(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody["expense"].([]dto.GetTransactionDTO)[0].Amount, resp["transactions"].(map[string]interface{})["expense"].([]interface{})[0].(map[string]interface{})["amount"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteTransaction) TestDeleteTransaction() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		ParamId            string
		MockReturnError    error
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success delete transaction",
			"DELETE",
			1,
			"1",
			nil,
			http.StatusOK,
			"success delete transaction",
		},
		{
			"fail get id",
			"DELETE",
			1,
			"w",
			nil,
			http.StatusBadRequest,
			"fail get id",
		},
		{
			"fail delete transaction with error auth",
			"DELETE",
			1,
			"1",
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail delete transaction",
		},
		{
			"fail delete transaction with error internal server",
			"DELETE",
			1,
			"1",
			errors.New("error"),
			http.StatusInternalServerError,
			"fail delete transaction",
		},
	}
	for _, v := range testCase {
		id, _ := strconv.Atoi(v.ParamId)
		mockCall := s.mock.On("DeleteTransaction", uint(id), v.userId).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/debts", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/debts")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamId)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.transactionController.DeleteTransaction(ctx)
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

func (s *suiteTransaction) TestUpdateTransaction() {
	testCase := []struct {
		Name               string
		Method             string
		UserId             uint
		Body               dto.TransactionDTO
		ParamID            string
		MockReturnError    error
		MockParamBody      dto.TransactionDTO
		HasReturnBody      bool
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success update transaction",
			"PUT",
			1,
			dto.TransactionDTO{
				Amount:    10000,
				AccountID: 1,
			},
			"1",
			nil,
			dto.TransactionDTO{
				Amount:    10000,
				AccountID: 1,
			},
			true,
			http.StatusOK,
			"success update transaction",
		},
		{
			"fail bind data",
			"PUT",
			1,
			dto.TransactionDTO{
				Amount:    10000,
				AccountID: 1,
			},
			"1",
			nil,
			dto.TransactionDTO{
				Amount:    10000,
				AccountID: 1,
			},
			true,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"fail get id",
			"PUT",
			1,
			dto.TransactionDTO{},
			"w",
			nil,
			dto.TransactionDTO{},
			true,
			http.StatusBadRequest,
			"fail get id",
		},
		{
			"fail update transaction",
			"PUT",
			1,
			dto.TransactionDTO{},
			"1",
			errors.New(constantError.ErrorNotAuthorized),
			dto.TransactionDTO{},
			true,
			http.StatusUnauthorized,
			"fail update transaction",
		},
		{
			"fail update transaction",
			"PUT",
			1,
			dto.TransactionDTO{},
			"1",
			errors.New("error"),
			dto.TransactionDTO{},
			true,
			http.StatusInternalServerError,
			"fail update transaction",
		},
	}
	for i, v := range testCase {
		id, _ := strconv.Atoi(v.ParamID)
		body := dto.TransactionDTO{
			ID:        uint(id),
			Amount:    v.Body.Amount,
			AccountID: v.Body.AccountID,
		}
		mockCall := s.mock.On("UpdateTransaction", body, v.UserId).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			// Create request
			r := httptest.NewRequest(v.Method, "/subcategories", bytes.NewBuffer(res))
			if i != 1 {
				r.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/subcategories")
			ctx.SetParamNames("id")
			ctx.SetParamValues(v.ParamID)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.transactionController.UpdateTransaction(ctx)
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

func (s *suiteTransaction) TestCreateTransaction() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		Body               dto.TransactionDTO
		MockReturnError    error
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success create transaction",
			"POST",
			1,
			dto.TransactionDTO{
				SubCategoryID: 8,
				AccountID: 1,
				Amount:        10000,
				Note:          "test",
			},
			nil,
			http.StatusOK,
			"success create transaction",
		},
		{
			"fail bind data",
			"POST",
			1,
			dto.TransactionDTO{
				SubCategoryID: 8,
				AccountID: 1,
				Amount:        10000,
				Note:          "test",
			},
			nil,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"There is an empty field",
			"POST",
			1,
			dto.TransactionDTO{
				AccountID: 1,
				Amount:        10000,
				Note:          "test",
			},
			nil,
			http.StatusBadRequest,
			"There is an empty field",
		},
		{
			"fail create transaction with error auth",
			"POST",
			1,
			dto.TransactionDTO{
				SubCategoryID: 8,
				AccountID: 1,
				Amount:        10000,
				Note:          "test",
			},
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail create transaction",
		},
		{
			"fail create transaction with error internal server",
			"POST",
			1,
			dto.TransactionDTO{
				SubCategoryID: 8,
				AccountID: 1,
				Amount:        10000,
				Note:          "test",
			},
			errors.New("error"),
			http.StatusInternalServerError,
			"fail create transaction",
		},
	}
	for i, v := range testCase {
		body := dto.TransactionDTO{
			SubCategoryID: v.Body.SubCategoryID,
			AccountID: v.Body.AccountID,
			Amount:        v.Body.Amount,
			Note:          v.Body.Note,
			UserID: v.userId,
		}
		mockCall := s.mock.On("CreateTransaction", body).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			// Create request
			r := httptest.NewRequest(v.Method, "/subcategories", bytes.NewBuffer(res))
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
			ctx.SetPath("/subcategories")
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.transactionController.CreateTransaction(ctx)
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

func TestSuiteTransaction(t *testing.T) {
	suite.Run(t, new(suiteTransaction))
}
