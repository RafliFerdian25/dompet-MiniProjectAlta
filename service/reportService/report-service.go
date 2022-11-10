package reportService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/reportRepository"
	"errors"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type ReportService interface {
	GetCashflow(userId uint, period string) (map[string]interface{}, error)
	GetReportbyCategory(userId uint, period string, numberPeriod int) (map[string]interface{}, error)
	GetAnalyticPeriod(userId uint, period string) (map[string]interface{}, error)
}

type reportService struct {
	reportRepo reportRepository.ReportRepository
}

// GetCashflow implements ReportService
func (rs *reportService) GetCashflow(userId uint, period string) (map[string]interface{}, error) {
	// check if period is month or week
	if period == "month" {
		period = "%M_%Y"
	} else if period == "week" {
		period = "%v_%x"
	} else {
		return nil, errors.New(constantError.ErrorInvalidPeriod)
	}

	// set limit
	limit := 1

	// call repository to get transaction by period
	categoryExpense := uint(2)
	ExpenseTransactionByPeriod, err := rs.reportRepo.GetTransactionPeriod(userId, period, categoryExpense, limit)
	if err != nil {
		return nil, err
	}
	// check if ExpenseTransactionByPeriod is empty
	if len(ExpenseTransactionByPeriod) == 0 {
		ExpenseTransactionByPeriod = append(ExpenseTransactionByPeriod, dto.TransactionReportPeriod{
			Period: "No_Data",
			Total:  0,
		})
	}

	categoryIncome := uint(3)
	IncomeTransactionByPeriod, err := rs.reportRepo.GetTransactionPeriod(userId, period, categoryIncome, limit)
	if err != nil {
		return nil, err
	}
	// check if IncomeTransactionByPeriod is empty
	if len(IncomeTransactionByPeriod) == 0 {
		IncomeTransactionByPeriod = append(IncomeTransactionByPeriod, dto.TransactionReportPeriod{
			Period: "No_Data",
			Total:  0,
		})
	}

	// calculate total income and expense
	totalExpense := ExpenseTransactionByPeriod[0].Total
	totalIncome := IncomeTransactionByPeriod[0].Total

	// calculate cashflow
	cashflow := totalIncome + totalExpense

	data := map[string]interface{}{
		"total_income_" + IncomeTransactionByPeriod[0].Period:  totalIncome,
		"total_expense_" + ExpenseTransactionByPeriod[0].Period: totalExpense,
		"cashflow_"+ ExpenseTransactionByPeriod[0].Period:      cashflow,
	}
	return data, nil
}

// GetReportbyCategory implements ReportService
func (rs *reportService) GetReportbyCategory(userId uint, period string, numberPeriod int) (map[string]interface{}, error) {
	// check if period is month or week
	var periodData map[string]interface{}
	if period == "month" {
		periodData = map[string]interface{}{
			"format":       "%M_%Y",
			"period":       "month",
			"numberPeriod": numberPeriod,
		}
	} else if period == "week" {
		periodData = map[string]interface{}{
			"format":       "%v_%x",
			"period":       "week",
			"numberPeriod": numberPeriod,
		}
	} else {
		return nil, errors.New(constantError.ErrorInvalidPeriod)
	}

	// call repository to get report expense by category
	categoryExpense := uint(2)
	expenseByCategory, err := rs.reportRepo.GetReportbyCategory(userId, periodData, categoryExpense)
	if err != nil {
		return nil, err
	}
	var totalExpense int64
	for _, expense := range expenseByCategory {
		totalExpense += expense.Total
	}
	// persentage expense by category
	for i, expense := range expenseByCategory {
		persentage := float64(expense.Total*100) / float64(totalExpense)
		expenseByCategory[i].Persentage, _ = decimal.NewFromFloatWithExponent(persentage, -2).Float64()
	}

	// call repository to get report income by category
	categoryIncome := uint(3)
	incomeByCategory, err := rs.reportRepo.GetReportbyCategory(userId, periodData, categoryIncome)
	if err != nil {
		return nil, err
	}
	var totalIncome int64
	for _, income := range incomeByCategory {
		totalIncome += income.Total
	}
	// persentage income by category
	for i, income := range incomeByCategory {
		persentage := float64(income.Total*100) / float64(totalIncome)
		incomeByCategory[i].Persentage, _ = decimal.NewFromFloatWithExponent(persentage, -2).Float64()
	}

	data := map[string]interface{}{
		"expense_by_category": expenseByCategory,
		"total_expense_" + periodData["period"].(string) + "_" + strconv.Itoa(periodData["numberPeriod"].(int)): totalExpense,
		"income_by_category": incomeByCategory,
		"total_income_" + periodData["period"].(string) + "_" + strconv.Itoa(periodData["numberPeriod"].(int)): totalIncome,
	}
	return data, nil
}

