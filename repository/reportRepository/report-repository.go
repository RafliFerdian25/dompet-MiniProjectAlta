package reportRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type ReportRepository interface {
	// GetTransactionReport(userID uint, categoryID uint) ([]dto.TransactionReport, error)
	GetTransactionPeriod(userID uint, period string, categoryID uint) ([]dto.TransactionReportPeriod, error)
}