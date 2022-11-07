package debtRepository

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"errors"

	"gorm.io/gorm"
)

type debtRepository struct {
	db *gorm.DB
}

// GetDebt implements DebtRepostory
func (dr *debtRepository) GetDebt(userId uint, subCategory int, debtStatus string) ([]dto.GetDebtTransactionResponse, error) {
	var debts []dto.GetDebtTransactionResponse
	err := dr.db.Model(&model.Debt{}).Select("debts.id, debts.name, transactions.sub_category_id, transactions.account_id, debts.total, debts.remaining, debts.note, debts.created_at, debts.debt_status").Joins("JOIN transactions ON transactions.debt_id = debts.id").Where("debts.debt_status = ? AND transactions.sub_category_id = ? AND transactions.user_id = ?", debtStatus, subCategory, userId).Scan(&debts).Error
	if err != nil {
		return nil, err
	}
	return debts, nil
}

// DeleteDebt implements DebtRepostory
func (dr *debtRepository) DeleteDebt(id uint, account dto.AccountDTO) error {
	err := dr.db.Transaction(func(tx *gorm.DB) error {
		// update account balance
		err := tx.Model(&model.Account{}).Where("id = ?", account.ID).Update("balance", account.Balance)
		if err.Error != nil {
			return err.Error
		}
		if err.RowsAffected <= 0 {
			return errors.New(constantError.ErrorAccountNotFound)
		}

		// delete transaction
		err = tx.Model(&model.Transaction{}).Where("debt_id = ?", id).Delete(&model.Transaction{})
		if err.Error != nil {
			return err.Error
		}
		if err.RowsAffected <= 0 {
			return errors.New(constantError.ErrorTransactionNotFound)
		}

		// delete debt
		err = tx.Model(&model.Debt{}).Where("id = ?", id).Delete(&model.Debt{})
		if err.Error != nil {
			return err.Error
		}
		if err.RowsAffected <= 0 {
			return errors.New(constantError.ErrorDebtNotFound)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
	err := dr.db.Transaction(func(tx *gorm.DB) error {
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
			err := tx.Model(&model.Debt{}).Create(&debtModel).Error
			if err != nil {
				return err
			}
		} else {
			debtModel.ID = debt.ID
			// update debt transaction
			err := tx.Model(&model.Debt{}).Select("Total", "Remaining", "DebtStatus").Where("id = ?", transaction.DebtID).Updates(model.Debt{
				Total:      debt.Total,
				Remaining:  debt.Remaining,
				DebtStatus: debt.DebtStatus,
			})
			if err.Error != nil {
				return err.Error
			}
		}

		// save transaction
		transaction.DebtID = debtModel.ID
		err := tx.Model(&model.Transaction{}).Create(&model.Transaction{
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
		errUpdate := tx.Model(&model.Account{}).Where("id = ?", account.ID).Update("balance", account.Balance)
		if errUpdate.Error != nil {
			return errUpdate.Error
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// NewDebtRepository creates a new
func NewDebtRepository(db *gorm.DB) DebtRepostory {
	return &debtRepository{
		db: db,
	}
}
