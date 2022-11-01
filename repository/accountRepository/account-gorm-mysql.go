package accountRepository

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/models/model"
	"errors"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

// DeleteAccount implements AccountRepository
func (ar *accountRepository) DeleteAccount(id uint) error {
	// delete data account from database by id
	err := ar.db.Model(&model.Account{}).Where("id = ?", id).Delete(&model.Account{})
	if err.Error != nil {
		return err.Error
	}
	if err.RowsAffected <= 0 {
		return errors.New("account not found")
	}

	return nil
}

// GetAccountByUser implements AccountRepository
func (ar *accountRepository) GetAccountByUser(userId string) ([]dto.AccountDTO, error) {
	var userAccounts []dto.AccountDTO
	// get data sub category from database by user
	err := ar.db.Model(&model.Account{}).Where("user_id = ?", userId).Or("user_id IS NULL").Find(&userAccounts).Error
	if err != nil {
		return nil, err
	}
	return userAccounts, nil
}

// GetAccountById implements AccountRepository
func (ar *accountRepository) GetAccountById(id uint) (dto.AccountDTO, error) {
	var account dto.AccountDTO
	err := ar.db.Model(&model.Account{}).First(&account, id).Error
	if err != nil {
		return dto.AccountDTO{}, err
	}
	return account, nil
}

// UpdateAccount implements AccountRepository
func (ar *accountRepository) UpdateAccount(account dto.AccountDTO) error {
	// update account with new data
	err := ar.db.Model(&model.Account{}).Where("id = ?", account.ID).Updates(&model.Account{
		Name: account.Name,
	})
	if err.Error != nil {
		return err.Error
	}
	if err.RowsAffected <= 0 {
		return errors.New("account not found")
	}

	return nil
}

// CreateAccount implements AccountRepository
func (ar *accountRepository) CreateAccount(account dto.AccountDTO) error {
	err := ar.db.Model(&model.Account{}).Create(&model.Account{
		UserID:  account.UserID,
		Name:    account.Name,
		Balance: account.Balance,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}
