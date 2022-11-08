package transactionRepository

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"errors"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

// GetTransaction implements TransactionRepository
func (tr*transactionRepository) GetTransaction(month int, userId, categoryID uint) ([]dto.GetTransactionDTO, error) {
	var transaction []dto.GetTransactionDTO
	err := tr.db.Model(&model.Transaction{}).Select("transactions.id, transactions.user_id, transactions.sub_category_id, sub_categories.category_id, transactions.account_id, transactions.amount, transactions.note, transactions.created_at").Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.user_id = ? AND MONTH(transactions.created_at) = ? AND sub_categories.category_id = ?", userId, month, categoryID).Scan(&transaction)
	if err.Error != nil {
		return []dto.GetTransactionDTO{}, err.Error
	}

	return transaction, nil
}

// DeleteTransaction implements TransactionRepository
func (tr *transactionRepository) DeleteTransaction(id uint, account dto.AccountDTO) error {
	err := tr.db.Transaction(func(tx *gorm.DB) error {
		// update account balance
		err := tx.Model(&model.Account{}).Where("id = ?", account.ID).Update("balance", account.Balance)
		if err.Error != nil {
			return err.Error
		}
		if err.RowsAffected <= 0 {
			return errors.New(constantError.ErrorAccountNotFound)
		}

		// delete transaction
		errDelete := tx.Delete(&model.Transaction{}, id)
		if errDelete.Error != nil {
			return errDelete.Error
		}
		if errDelete.RowsAffected <= 0 {
			return errors.New(constantError.ErrorTransactionNotFound)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateTransaction implements TransactionRepository
func (tr *transactionRepository) UpdateTransaction(newTransaction dto.TransactionDTO, oldTransaction dto.TransactionJoin, newAccount, oldAccount dto.AccountDTO) error {
	err := tr.db.Transaction(func(tx *gorm.DB) error {
		// update transaction
		err := tx.Model(&model.Transaction{}).Where("id = ?", newTransaction.ID).Updates(model.Transaction{
			SubCategoryID: newTransaction.SubCategoryID,
			AccountID:     newTransaction.AccountID,
			Amount:        newTransaction.Amount,
			Note:          newTransaction.Note,
		})
		if err.Error != nil {
			return err.Error
		}

		// update old account balance transaction
		err = tx.Model(&model.Account{}).Where("id = ?", oldAccount.ID).Update("balance", oldAccount.Balance)
		if err.Error != nil {
			return err.Error
		}
		if err.RowsAffected <= 0 {
			return errors.New(constantError.ErrorOldAccountNotFound)
		}

		// update new account balance transaction
		errUpdate := tx.Model(&model.Account{}).Where("id = ?", newTransaction.AccountID).Update("balance", newAccount.Balance)
		if errUpdate.Error != nil {
			return errUpdate.Error
		}
		if errUpdate.RowsAffected <= 0 {
			return errors.New(constantError.ErrorNewAccountNotFound)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetTransactionById implements TransactionRepository
func (tr *transactionRepository) GetTransactionById(id uint) (dto.TransactionJoin, error) {
	var transaction dto.TransactionJoin
	err := tr.db.Model(&model.Transaction{}).Select("transactions.id, transactions.user_id, transactions.sub_category_id, sub_categories.category_id, transactions.account_id, transactions.amount").Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.id = ?", id).Scan(&transaction)
	if err.Error != nil {
		return dto.TransactionJoin{}, err.Error
	}
	if transaction.ID != id {
		return dto.TransactionJoin{}, gorm.ErrRecordNotFound
	}
	return transaction, nil
}

// CreateTransaction implements TransactionRepository
func (tr *transactionRepository) CreateTransaction(transaction dto.TransactionDTO, categoryId uint, account dto.AccountDTO) error {
	err := tr.db.Transaction(func(tx *gorm.DB) error {
		// save transaction
		if err := tx.Model(&model.Transaction{}).Create(&model.Transaction{
			UserID:        transaction.UserID,
			SubCategoryID: transaction.SubCategoryID,
			AccountID:     transaction.AccountID,
			Amount:        transaction.Amount,
			Note:          transaction.Note,
		}).Error; err != nil {
			return err
		}
		// update account balance
		account.Balance += transaction.Amount
		err := tx.Model(&model.Account{}).Where("id = ?", account.ID).Update("balance", account.Balance)
		if err.Error != nil {
			return err.Error
		}
		if err.RowsAffected <= 0 {
			return errors.New(constantError.ErrorAccountNotFound)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}
