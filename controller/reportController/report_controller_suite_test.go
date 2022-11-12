package reportController

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/helper"
	"dompet-miniprojectalta/models/dto"
	reportMockService "dompet-miniprojectalta/service/reportService/reportMock"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type suiteReport struct {
	suite.Suite
	reportController *ReportController
	mock             *reportMockService.ReportMock
}

func (s *suiteReport) SetupTest() {
	mock := &reportMockService.ReportMock{}
	s.mock = mock
	s.reportController = &ReportController{
		ReportService: s.mock,
	}
}

func (s *suiteReport) TestGetCashflow() {
	// Setup
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		ParamPeriod        string
		MockReturnBody     map[string]int64
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       map[string]int64
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get cashflow",
			"GET",
			1,
			"month",
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			nil,
			true,
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			http.StatusOK,
			"success get cashflow",
		},
		{
			"fail get period",
			"GET",
			1,
			"year",
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			nil,
			false,
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			http.StatusBadRequest,
			"fail get period",
		},
		{
			"fail get cashflow",
			"GET",
			1,
			"month",
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			http.StatusUnauthorized,
			"fail get cashflow",
		},
		{
			"fail get cashflow",
			"GET",
			1,
			"month",
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			errors.New("error"),
			false,
			map[string]int64{
				"total_income_45_2022":  1000000,
				"total_expense_45_2022": 1000000,
				"cashflow_45_2022":      1000000,
			},
			http.StatusInternalServerError,
			"fail get cashflow",
		},
	}
	for _, v := range testCase {
		mockCall := s.mock.On("GetCashflow", v.userId, v.ParamPeriod).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/reports", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/reports")
			ctx.QueryParams().Add("period", v.ParamPeriod)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.reportController.GetCashflow(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody["total_income_45_2022"], int64(resp["data"].(map[string]interface{})["total_income_45_2022"].(float64)))
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteReport) TestGetReportbyCategory() {
	_, numberWeek := time.Now().ISOWeek()
	testCase := []struct {
		Name               string
		Method             string
		UserId             uint
		NumberPeriod       int
		ParamPeriod        string
		ParamNumberPeriod  string
		MockReturnBody     map[string]interface{}
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       map[string]interface{}
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get report by category",
			"GET",
			1,
			11,
			"month",
			"11",
			map[string]interface{}{
				"expense_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       -1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_expense": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total":  float64(1000000),
				},
				"income_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_income": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total":  float64(1000000),
				},
			},
			nil,
			true,
			map[string]interface{}{
				"expense_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       -1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_expense": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total":  float64(1000000),
				},
				"income_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_income": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total": float64(1000000),
				},
			},
			http.StatusOK,
			"success get report by category",
		},
		{
			"success get report by category",
			"GET",
			1,
			numberWeek,
			"week",
			"",
			map[string]interface{}{
				"expense_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       -1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_expense": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total":  float64(1000000),
				},
				"income_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_income": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total": float64(1000000),
				},
			},
			nil,
			true,
			map[string]interface{}{
				"expense_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       -1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_expense": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total":  float64(1000000),
				},
				"income_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "test",
						Total:       1000000,
						Period:      "November_2022",
						Persentage:  float64(50),
					},
				},
				"total_income": map[string]interface{}{
					"period": "month",
					"number": 11,
					"total": float64(1000000),
				},
			},
			http.StatusOK,
			"success get report by category",
		},
		{
			"fail get period",
			"GET",
			1,
			11,
			"year",
			"11",
			map[string]interface{}{},
			nil,
			false,
			map[string]interface{}{},
			http.StatusBadRequest,
			"fail get period",
		},
		{
			"fail get report by category",
			"GET",
			1,
			int(time.Now().Month()),
			"month",
			"",
			map[string]interface{}{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			map[string]interface{}{},
			http.StatusUnauthorized,
			"fail get report by category",
		},
		{
			"fail get report by category",
			"GET",
			1,
			11,
			"month",
			"",
			map[string]interface{}{},
			errors.New("error"),
			false,
			map[string]interface{}{},
			http.StatusInternalServerError,
			"fail get report by category",
		},
		{
			"invalid number period",
			"GET",
			1,
			11,
			"month",
			"qwe",
			map[string]interface{}{},
			nil,
			false,
			map[string]interface{}{},
			http.StatusBadRequest,
			"invalid number period",
		},
		{
			"invalid number period",
			"GET",
			1,
			100,
			"month",
			"100",
			map[string]interface{}{},
			nil,
			false,
			map[string]interface{}{},
			http.StatusBadRequest,
			"invalid number period",
		},
		{
			"invalid number period",
			"GET",
			1,
			100,
			"week",
			"100",
			map[string]interface{}{},
			nil,
			false,
			map[string]interface{}{},
			http.StatusBadRequest,
			"invalid number period",
		},
	}
	for _, v := range testCase {
		mockCall := s.mock.On("GetReportbyCategory", v.UserId, v.ParamPeriod, v.NumberPeriod).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/reports/subcategory", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/reports/subcategory")
			ctx.QueryParams().Add("number_period", v.ParamNumberPeriod)
			ctx.QueryParams().Add("period", v.ParamPeriod)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.reportController.GetReportbyCategory(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
			if v.HasReturnBody {
				s.Equal(v.ExpectedBody["expense_by_category"].([]dto.ReportSpendingCategoryPeriod)[0].Total, int64(resp["data"].(map[string]interface{})["expense_by_category"].([]interface{})[0].(map[string]interface{})["total"].(float64)))
				s.Equal(v.ExpectedBody["total_expense"].(map[string]interface{})["total"], resp["data"].(map[string]interface{})["total_expense"].(map[string]interface{})["total"])
				s.Equal(v.ExpectedBody["income_by_category"].([]dto.ReportSpendingCategoryPeriod)[0].Total, int64(resp["data"].(map[string]interface{})["income_by_category"].([]interface{})[0].(map[string]interface{})["total"].(float64)))
				s.Equal(v.ExpectedBody["total_income"].(map[string]interface{})["total"], resp["data"].(map[string]interface{})["total_income"].(map[string]interface{})["total"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func (s *suiteReport) TestGetAnalyticPeriod() {
	// Setup
	testCase := []struct {
		Name               string
		Method             string
		userId             uint
		ParamPeriod        string
		MockReturnBody     map[string]interface{}
		MockReturnError    error
		HasReturnBody      bool
		ExpectedBody       map[string]interface{}
		ExpectedStatusCode int
		ExpectedMesaage    string
	}{
		{
			"success get report",
			"GET",
			1,
			"month",
			map[string]interface{}{
				"expense_period": []dto.TransactionReportPeriod{
					{
						Period: "November_2022",
						Total:  -50000,
					},
				},
				"income_period": []dto.TransactionReportPeriod{
					{
						Period: "November_2022",
						Total:  100000,
					},
				},
				"net_income" : map[string]interface{}{
					"period": "november 2022",
					"result": 50000,
				},
				"comparison_expense": map[string]interface{}{
					"period_after":  "november 2022",
					"period_before": "october 2022",
					"result":        "-50%",
				},
				"comparison_income": map[string]interface{}{
					"period_after":  "november 2022",
					"period_before": "october 2022",
					"result":        "-50%",
				},
			},
			nil,
			true,
			map[string]interface{}{
				"expense_period": []dto.TransactionReportPeriod{
					{
						Period: "November_2022",
						Total:  -50000,
					},
				},
				"income_period": []dto.TransactionReportPeriod{
					{
						Period: "November_2022",
						Total:  100000,
					},
				},
				"net_income" : map[string]interface{}{
					"period": "november 2022",
					"result": 50000,
				},
				"comparison_expense": map[string]interface{}{
					"period_after":  "november 2022",
					"period_before": "october 2022",
					"result":        "-50%",
				},
				"comparison_income": map[string]interface{}{
					"period_after":  "november 2022",
					"period_before": "october 2022",
					"result":        "-50%",
				},
			},
			http.StatusOK,
			"success get report month",
		},
		{
			"fail get period",
			"GET",
			1,
			"year",
			map[string]interface{}{},
			nil,
			false,
			map[string]interface{}{},
			http.StatusBadRequest,
			"fail get period",
		},
		{
			"fail get analytic period",
			"GET",
			1,
			"month",
			map[string]interface{}{},
			errors.New(constantError.ErrorNotAuthorized),
			false,
			map[string]interface{}{},
			http.StatusUnauthorized,
			"fail get analytic period",
		},
		{
			"fail get analytic period",
			"GET",
			1,
			"month",
			map[string]interface{}{},
			errors.New("error"),
			false,
			map[string]interface{}{},
			http.StatusInternalServerError,
			"fail get analytic period",
		},
	}
	for _, v := range testCase {
		mockCall := s.mock.On("GetAnalyticPeriod", v.userId, v.ParamPeriod).Return(v.MockReturnBody, v.MockReturnError)
		s.T().Run(v.Name, func(t *testing.T) {
			// Create request
			r := httptest.NewRequest(v.Method, "/reports/analytic", nil)
			// Create response recorder
			w := httptest.NewRecorder()

			// handler echo
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath("/reports/analytic")
			ctx.QueryParams().Add("period", v.ParamPeriod)
			ctx.Set("user", &jwt.Token{Claims: &helper.JWTCustomClaims{UserID: 1, Name: "rafliferdian"}})

			err := s.reportController.GetAnalyticPeriod(ctx)
			s.NoError(err)
			s.Equal(v.ExpectedStatusCode, w.Code)

			var resp map[string]interface{}
			err = json.NewDecoder(w.Result().Body).Decode(&resp)
			s.NoError(err)

			s.Equal(v.ExpectedMesaage, resp["message"])
			if v.HasReturnBody {
				fmt.Println(v.ExpectedBody["net_income"].(map[string]interface{})["result"])
				s.Equal(v.ExpectedBody["expense_period"].([]dto.TransactionReportPeriod)[0].Total, int64(resp["data"].(map[string]interface{})["expense_period"].([]interface{})[0].(map[string]interface{})["total"].(float64)))
				s.Equal(v.ExpectedBody["income_period"].([]dto.TransactionReportPeriod)[0].Total, int64(resp["data"].(map[string]interface{})["income_period"].([]interface{})[0].(map[string]interface{})["total"].(float64)))
				s.Equal(v.ExpectedBody["net_income"].(map[string]interface{})["result"], int(resp["data"].(map[string]interface{})["net_income"].(map[string]interface{})["result"].(float64)))
				s.Equal(v.ExpectedBody["comparison_expense"].(map[string]interface{})["result"], resp["data"].(map[string]interface{})["comparison_expense"].(map[string]interface{})["result"])
				s.Equal(v.ExpectedBody["comparison_income"].(map[string]interface{})["result"], resp["data"].(map[string]interface{})["comparison_income"].(map[string]interface{})["result"])
			}
		})
		// remove mock
		mockCall.Unset()
	}
}

func TestSuiteReport(t *testing.T) {
	suite.Run(t, new(suiteReport))
}
