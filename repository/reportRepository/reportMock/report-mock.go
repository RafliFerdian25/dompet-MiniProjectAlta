package reportMockService

import (
	"dompet-miniprojectalta/models/dto"

	"github.com/stretchr/testify/mock"
)

type ReportMock struct {
	mock.Mock
}
func (r *ReportMock) GetReportbyCategory(userId uint, period map[string]interface{}, categoryID uint) ([]dto.ReportSpendingCategoryPeriod, error) {
	args := r.Called(userId, period, categoryID)

	return args.Get(0).([]dto.ReportSpendingCategoryPeriod), args.Error(1)
}

func (r *ReportMock) GetTransactionPeriod(userID uint, period string, categoryID uint, limit int) ([]dto.TransactionReportPeriod, error) {
	args := r.Called(userID, period, categoryID, limit)

	return args.Get(0).([]dto.TransactionReportPeriod), args.Error(1)
}