package debtController

import (
	"bytes"
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	debtMockService "dompet-miniprojectalta/service/debtService/debtMock"
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

type suiteDebt struct {
	suite.Suite
	debtController *DebtController
	mock           *debtMockService.DebtMock
}

func (s *suiteDebt) SetupTest() {
	mock := &debtMockService.DebtMock{}
	s.mock = mock
	s.debtController = &DebtController{
		DebtService: s.mock,
	}
}

func (s *suiteDebt) TestGetDebt() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		ParamDebtStatus    string
		MockReturnBody     map[string][]dto.GetDebtTransactionResponse
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       map[string][]dto.GetDebtTransactionResponse
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get debt",
			"GET",
			1,
			"unpaid",
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			nil,
			true,
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			http.StatusOK,
			"success get debt",
		},
		{
			"fail get debt status",
			"GET",
			1,
			"",
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			nil,
			false,
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			http.StatusBadRequest,
			"fail get debt status",
		},
		{
			"fail get debt status",
			"GET",
			1,
			"allpaid",
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			nil,
			false,
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			http.StatusBadRequest,
			"fail get debt status",
		},
		{
			"fail get debt",
			"GET",
			1,
			"unpaid",
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			http.StatusUnauthorized,
			"fail get debt",
		},
		{
			"fail get debt",
			"GET",
			1,
			"unpaid",
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			errors.New("error"),
			false,
			map[string][]dto.GetDebtTransactionResponse{
				"debt": {
					{
						ID:            1,
						Name:          "rafli",
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
						Name:          "rafli",
						SubCategoryID: 1,
						AccountID:     1,
						Total:         10000,
						Remaining:     10000,
						Note:          "test",
						CreatedAt:     time.Now(),
						DebtStatus:    "unpaid",
					},
				},
			},
			http.StatusInternalServerError,
			"fail get debt",
		},
	}
	for _, v := range testCase {
		mockCall := s.mock.On("GetDebt", v.userId, v.ParamDebtStatus).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/debts", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/debts")
			ctx.QueryParams().Add("debt_status", v.ParamDebtStatus)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.debtController.GetDebt(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody["debt"][0].Name, resp["data"].(map[string]interface{})["debt"].([]interface{})[0].(map[string]interface{})["name_lender"])
				s.Equal(v.ExpectedBody["debt"][0].Total, resp["data"].(map[string]interface{})["debt"].([]interface{})[0].(map[string]interface{})["total"])
				s.Equal(v.ExpectedBody["loan"][0].Name, resp["data"].(map[string]interface{})["loan"].([]interface{})[0].(map[string]interface{})["name_lender"])
				s.Equal(v.ExpectedBody["loan"][0].Total, resp["data"].(map[string]interface{})["loan"].([]interface{})[0].(map[string]interface{})["total"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteDebt) TestDeleteDebt() {
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
			"success delete debt",
			"DELETE",
			1,
			"1",
			nil,
			http.StatusOK,
			"success delete debt",
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
			"fail delete debt",
			"DELETE",
			1,
			"1",
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail delete debt",
		},
		{
			"fail delete debt",
			"DELETE",
			1,
			"1",
			errors.New("error"),
			http.StatusInternalServerError,
			"fail delete debt",
		},
	}
	for _, v := range testCase {
		id, _ := strconv.Atoi(v.ParamId)
		mockCall := s.mock.On("DeleteDebt", uint(id), v.userId).Return(v.MockReturnError)
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

			err := s.debtController.DeleteDebt(ctx)
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

func (s *suiteDebt) TestCreateDebt() {
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		Body            dto.DebtTransactionDTO
		MockReturnError    error
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success create debt",
			"POST",
			1,
			dto.DebtTransactionDTO{
				Name: "rafli",
				UserID: 1,
				SubCategoryID: 1,
				AccountID: 1,
				Amount: 10000,
			},
			nil,
			http.StatusOK,
			"success create debt",
		},
		{
			"fail bind data",
			"POST",
			1,
			dto.DebtTransactionDTO{
				Name: "rafli",
				UserID: 1,
				SubCategoryID: 1,
				AccountID: 1,
				Amount: 10000,
			},
			nil,
			http.StatusInternalServerError,
			"fail bind data",
		},
		{
			"There is an empty field",
			"POST",
			1,
			dto.DebtTransactionDTO{
				UserID: 1,
				SubCategoryID: 1,
				AccountID: 1,
				Amount: 10000,
			},
			nil,
			http.StatusBadRequest,
			"There is an empty field",
		},
		{
			"debt id is required",
			"POST",
			1,
			dto.DebtTransactionDTO{
				Name: "rafli",
				UserID: 1,
				SubCategoryID: 2,
				AccountID: 1,
				Amount: 10000,
			},
			nil,
			http.StatusBadRequest,
			"fail create debt",
		},
		{
			"fail create debt",
			"POST",
			1,
			dto.DebtTransactionDTO{
				Name: "rafli",
				UserID: 1,
				SubCategoryID: 1,
				AccountID: 1,
				Amount: 10000,
			},
			errors.New(constantError.ErrorNotAuthorized),
			http.StatusUnauthorized,
			"fail create debt",
		},
		{
			"fail create debt",
			"POST",
			1,
			dto.DebtTransactionDTO{
				Name: "rafli",
				UserID: 1,
				SubCategoryID: 1,
				AccountID: 1,
				Amount: 10000,
			},
			errors.New("error"),
			http.StatusInternalServerError,
			"fail create debt",
		},
	}
	for i, v := range testCase {
		mockCall := s.mock.On("CreateDebt", v.Body).Return(v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			res, _ := json.Marshal(v.Body)
			// Create request
			r := httptest.NewRequest(v.Method, "/debts", bytes.NewBuffer(res))
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
			ctx.SetPath("/debts")
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.debtController.CreateDebt(ctx)
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

func TestSuiteDebt(t *testing.T) {
	suite.Run(t, new(suiteDebt))
}
