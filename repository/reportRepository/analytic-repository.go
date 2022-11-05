package reportRepository

import (
	"dompet-miniprojectalta/models/dto"
	"time"
)

type ReportRepository interface {
	// GetTransactionReport(userID uint, categoryID uint) ([]dto.TransactionReport, error)
	GetTransactionPeriod(userID uint, periodDate map[string]time.Time, categoryID uint) ([]dto.TransactionReport, error)
}