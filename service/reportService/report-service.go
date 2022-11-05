package reportService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/reportRepository"
	"time"
)

type ReportService interface {
	GetReportByUser(userID uint) ([]dto.TransactionReport, error)
	GetTransaction(userID uint, periodDate map[string]time.Time, categoryID uint) ([]dto.TransactionReport, float64, error)
	GetAnalyticPeriod(userId uint, period string) (map[string]interface{}, error)
}

type reportService struct {
	reportRepo reportRepository.ReportRepository
}

// GetReportByUser implements ReportService
func (*reportService) GetReportByUser(userID uint) ([]dto.TransactionReport, error) {
	return nil, nil
}

// GetTransaction implements ReportService
func (as *reportService) GetTransaction(userID uint, periodDate map[string]time.Time, categoryID uint) ([]dto.TransactionReport, float64, error) {
	transaction, err := as.reportRepo.GetTransactionPeriod(userID, periodDate, categoryID)
	if err != nil {
		return nil, 0, err
	}

	var total float64 = 0
	for _, v := range transaction {
		total += v.Amount
	}

	return transaction, total, nil
}

// GetAnalyticPeriod implements ReportService
func (as *reportService) GetAnalyticPeriod(userId uint, period string) (map[string]interface{}, error) {
	// make period time
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	var periodDateCurrent map[string]time.Time
	var periodDatePrevious map[string]time.Time

	// check if period is month or week
	if period == "month" {
		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		firstOfMonthPrevious := time.Date(currentYear, currentMonth-1, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonthPrevious := firstOfMonth.AddDate(0, 0, -1)

		periodDateCurrent = map[string]time.Time{
			"firstOfDate": firstOfMonth,
			"lastOfDate":  lastOfMonth,
		}
		periodDatePrevious = map[string]time.Time{
			"firstOfDate": firstOfMonthPrevious,
			"lastOfDate":  lastOfMonthPrevious,
		}
	} else if period == "week" {
	}

	// call repository to get report expense current
	var categoryExpense uint = 2
	_, totalCurrentExpense, err := as.GetTransaction(userId, periodDateCurrent, categoryExpense)
	if err != nil {
		return nil, err
	}

	// call repository to get report expense pervious
	_, totalPreviousExpense, err := as.GetTransaction(userId, periodDatePrevious, categoryExpense)
	if err != nil {
		return nil, err
	}

	// call repository to get report income current
	var categoryIncome uint = 3
	_, totalCurrentIncome, err := as.GetTransaction(userId, periodDateCurrent, categoryIncome)
	if err != nil {
		return nil, err
	}

	// call repository to get report income pervious
	_, totalPreviousIncome, err := as.GetTransaction(userId, periodDatePrevious, categoryIncome)
	if err != nil {
		return nil, err
	}

	// Net income
	netIncome := totalCurrentIncome + totalCurrentExpense

	// comparison between current and previous
	var comparisonExpense int64
	if totalPreviousExpense == 0 {
		comparisonExpense = 0
	} else {
		comparisonExpense = (int64(totalCurrentExpense) - int64(totalPreviousExpense)) * 100 / int64(totalPreviousExpense)
	}
	var comparisonIncome int64
	if totalPreviousIncome == 0 {
		comparisonIncome = 0
	} else {
		comparisonIncome = (int64(totalCurrentIncome) - int64(totalPreviousIncome)) * 100 / int64(totalPreviousIncome)
	}

	data := map[string]interface{}{
		"Expense":            totalCurrentExpense,
		"Income":             totalCurrentIncome,
		"Net income":         netIncome,
		"Comparison expense": comparisonExpense,
		"Comparison income":  comparisonIncome,
	}

	return data, err
}

func NewReportService(reportRepo reportRepository.ReportRepository) ReportService {
	return &reportService{
		reportRepo: reportRepo,
	}
}
