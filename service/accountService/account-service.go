package accountService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/accountRepository"
	"errors"
)

type AccountService interface {
	DeleteAccount(id, userId uint) error
	GetAccountByUser(userId uint) ([]dto.AccountDTO, error)
	UpdateAccount(account dto.AccountDTO) error
	CreateAccount(account dto.AccountDTO) error
}

type accountService struct {
	accRepo accountRepository.AccountRepository
}

// DeleteAccount implements AccountService
func (as *accountService) DeleteAccount(id, userId uint) error {
	account, err := as.accRepo.GetAccountById(id)
	if err != nil {
		return err
	}
	// check if user id in the account is the same as the user id in the token
	if account.UserID != userId {
		return errors.New("you are not authorized to delete this account")
	}
	err = as.accRepo.DeleteAccount(id)
	if err != nil {
		return err
	}
	return nil
}

// GetAccountByUser implements AccountService
func (as *accountService) GetAccountByUser(userId uint) ([]dto.AccountDTO, error) {
	userAccounts, err := as.accRepo.GetAccountByUser(userId)
	if err != nil {
		return nil, err
	}
	return userAccounts, nil
}

// UpdateAccount implements AccountService
func (as *accountService) UpdateAccount(account dto.AccountDTO) error {
	dataAccount, err := as.accRepo.GetAccountById(account.ID)
	if err != nil {
		return err
	}
	// check if user id in the account is the same as the user id in the token
	if account.UserID != dataAccount.UserID {
		return errors.New("you are not authorized to delete this subcategory")
	}

	err = as.accRepo.UpdateAccount(account)
	if err != nil {
		return err
	}
	return nil
}

// CreateAccount implements AccountService
func (as *accountService) CreateAccount(account dto.AccountDTO) error {
	err := as.accRepo.CreateAccount(account)
	if err != nil {
		return err
	}
	return nil
}

func NewAccountService(accRepo accountRepository.AccountRepository) AccountService {
	return &accountService{
		accRepo: accRepo,
	}
}
