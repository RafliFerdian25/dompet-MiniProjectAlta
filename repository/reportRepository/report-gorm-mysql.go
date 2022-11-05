package reportRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"time"

	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

// GetReport implements ReportRepository
func (ar *reportRepository) GetTransactionPeriod(userID uint, periodDate map[string]time.Time, categoryID uint) ([]dto.TransactionReport, error) {
	// get current report
	var report []dto.TransactionReport
	err := ar.db.Model(&model.Transaction{}).Select("transactions.id, transactions.user_id, transactions.sub_category_id, sub_categories.category_id, transactions.account_id, transactions.amount, transactions.created_at").Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.created_at BETWEEN ? AND ? AND transactions.user_id = ? AND sub_categories.category_id = ?", periodDate["firstOfDate"], periodDate["lastOfDate"], userID, categoryID).Scan(&report).Error
	if err != nil {
		return nil, err
	}
	
	return report, nil
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}
