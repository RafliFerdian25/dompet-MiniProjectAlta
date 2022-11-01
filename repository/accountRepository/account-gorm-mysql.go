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
