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
