package reportService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	reportMockRepository "dompet-miniprojectalta/repository/reportRepository/reportMock"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type suiteReports struct {
	suite.Suite
	reportService *reportService
	reportMock    *reportMockRepository.ReportMock
}

func (s *suiteReports) SetupSuite() {
	reportMock := &reportMockRepository.ReportMock{}
	s.reportMock = reportMock

	s.reportService = &reportService{
		reportRepo: s.reportMock,
	}
}

func (s *suiteReports) TestGetCashflow() {
	monthPeriod := strings.ToLower(time.Month(time.Now().Month()).String()) + "_" + strconv.Itoa(time.Now().Year())
	_, week := time.Now().ISOWeek()
	weekPeriod := strings.ToLower(strconv.Itoa(week)) + "_" + strconv.Itoa(time.Now().Year())
	testCase := []struct {
		Name                   string
		MockReturnErrorExpense error
		MockReturnBodyExpense  []dto.TransactionReportPeriod
		MockReturnErrorIncome  error
		MockReturnBodyIncome   []dto.TransactionReportPeriod
		ParamUserId            uint
		ParamPeriod            string
		HasReturnBody          bool
		ExpectedBody           map[string]int64
		ExpectedError          error
	}{
		{
			Name:                   "success period month",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense: []dto.TransactionReportPeriod{
				{
					Period: monthPeriod,
					Total:  -10000,
				},
			},
			MockReturnErrorIncome: nil,
			MockReturnBodyIncome: []dto.TransactionReportPeriod{
				{
					Period: monthPeriod,
					Total:  20000,
				},
			},
			ParamUserId:   1,
			ParamPeriod:   "month",
			HasReturnBody: true,
			ExpectedBody: map[string]int64{
				"total_income_" + monthPeriod:  20000,
				"total_expense_" + monthPeriod: -10000,
				"cashflow_" + monthPeriod:      10000,
			},
			ExpectedError: nil,
		},
		{
			Name:                   "success period week",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense: []dto.TransactionReportPeriod{
				{
					Period: weekPeriod,
					Total:  -10000,
				},
			},
			MockReturnErrorIncome: nil,
			MockReturnBodyIncome: []dto.TransactionReportPeriod{
				{
					Period: weekPeriod,
					Total:  20000,
				},
			},
			ParamUserId:   1,
			ParamPeriod:   "week",
			HasReturnBody: true,
			ExpectedBody: map[string]int64{
				"total_income_" + weekPeriod:  20000,
				"total_expense_" + weekPeriod: -10000,
				"cashflow_" + weekPeriod:      10000,
			},
			ExpectedError: nil,
		},
		{
			Name:                   "failed period invalid",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense:  []dto.TransactionReportPeriod{},
			MockReturnErrorIncome:  nil,
			MockReturnBodyIncome:   []dto.TransactionReportPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "year",
			HasReturnBody:          false,
			ExpectedBody: map[string]int64{
				"total_income_" + weekPeriod:  20000,
				"total_expense_" + weekPeriod: -10000,
				"cashflow_" + weekPeriod:      10000,
			},
			ExpectedError: errors.New(constantError.ErrorInvalidPeriod),
		},
		{
			Name:                   "failed get transaction expense",
			MockReturnErrorExpense: errors.New("error"),
			MockReturnBodyExpense:  []dto.TransactionReportPeriod{},
			MockReturnErrorIncome:  nil,
			MockReturnBodyIncome:   []dto.TransactionReportPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "month",
			HasReturnBody:          false,
			ExpectedBody: map[string]int64{
				"total_income_" + weekPeriod:  20000,
				"total_expense_" + weekPeriod: -10000,
				"cashflow_" + weekPeriod:      10000,
			},
			ExpectedError: errors.New("error"),
		},
		{
			Name:                   "failed get transaction income",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense:  []dto.TransactionReportPeriod{},
			MockReturnErrorIncome:  errors.New("error"),
			MockReturnBodyIncome:   []dto.TransactionReportPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "month",
			HasReturnBody:          false,
			ExpectedBody: map[string]int64{
				"total_income_" + weekPeriod:  20000,
				"total_expense_" + weekPeriod: -10000,
				"cashflow_" + weekPeriod:      10000,
			},
			ExpectedError: errors.New("error"),
		},
		{
			Name:                   "transaction is empty",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense:  []dto.TransactionReportPeriod{},
			MockReturnErrorIncome:  nil,
			MockReturnBodyIncome:   []dto.TransactionReportPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "month",
			HasReturnBody:          true,
			ExpectedBody: map[string]int64{
				"total_income_no_data":  0,
				"total_expense_no_data": 0,
				"cashflow_no_data":      0,
			},
			ExpectedError: nil,
		},
	}

	for _, v := range testCase {
		categoryExpense := uint(2)
		categoryIncome := uint(3)
		limit := 1
		var period string
		if v.ParamPeriod == "month" {
			period = "%M_%Y"
		} else if v.ParamPeriod == "week" {
			period = "%v_%x"
		}
		var mockCallExpense = s.reportMock.On("GetTransactionPeriod", v.ParamUserId, period, categoryExpense, limit).Return(v.MockReturnBodyExpense, v.MockReturnErrorExpense)
		var mockCallIncome = s.reportMock.On("GetTransactionPeriod", v.ParamUserId, period, categoryIncome, limit).Return(v.MockReturnBodyIncome, v.MockReturnErrorIncome)
		s.T().Run(v.Name, func(t *testing.T) {
			debts, err := s.reportService.GetCashflow(v.ParamUserId, v.ParamPeriod)
			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, debts)
			} else {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			}
		})
		// remove mock
		mockCallExpense.Unset()
		mockCallIncome.Unset()
	}
}

