package transactionService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/transactionRepository"
	"errors"
)

type TransactionService interface {
	CreateTransaction(transaction dto.TransactionDTO) error
}

type transactionService struct {
	transactionRepo transactionRepository.TransactionRepository
}

// CreateTransaction implements TransactionService
func (ts *transactionService) CreateTransaction(transaction dto.TransactionDTO) error {
	// Get data subcategpry
	subCategory, err := ts.transactionRepo.GetSubCategoryById(transaction.SubCategoryID)
	if err != nil {
		return err
	}
	// check if user id in the subcategory is the same as the user id in the transaction
	if subCategory.UserID != transaction.UserID && subCategory.UserID != 0 {
		return errors.New("you are not authorized to use this sub category")
	}

	// get data account
	account, err := ts.transactionRepo.GetAccountById(transaction.AccountID)
	if err != nil {
		return err
	}

	// check if user id in the account is the same as the user id in the transaction
	if account.UserID != transaction.UserID {
		return errors.New("you are not authorized to use this account")
	}

	// check if category id is 1 (debt & loan) or 2 (expense)
	if subCategory.CategoryID == 1 {
		return errors.New("you are not authorized to use this category in transactions routes")
	} else if subCategory.CategoryID == 2 {
		// check if balance is enough
		if transaction.Amount > account.Balance {
			return errors.New("Not enough balance")
		}
		transaction.Amount *= -1
	}
	// call repository to save transaction
	err = ts.transactionRepo.CreateTransaction(transaction, subCategory.CategoryID, account)
	if err != nil {
		return err
	}
	return nil
}

func NewTransactionService(transactionRepository transactionRepository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepository,
	}
}
