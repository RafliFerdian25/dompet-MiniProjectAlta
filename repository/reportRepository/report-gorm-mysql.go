package reportRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"

	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

// GetReport implements ReportRepository
func (ar *reportRepository) GetTransactionPeriod(userID uint, period string, categoryID uint) ([]dto.TransactionReportPeriod, error) {
	// get current report
	var reportPeriod []dto.TransactionReportPeriod
	err := ar.db.Model(&model.Transaction{}).Select("date_format(transactions.created_at, ?) as period, sum(transactions.amount) as total", period).Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.user_id = ? AND sub_categories.category_id = ?", userID, categoryID).Group("period").Order("transactions.created_at DESC").Scan(&reportPeriod).Error
	if err != nil {
		return nil, err
	}
	
	return reportPeriod, nil
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}
