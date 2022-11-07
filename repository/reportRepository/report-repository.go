package reportRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type ReportRepository interface {
	GetReportbyCategory(userId uint, period map[string]interface{}, categoryID uint) ([]dto.ReportSpendingCategoryPeriod, error)
	GetTransactionPeriod(userID uint, period string, categoryID uint, limit int) ([]dto.TransactionReportPeriod, error)
}