package reportRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"

	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

// GetReportbyCategory implements ReportRepository
func (ar *reportRepository) GetReportbyCategory(userID uint, period map[string]interface{}, categoryID uint) ([]dto.ReportSpendingCategoryPeriod, error) {
	// get current report
	var reportPeriod []dto.ReportSpendingCategoryPeriod
	var err error
	if period["period"] == "month" {
		err = ar.db.Model(&model.Transaction{}).Select("sub_categories.name as sub_category, date_format(transactions.created_at, ?) as period, sum(transactions.amount) as total", period["format"]).Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.user_id = ? AND sub_categories.category_id = ? AND MONTH(transactions.created_at) = ?", userID, categoryID, period["numberPeriod"]).Group("period, transactions.sub_category_id").Order("transactions.created_at DESC").Scan(&reportPeriod).Error
	} else if period["period"] == "week" {
		err = ar.db.Model(&model.Transaction{}).Select("sub_categories.name as sub_category, date_format(transactions.created_at, ?) as period, sum(transactions.amount) as total", period["format"]).Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.user_id = ? AND sub_categories.category_id = ? AND WEEK(transactions.created_at) = ?", userID, categoryID, period["numberPeriod"]).Group("period, transactions.sub_category_id").Order("transactions.created_at DESC").Scan(&reportPeriod).Error
	}
	if err != nil {
		return nil, err
	}

	return reportPeriod, nil
}

// GetReport implements ReportRepository
func (ar *reportRepository) GetTransactionPeriod(userID uint, period string, categoryID uint, limit int) ([]dto.TransactionReportPeriod, error) {
	// get current report
	var reportPeriod []dto.TransactionReportPeriod
	err := ar.db.Model(&model.Transaction{}).Select("date_format(transactions.created_at, ?) as period, sum(transactions.amount) as total", period).Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.user_id = ? AND sub_categories.category_id = ?", userID, categoryID).Group("period").Limit(limit).Order("transactions.created_at DESC").Scan(&reportPeriod).Error
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
