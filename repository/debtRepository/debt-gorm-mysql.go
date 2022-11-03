package debtRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"errors"

	"gorm.io/gorm"
)

type debtRepository struct {
	db *gorm.DB
}

// GetDebtById implements DebtRepostory
func (dr *debtRepository) GetDebtById(id uint) (dto.Debt, error) {
	var debt dto.Debt
	err := dr.db.Model(&model.Debt{}).Preload("Transactions").First(&debt, id).Error
	if err != nil {
		return dto.Debt{}, err
	}
	return debt, nil
}

// CreateDebt implements DebtRepository
func (dr *debtRepository) CreateDebt(debt dto.Debt, transaction dto.TransactionDTO, account dto.AccountDTO) error {
	dr.db.Transaction(func(tx *gorm.DB) error {
		var debtModel model.Debt
		// check if debt is new or not
		if transaction.DebtID == 0 {
			debtModel = model.Debt{
				Name:       debt.Name,
				Total:      debt.Total,
				Remaining:  debt.Remaining,
				Note:       debt.Note,
				DebtStatus: debt.DebtStatus,
				Status:     debt.Status,
			}
			// save debt transaction
			err := dr.db.Model(&model.Debt{}).Create(&debtModel).Error
			if err != nil {
				return err
			}
		} else {
			debtModel.ID = debt.ID
			// update debt transaction
			err := dr.db.Model(&model.Debt{}).Select("Total","Remaining","DebtStatus").Where("id = ?", transaction.DebtID).Updates(model.Debt{
				Total:      debt.Total,
				Remaining:  debt.Remaining,
				DebtStatus: debt.DebtStatus,
			}).Error
			if err != nil {
				return err
			}
		}

		// save transaction
		transaction.DebtID = debtModel.ID
		err := dr.db.Model(&model.Transaction{}).Create(&model.Transaction{
			UserID:        transaction.UserID,
			DebtID:        transaction.DebtID,
			SubCategoryID: transaction.SubCategoryID,
			AccountID:     transaction.AccountID,
			Amount:        transaction.Amount,
			Note:          transaction.Note,
		}).Error
		if err != nil {
			return err
		}

		// update account balance
		errUpdate := dr.db.Model(&model.Account{}).Where("id = ?", account.ID).Update("balance", account.Balance)
		if errUpdate.Error != nil {
			return errUpdate.Error
		}
		if errUpdate.RowsAffected <= 0 {
			return errors.New("subcategory not found")
		}

		return nil
	})
	return nil
}

// NewDebtRepository creates a new
func NewDebtRepository(db *gorm.DB) DebtRepostory {
	return &debtRepository{
		db: db,
	}
}
