package reportService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/reportRepository"
	"errors"
	"strconv"
)

type ReportService interface {
	GetAnalyticPeriod(userId uint, period string) (map[string]interface{}, error)
}

type reportService struct {
	reportRepo reportRepository.ReportRepository
}


// GetAnalyticPeriod implements ReportService
func (as *reportService) GetAnalyticPeriod(userId uint, period string) (map[string]interface{}, error) {

	// check if period is month or week
	if period == "month" {
		period = "%M_%Y"
	} else if period == "week" {
		period = "%v_%x"
	} else {
		return nil, errors.New(constantError.ErrorInvalidPeriod)
	}

	// call repository to get report expense period
	var categoryExpense uint = 2
	expensePeriod, err := as.reportRepo.GetTransactionPeriod(userId, period, categoryExpense)
	if err != nil {
		return nil, err
	}
	// check if expensePeriod is empty
	if len(expensePeriod) == 0 {
		expensePeriod = append(expensePeriod, dto.TransactionReportPeriod{
			Period: "No Data",
			Total:  0,
		})
	}

	// call repository to get report income period
	var categoryIncome uint = 3
	incomePeriod, err := as.reportRepo.GetTransactionPeriod(userId,period, categoryIncome)
	if err != nil {
		return nil, err
	}
	// check if incomePeriod is empty
	if len(incomePeriod) == 0 {
		incomePeriod = append(incomePeriod, dto.TransactionReportPeriod{
			Period: "No Data",
			Total:  0,
		})
	}

	// get net income
	var netIncome int64 = 0
	netIncome = incomePeriod[0].Total + expensePeriod[0].Total

	// get expense and income current
	var expenseCurrent, incomeCurrent int64
	expenseCurrent = expensePeriod[0].Total
	incomeCurrent = incomePeriod[0].Total

	// check if len expensePeriod is 1
	if len(expensePeriod) == 1 {
		expensePeriod = append(expensePeriod, dto.TransactionReportPeriod{
			Period: "No Data",
			Total:  0,
		})
	}

	// check if len incomePeriod is 1
	if len(incomePeriod) == 1 {
		incomePeriod = append(incomePeriod, dto.TransactionReportPeriod{
			Period: "No Data",
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
		"expense_period": expensePeriod,
		"income_period":  incomePeriod,
		"net_income_" + incomePeriod[0].Period: netIncome,
		"comparison_expense_" + incomePeriod[0].Period + "_and_" + incomePeriod[1].Period: compareExpenseString + "%",
		"comparison_income_" + incomePeriod[0].Period + "_and_" + incomePeriod[1].Period:  compareIncomeString + "%",
	}

	return data, err
}

func NewReportService(reportRepo reportRepository.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}
