package transactionAccRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"

	"gorm.io/gorm"
)

type transactionAccRepository struct {
	db *gorm.DB
}

// GetTransactionAccount implements TransactionAccRepository
func (tar *transactionAccRepository) GetTransactionAccount(userId uint, month int) ([]dto.GetTransactionAccountDTO, error) {
	// get transaction account
	var transAcc []dto.GetTransactionAccountDTO
	err := tar.db.Model(&model.TransactionAccount{}).Where("user_id = ?", userId).Where("MONTH(created_at) = ?", month).Find(&transAcc)
	if err.Error != nil {
		return []dto.GetTransactionAccountDTO{}, err.Error
	}
	return transAcc, nil
}

// DeleteTransactionAcc implements TransactionAccRepository
func (tar *transactionAccRepository) DeleteTransactionAccount(id uint, accountFrom, accountTo dto.AccountDTO) error {
	err := tar.db.Transaction(func(tx *gorm.DB) error {
		// update 'account from' balance
		if err := tx.Model(&model.Account{}).Where("id = ?", accountFrom.ID).Update("balance", accountFrom.Balance).Error; err != nil {
			return err
		}

		// update 'account to' balance
		if err := tx.Model(&model.Account{}).Where("id = ?", accountTo.ID).Update("balance", accountTo.Balance).Error; err != nil {
			return err
		}

		// delete transaction
		if err := tx.Model(&model.TransactionAccount{}).Where("id = ?", id).Delete(&model.TransactionAccount{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetTransactionAccountById implements TransactionAccRepository
func (tar *transactionAccRepository) GetTransactionAccountById(id uint) (dto.TransactionAccount, error) {
	// get transaction account
	var transAcc dto.TransactionAccount
	err := tar.db.Model(&model.TransactionAccount{}).First(&transAcc, id)
	if err.Error != nil {
		return dto.TransactionAccount{}, err.Error
	}
	return transAcc, nil
}

// CreateTransactionAccount implements TransactionAccRepository
func (tar *transactionAccRepository) CreateTransactionAccount(transAcc dto.TransactionAccount, accountFrom dto.AccountDTO, accountTo dto.AccountDTO) error {
	err := tar.db.Transaction(func(tx *gorm.DB) error {
		// create transaction account
		if err := tx.Create(&model.TransactionAccount{
			AccountFromID: accountFrom.ID,
			AccountToID:   accountTo.ID,
			UserID:        transAcc.UserID,
			Amount:        transAcc.Amount,
			Note:          transAcc.Note,
			AdminFee:      transAcc.AdminFee,
		}).Error; err != nil {
			return err
		}

		// update balance account from
		if err := tx.Model(&model.Account{}).Where("id = ?", accountFrom.ID).Update("balance", accountFrom.Balance).Error; err != nil {
			return err
		}

		// Update balance account to
		if err := tx.Model(&model.Account{}).Where("id = ?", accountTo.ID).Update("balance", accountTo.Balance).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// newTransAcc
func NewTransAccRepo(db *gorm.DB) TransactionAccRepository {
	return &transactionAccRepository{
		db: db,
	}
}
