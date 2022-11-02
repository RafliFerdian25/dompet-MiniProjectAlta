package transactionRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"errors"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

// UpdateTransaction implements TransactionRepository
func (tr *transactionRepository) UpdateTransaction(newTransaction dto.TransactionDTO, oldTransaction dto.TransactionJoin, newAccount dto.AccountDTO) error {
	tr.db.Transaction(func(tx *gorm.DB) error {
		// update transaction
		err := tr.db.Model(&model.Transaction{}).Where("id = ?", newTransaction.ID).Updates(model.Transaction{
			SubCategoryID: newTransaction.SubCategoryID,
			AccountID:     newTransaction.AccountID,
			Amount:        newTransaction.Amount,
			Note:          newTransaction.Note,
		}).Error
		if err != nil {
			return err
		}

		// check if account change
		if newTransaction.AccountID != oldTransaction.AccountID {
			// update old account balance transaction
			oldAccount, errGetOldAccount := tr.GetAccountById(oldTransaction.AccountID)
			if errGetOldAccount != nil {
				return errGetOldAccount
			}
			oldAccount.Balance -= oldTransaction.Amount
			err := tr.db.Model(&model.Account{}).Where("id = ?", oldTransaction.AccountID).Update("balance", oldAccount.Balance)
			if err.Error != nil {
				return err.Error
			}
			if err.RowsAffected <= 0 {
				return errors.New("old account not found")
			}
		} else {
			// change new transaction amount
			newTransaction.Amount += (oldTransaction.Amount * -1)
		}

		// update new account balance transaction
		newAccount.Balance += newTransaction.Amount
		errUpdate := tr.db.Model(&model.Account{}).Where("id = ?", newTransaction.AccountID).Update("balance", newAccount.Balance)
		if errUpdate.Error != nil {
			return errUpdate.Error
		}
		if errUpdate.RowsAffected <= 0 {
			return errors.New("zaccount not found")
		}

		return nil
	})
	return nil

}

// GetTransactionById implements TransactionRepository
func (tr *transactionRepository) GetTransactionById(id uint) (dto.TransactionJoin, error) {
	var transaction dto.TransactionJoin
	err := tr.db.Model(&model.Transaction{}).Select("transactions.id, transactions.user_id, transactions.sub_category_id, sub_categories.category_id, transactions.account_id, transactions.amount").Joins("JOIN sub_categories On transactions.sub_category_id = sub_categories.id").Where("transactions.id = ?", id).Scan(&transaction)
	if err.Error != nil {
		return dto.TransactionJoin{}, err.Error
	}
	return transaction, nil
}

// GetSubCategoryById implements TransactionRepository
func (tr *transactionRepository) GetSubCategoryById(id uint) (dto.SubCategoryDTO, error) {
	var subCategory dto.SubCategoryDTO
	err := tr.db.Model(&model.SubCategory{}).First(&subCategory, id).Error
	if err != nil {
		return dto.SubCategoryDTO{}, errors.New("something went wrong when get sub category by id : " + err.Error())
	}
	return subCategory, nil
}

// GetAccountById implements TransactionRepository
func (tr *transactionRepository) GetAccountById(id uint) (dto.AccountDTO, error) {
	var account dto.AccountDTO
	err := tr.db.Model(&model.Account{}).First(&account, id).Error
	if err != nil {
		return dto.AccountDTO{}, err
	}
	return account, nil
}

// CreateTransaction implements TransactionRepository
func (tr *transactionRepository) CreateTransaction(transaction dto.TransactionDTO, categoryId uint, account dto.AccountDTO) error {
	tr.db.Transaction(func(tx *gorm.DB) error {
		// save transaction
		if err := tr.db.Model(&model.Transaction{}).Create(&model.Transaction{
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
		err := tr.db.Model(&model.Account{}).Where("id = ?", account.ID).Update("balance", account.Balance)
		if err.Error != nil {
			return err.Error
		}
		if err.RowsAffected <= 0 {
			return errors.New("subcategory not found")
		}

		return nil
	})
	return nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}