// GetAnalyticPeriod implements ReportService
func (rs *reportService) GetAnalyticPeriod(userId uint, period string) (map[string]interface{}, error) {
	// check if period is month or week
	if period == "month" {
		period = "%M_%Y"
	} else if period == "week" {
		period = "%v_%x"
	} else {
		return nil, errors.New(constantError.ErrorInvalidPeriod)
	}

	// set limit
	limit := -1

	// call repository to get report expense period
	var categoryExpense uint = 2
	expensePeriod, err := rs.reportRepo.GetTransactionPeriod(userId, period, categoryExpense, limit)
	if err != nil {
		return nil, err
	}
	// check if expensePeriod is empty
	if len(expensePeriod) == 0 {
		expensePeriod = append(expensePeriod, dto.TransactionReportPeriod{
			Period: "No_Data",
			Total:  0,
		})
	}

	// call repository to get report income period
	var categoryIncome uint = 3
	incomePeriod, err := rs.reportRepo.GetTransactionPeriod(userId, period, categoryIncome, limit)
	if err != nil {
		return nil, err
	}
	// check if incomePeriod is empty
	if len(incomePeriod) == 0 {
		incomePeriod = append(incomePeriod, dto.TransactionReportPeriod{
			Period: "No_Data",
			Total:  0,
		})
	}

	// get net income
	netIncome := incomePeriod[0].Total + expensePeriod[0].Total

	// get expense and income current
	var expenseCurrent, incomeCurrent int64
	expenseCurrent = expensePeriod[0].Total
	incomeCurrent = incomePeriod[0].Total

	// check if len expensePeriod is 1
	if len(expensePeriod) == 1 {
		expensePeriod = append(expensePeriod, dto.TransactionReportPeriod{
			Period: "No_Data",
			Total:  0,
		})
	}

	// check if len incomePeriod is 1
	if len(incomePeriod) == 1 {
		incomePeriod = append(incomePeriod, dto.TransactionReportPeriod{
			Period: "No_Data",
			Total:  0,
		})
	}

	// get expense and income last
	var expenseLast, incomeLast int64
	expenseLast = expensePeriod[1].Total
	incomeLast = incomePeriod[1].Total

	// compare expense and income current with expense and income last
	var compareExpense, compareIncome int
	// persentase expense current with expense last
	if expenseLast == 0 {
		compareExpense = 0
	} else {
		compareExpense = int(expenseCurrent-expenseLast) * 100 / int(expenseLast)
	}
	// persentase income current with income last
	if incomeLast == 0 {
		compareIncome = 0
	} else {
		compareIncome = int(incomeCurrent-incomeLast) * 100 / int(incomeLast)
	}
	// convert to string
	var compareExpenseString, compareIncomeString string
	compareExpenseString = strconv.Itoa(compareExpense)
	compareIncomeString = strconv.Itoa(compareIncome)

	data := map[string]interface{}{
		"expense_period":                       expensePeriod,
		"income_period":                        incomePeriod,
		"net_income_" + strings.ToLower(incomePeriod[0].Period): netIncome,
		"comparison_expense_" + strings.ToLower(expensePeriod[0].Period) + "_and_" + strings.ToLower(expensePeriod[1].Period): compareExpenseString + "%",
		"comparison_income_" + strings.ToLower(incomePeriod[0].Period) + "_and_" + strings.ToLower(incomePeriod[1].Period):    compareIncomeString + "%",
	}

	return data, nil
}

func NewReportService(reportRepo reportRepository.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}