func (s *suiteReports) TestGetReportbyCategory() {
	_, month, _ := time.Now().Date()
	monthPeriod := strconv.Itoa(int(month))
	_, week := time.Now().ISOWeek()
	weekPeriod := strconv.Itoa(week)
	testCase := []struct {
		Name                   string
		MockReturnErrorExpense error
		MockReturnBodyExpense  []dto.ReportSpendingCategoryPeriod
		MockReturnErrorIncome  error
		MockReturnBodyIncome   []dto.ReportSpendingCategoryPeriod
		ParamUserId            uint
		ParamPeriod            string
		ParamNumberPeriod      int
		HasReturnBody          bool
		ExpectedBody           map[string]interface{}
		ExpectedError          error
	}{
		{
			Name:                   "success period month",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense: []dto.ReportSpendingCategoryPeriod{
				{
					SubCategory: "Entainment",
					Period:      monthPeriod,
					Total:       -10000,
					Persentage:  100,
				},
			},
			MockReturnErrorIncome: nil,
			MockReturnBodyIncome: []dto.ReportSpendingCategoryPeriod{
				{
					SubCategory: "Salary",
					Period:      monthPeriod,
					Total:       10000,
					Persentage:  100,
				},
			},
			ParamUserId:       1,
			ParamPeriod:       "month",
			ParamNumberPeriod: int(month),
			HasReturnBody:     true,
			ExpectedBody: map[string]interface{}{
				"expense_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "Entainment",
						Period:      monthPeriod,
						Total:       -10000,
						Persentage:  100,
					},
				},
				"income_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "Salary",
						Period:      monthPeriod,
						Total:       10000,
						Persentage:  100,
					},
				},
				"total_income_month_" + monthPeriod:  int64(10000),
				"total_expense_month_" + monthPeriod: int64(-10000),
			},
			ExpectedError: nil,
		},
		{
			Name:                   "success period week",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense: []dto.ReportSpendingCategoryPeriod{
				{
					SubCategory: "Entainment",
					Period:      weekPeriod,
					Total:       -10000,
					Persentage:  100,
				},
			},
			MockReturnErrorIncome: nil,
			MockReturnBodyIncome: []dto.ReportSpendingCategoryPeriod{
				{
					SubCategory: "Salary",
					Period:      weekPeriod,
					Total:       10000,
					Persentage:  100,
				},
			},
			ParamUserId:       1,
			ParamPeriod:       "week",
			ParamNumberPeriod: week,
			HasReturnBody:     true,
			ExpectedBody: map[string]interface{}{
				"expense_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "Entainment",
						Period:      weekPeriod,
						Total:       -10000,
						Persentage:  100,
					},
				},
				"income_by_category": []dto.ReportSpendingCategoryPeriod{
					{
						SubCategory: "Salary",
						Period:      weekPeriod,
						Total:       10000,
						Persentage:  100,
					},
				},
				"total_income_week_" + weekPeriod:  int64(10000),
				"total_expense_week_" + weekPeriod: int64(-10000),
			},
			ExpectedError: nil,
		},
		{
			Name:                   "failed period invalid",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense:  []dto.ReportSpendingCategoryPeriod{},
			MockReturnErrorIncome:  nil,
			MockReturnBodyIncome:   []dto.ReportSpendingCategoryPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "year",
			ParamNumberPeriod:      0,
			HasReturnBody:          false,
			ExpectedBody:           map[string]interface{}{},
			ExpectedError:          errors.New(constantError.ErrorInvalidPeriod),
		},
		{
			Name:                   "failed get category expense",
			MockReturnErrorExpense: errors.New("error"),
			MockReturnBodyExpense:  []dto.ReportSpendingCategoryPeriod{},
			MockReturnErrorIncome:  nil,
			MockReturnBodyIncome:   []dto.ReportSpendingCategoryPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "month",
			HasReturnBody:          false,
			ExpectedBody:           map[string]interface{}{},
			ExpectedError:          errors.New("error"),
		},
		{
			Name:                   "failed get category income",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense:  []dto.ReportSpendingCategoryPeriod{},
			MockReturnErrorIncome:  errors.New("error"),
			MockReturnBodyIncome:   []dto.ReportSpendingCategoryPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "month",
			HasReturnBody:          false,
			ExpectedBody:           map[string]interface{}{},
			ExpectedError:          errors.New("error"),
		},
	}

	for _, v := range testCase {
		categoryExpense := uint(2)
		categoryIncome := uint(3)
		var periodData map[string]interface{}
		if v.ParamPeriod == "month" {
			periodData = map[string]interface{}{
				"format":       "%M_%Y",
				"period":       "month",
				"numberPeriod": v.ParamNumberPeriod,
			}
		} else if v.ParamPeriod == "week" {
			periodData = map[string]interface{}{
				"format":       "%v_%x",
				"period":       "week",
				"numberPeriod": v.ParamNumberPeriod,
			}
		}
		var mockCallExpense = s.reportMock.On("GetReportbyCategory", v.ParamUserId, periodData, categoryExpense).Return(v.MockReturnBodyExpense, v.MockReturnErrorExpense)
		var mockCallIncome = s.reportMock.On("GetReportbyCategory", v.ParamUserId, periodData, categoryIncome).Return(v.MockReturnBodyIncome, v.MockReturnErrorIncome)
		s.T().Run(v.Name, func(t *testing.T) {
			debts, err := s.reportService.GetReportbyCategory(v.ParamUserId, v.ParamPeriod, v.ParamNumberPeriod)
			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, debts)
			} else {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			}
		})
		// remove mock
		mockCallExpense.Unset()
		mockCallIncome.Unset()
	}
}

