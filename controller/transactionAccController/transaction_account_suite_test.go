package transactionAccController

import (
	"bytes"
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	transactionAccMockService "dompet-miniprojectalta/service/transactionAccService/transactionAccMock"
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

type suiteTransactionAcc struct {
	suite.Suite
	transactionAccController *TransactionAccController
	mock                     *transactionAccMockService.TransactionAccMock
}

func (s *suiteTransactionAcc) SetupTest() {
	mock := &transactionAccMockService.TransactionAccMock{}
	s.mock = mock
	s.transactionAccController = &TransactionAccController{
		TransAccService: s.mock,
	}
}

func (s *suiteTransactionAcc) TestGetTransactionAccount() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		ParamMonth         string
		MockReturnBody     []dto.GetTransactionAccountDTO
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       []dto.GetTransactionAccountDTO
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get transaction account",
			"GET",
			1,
			"11",
			[]dto.GetTransactionAccountDTO{
				{
					ID:        1,
					AccountFromID: 1,
					AccountToID: 2,
					Amount: 100000,
					UserID: 1,
					Note: "test",
					CreatedAt: time.Now(),
					AdminFee: 0,
				},
			},
			nil,
			true,
			[]dto.GetTransactionAccountDTO{
				{
					ID:        1,
					AccountFromID: 1,
					AccountToID: 2,
					Amount: 100000,
					UserID: 1,
					Note: "test",
					CreatedAt: time.Now(),
					AdminFee: 0,
				},
			},
			http.StatusOK,
			"success get transaction account",
		},
		{
			"success get transaction account but no param month",
			"GET",
			1,
			"",
			[]dto.GetTransactionAccountDTO{
				{
					ID:        1,
					AccountFromID: 1,
					AccountToID: 2,
					Amount: 100000,
					UserID: 1,
					Note: "test",
					CreatedAt: time.Now(),
					AdminFee: 0,
				},
			},
			nil,
			true,
			[]dto.GetTransactionAccountDTO{
				{
					ID:        1,
					AccountFromID: 1,
					AccountToID: 2,
					Amount: 100000,
					UserID: 1,
					Note: "test",
					CreatedAt: time.Now(),
					AdminFee: 0,
				},
			},
			http.StatusOK,
			"success get transaction account",
		},
		{
			"fail get month",
			"GET",
			1,
			"w",
			[]dto.GetTransactionAccountDTO{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			[]dto.GetTransactionAccountDTO{},
			http.StatusBadRequest,
			"fail get month",
		},
		{
			"fail get transaction account with error auth",
			"GET",
			1,
			"11",
			[]dto.GetTransactionAccountDTO{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			[]dto.GetTransactionAccountDTO{},
			http.StatusUnauthorized,
			"fail get transaction account",
		},
		{
			"fail get transaction account with error internal server",
			"GET",
			1,
			"11",
			[]dto.GetTransactionAccountDTO{},
			errors.New("error"),
			false,
			[]dto.GetTransactionAccountDTO{},
			http.StatusInternalServerError,
			"fail get transaction account",
		},
	}
	for _, v := range testCase {
		var month int
		if v.ParamMonth == "" {
			month = int(time.Now().Month())
		}else {
			month, _ = strconv.Atoi(v.ParamMonth)
		}
		mockCall := s.mock.On("GetTransactionAccount", v.userId, month).Return(v.MockReturnBody, v.MockReturnError)
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

			err := s.transactionAccController.GetTransactionAccount(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody[0].Amount, resp["transaction_account_month_" + strconv.Itoa(month)].([]interface{})[0].(map[string]interface{})["amount"])
				s.Equal(v.ExpectedBody[0].AccountFromID, uint(resp["transaction_account_month_" + strconv.Itoa(month)].([]interface{})[0].(map[string]interface{})["account_from_id"].(float64)))
				s.Equal(v.ExpectedBody[0].AccountToID, uint(resp["transaction_account_month_" + strconv.Itoa(month)].([]interface{})[0].(map[string]interface{})["account_to_id"].(float64)))
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteTransactionAcc) TestDeleteTransactionAccount() {
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
			"success delete transaction account",
			"DELETE",
			1,
			"1",
			nil,
			http.StatusOK,
			"success delete transaction account",
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
			"fail delete transaction account with error auth",
			"DELETE",
			1,
			"1",
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail delete transaction account",
		},
		{
			"fail delete transaction account with error internal server",
			"DELETE",
			1,
			"1",
			errors.New("error"),
			http.StatusInternalServerError,
			"fail delete transaction account",
		},
	}
	for _, v := range testCase {
		id, _ := strconv.Atoi(v.ParamId)
		mockCall := s.mock.On("DeleteTransactionAccount", uint(id), v.userId).Return(v.MockReturnError)
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

			err := s.transactionAccController.DeleteTransactionAccount(ctx)
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

func (s *suiteTransactionAcc) TestCreateTransactionAccount() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		Body               dto.TransactionAccount
		MockReturnError    error
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success create transaction account",
			"POST",
			1,
			dto.TransactionAccount{
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee: 	0,
			},
			nil,
			http.StatusOK,
			"success create transaction account",
		},
		{
			"fail bind data",
			"POST",
			1,
			dto.TransactionAccount{
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee: 	0,
			},
			nil,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"There is an empty field",
			"POST",
			1,
			dto.TransactionAccount{
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee: 	0,
			},
			nil,
			http.StatusBadRequest,
			"There is an empty field",
		},
		{
			"fail create transaction account with error auth",
			"POST",
			1,
			dto.TransactionAccount{
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee: 	0,
			},
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail create transaction account",
		},
		{
			"fail create transaction account with error internal server",
			"POST",
			1,
			dto.TransactionAccount{
				AccountFromID: 1,
				AccountToID:   2,
				Amount:        10000,
				Note:          "test",
				AdminFee: 	0,
			},
			errors.New("error"),
			http.StatusInternalServerError,
			"fail create transaction account",
		},
	}
	for i, v := range testCase {
		body := dto.TransactionAccount{
			AccountFromID: v.Body.AccountFromID,
			AccountToID:   v.Body.AccountToID,
			Amount:        v.Body.Amount,
			Note:          v.Body.Note,
			AdminFee: 	v.Body.AdminFee,
			UserID: v.userId,
		}
		mockCall := s.mock.On("CreateTransactionAccount", body).Return(v.MockReturnError)
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

			err := s.transactionAccController.CreateTransactionAccount(ctx)
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

func TestSuiteTransactionAcc(t *testing.T) {
	suite.Run(t, new(suiteTransactionAcc))
}