func (s *suiteReports) TestGetAnalyticPeriod() {
	monthPeriod1 := strings.ToLower(time.Month(time.Now().Month()).String()) + "_" + strconv.Itoa(time.Now().Year())
	monthPeriod2 := strings.ToLower(time.Month(time.Now().AddDate(0, -1, 0).Month()).String()) + "_" + strconv.Itoa(time.Now().Year())
	_, week := time.Now().ISOWeek()
	weekPeriod1 := strings.ToLower(strconv.Itoa(week-1)) + "_" + strconv.Itoa(time.Now().Year())
	weekPeriod2 := strings.ToLower(strconv.Itoa(week-1)) + "_" + strconv.Itoa(time.Now().Year())
	testCase := []struct {
		Name                   string
		MockReturnErrorExpense error
		MockReturnBodyExpense  []dto.TransactionReportPeriod
		MockReturnErrorIncome  error
		MockReturnBodyIncome   []dto.TransactionReportPeriod
		ParamUserId            uint
		ParamPeriod            string
		HasReturnBody          bool
		ExpectedBody           map[string]interface{}
		ExpectedError          error
	}{
		{
			Name:                   "success period month",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense: []dto.TransactionReportPeriod{
				{
					Period: monthPeriod1,
					Total:  -10000,
				},
				{
					Period: monthPeriod2,
					Total:  -20000,
				},
			},
			MockReturnErrorIncome: nil,
			MockReturnBodyIncome: []dto.TransactionReportPeriod{
				{
					Period: monthPeriod1,
					Total:  20000,
				},
				{
					Period: monthPeriod2,
					Total:  40000,
				},
			},
			ParamUserId:   1,
			ParamPeriod:   "month",
			HasReturnBody: true,
			ExpectedBody: map[string]interface{}{
				"expense_period": []dto.TransactionReportPeriod{
					{
						Period: monthPeriod1,
						Total:  -10000,
					},
					{
						Period: monthPeriod2,
						Total:  -20000,
					},
				},
				"income_period": []dto.TransactionReportPeriod{
					{
						Period: monthPeriod1,
						Total:  20000,
					},
					{
						Period: monthPeriod2,
						Total:  40000,
					},
				},
				"net_income_" + monthPeriod1:                                  int64(10000),
				"comparison_expense_" + monthPeriod1 + "_and_" + monthPeriod2: "-50%",
				"comparison_income_" + monthPeriod1 + "_and_" + monthPeriod2:  "-50%",
			},
			ExpectedError: nil,
		},
		{
			Name:                   "success period week",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense: []dto.TransactionReportPeriod{
				{
					Period: weekPeriod1,
					Total:  -10000,
				},
				{
					Period: weekPeriod2,
					Total:  -20000,
				},
			},
			MockReturnErrorIncome: nil,
			MockReturnBodyIncome: []dto.TransactionReportPeriod{
				{
					Period: weekPeriod1,
					Total:  20000,
				},
				{
					Period: weekPeriod2,
					Total:  40000,
				},
			},
			ParamUserId:   1,
			ParamPeriod:   "week",
			HasReturnBody: true,
			ExpectedBody: map[string]interface{}{
				"expense_period": []dto.TransactionReportPeriod{
					{
						Period: weekPeriod1,
						Total:  -10000,
					},
					{
						Period: weekPeriod2,
						Total:  -20000,
					},
				},
				"income_period": []dto.TransactionReportPeriod{
					{
						Period: weekPeriod1,
						Total:  20000,
					},
					{
						Period: weekPeriod2,
						Total:  40000,
					},
				},
				"net_income_" + weekPeriod1:                                 int64(10000),
				"comparison_expense_" + weekPeriod1 + "_and_" + weekPeriod2: "-50%",
				"comparison_income_" + weekPeriod1 + "_and_" + weekPeriod2:  "-50%",
			},
			ExpectedError: nil,
		},
		{
			Name:                   "failed period invalid",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense:  []dto.TransactionReportPeriod{},
			MockReturnErrorIncome:  nil,
			MockReturnBodyIncome:   []dto.TransactionReportPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "year",
			HasReturnBody:          false,
			ExpectedBody:           map[string]interface{}{},
			ExpectedError:          errors.New(constantError.ErrorInvalidPeriod),
		},
		{
			Name:                   "failed get transaction expense",
			MockReturnErrorExpense: errors.New("error"),
			MockReturnBodyExpense:  []dto.TransactionReportPeriod{},
			MockReturnErrorIncome:  nil,
			MockReturnBodyIncome:   []dto.TransactionReportPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "month",
			HasReturnBody:          false,
			ExpectedBody:           map[string]interface{}{},
			ExpectedError:          errors.New("error"),
		},
		{
			Name:                   "failed get transaction income",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense:  []dto.TransactionReportPeriod{},
			MockReturnErrorIncome:  errors.New("error"),
			MockReturnBodyIncome:   []dto.TransactionReportPeriod{},
			ParamUserId:            1,
			ParamPeriod:            "month",
			HasReturnBody:          false,
			ExpectedBody:           map[string]interface{}{},
			ExpectedError:          errors.New("error"),
		},
		{
			Name:                   "income transaction is empty",
			MockReturnErrorExpense: nil,
			MockReturnBodyExpense: []dto.TransactionReportPeriod{
				{
					Period: weekPeriod1,
					Total:  -10000,
				},
			},
			MockReturnErrorIncome: nil,
			MockReturnBodyIncome:  []dto.TransactionReportPeriod{},
			ParamUserId:           1,
			ParamPeriod:           "month",
			HasReturnBody:         true,
			ExpectedBody:          map[string]interface{}{
				"expense_period": []dto.TransactionReportPeriod{
					{
						Period: weekPeriod1,
						Total:  -10000,
					},
					{
						Period: "No_Data",
						Total:  0,
					},
				},
				"income_period": []dto.TransactionReportPeriod{
					{
						Period: "No_Data",
						Total:  0,
					},
					{
						Period: "No_Data",
						Total:  0,
					},
				},
				"net_income_no_data":   int64(-10000),
				"comparison_expense_" + weekPeriod1 + "_and_no_data": "0%",
				"comparison_income_no_data_and_no_data":  "0%",
			},
			ExpectedError:         nil,
		},
	}

	for _, v := range testCase {
		categoryExpense := uint(2)
		categoryIncome := uint(3)
		limit := -1
		var period string
		if v.ParamPeriod == "month" {
			period = "%M_%Y"
		} else if v.ParamPeriod == "week" {
			period = "%v_%x"
		}
		var mockCallExpense = s.reportMock.On("GetTransactionPeriod", v.ParamUserId, period, categoryExpense, limit).Return(v.MockReturnBodyExpense, v.MockReturnErrorExpense)
		var mockCallIncome = s.reportMock.On("GetTransactionPeriod", v.ParamUserId, period, categoryIncome, limit).Return(v.MockReturnBodyIncome, v.MockReturnErrorIncome)
		s.T().Run(v.Name, func(t *testing.T) {
			debts, err := s.reportService.GetAnalyticPeriod(v.ParamUserId, v.ParamPeriod)
			if v.HasReturnBody {
				s.NoError(err)
				s.Equal(v.ExpectedBody, debts)
			} else {
				s.Error(err)
				s.Equal(v.ExpectedError, err)
			}
		})
		// remove mock
		mockCallExpense.Unset()
		mockCallIncome.Unset()
	}
}

func TestSuiteReports(t *testing.T) {
	suite.Run(t, new(suiteReports))
}
