package transactionService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/transactionRepository"
	"errors"
	"fmt"
)

type TransactionService interface {
	DeleteTransaction(id uint, userID uint) error
	UpdateTransaction(transaction dto.TransactionDTO, userId uint) error
	CreateTransaction(transaction dto.TransactionDTO) error
}

type transactionService struct {
	transactionRepo transactionRepository.TransactionRepository
}

// DeleteTransaction implements TransactionService
func (ts *transactionService) DeleteTransaction(id uint, userID uint) error {
	// get old transaction
	transaction, err := ts.transactionRepo.GetTransactionById(id)
	if err != nil {
		return err
	}
	// check if user id in the transaction is the same as the user id in the token
	if transaction.UserID != userID {
		return errors.New("you are not authorized to delete this transaction")
	}
	// call repository to delete
	err = ts.transactionRepo.DeleteTransaction(id, transaction.AccountID, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTransaction implements TransactionService
func (ts *transactionService) UpdateTransaction(newTransaction dto.TransactionDTO, userID uint) error {
	// get old transaction
	oldTransaction, err := ts.transactionRepo.GetTransactionById(newTransaction.ID)
	if err != nil {
		return err
	}
	// check if user id in the transaction is the same as the user id in the token
	if oldTransaction.UserID != userID {
		return errors.New("you are not authorized to update this transaction")
	}

	// Get data subcategpry
	newSubCategory, err := ts.transactionRepo.GetSubCategoryById(newTransaction.SubCategoryID)
	if newTransaction.SubCategoryID != 0 {
		if err != nil {
			return err
		}
		// check if user id in the subcategory is the same as the user id in the transaction
		if newSubCategory.UserID != userID && newSubCategory.UserID != 0 {
			return errors.New("you are not authorized to use this sub category")
		}
	}

	// check if category change
	if newSubCategory.CategoryID != oldTransaction.CategoryID {
		return errors.New("you cannot change category")
	}

	// get data account
	var account dto.AccountDTO
	if newTransaction.AccountID != 0 && newTransaction.AccountID != oldTransaction.AccountID {
		account, err = ts.transactionRepo.GetAccountById(newTransaction.AccountID)
		if err != nil {
			return err
		}
		// check if user id in the account is the same as the user id in the transaction
		if account.UserID != userID {
			return errors.New("you are not authorized to use this account")
		}
		fmt.Println(account)
		fmt.Println(account.Balance)
	} else {
		account, err = ts.transactionRepo.GetAccountById(oldTransaction.AccountID)
		if err != nil {
			return err
		}
	}

	// check if category id
	if newSubCategory.CategoryID == 2 {
		// check if balance is enough
		if oldTransaction.AccountID == account.ID {
			if newTransaction.Amount > account.Balance+(oldTransaction.Amount*-1) {
				return errors.New("Not enough balance old account transaction")
			}
		} else {
			if newTransaction.Amount > account.Balance {
				return errors.New("Not enough balance new account transaction")
			}
		}
		newTransaction.Amount *= -1
	}
	// call repository to save transaction
	err = ts.transactionRepo.UpdateTransaction(newTransaction, oldTransaction, account)
	if err != nil {
		return err
	}
	return nil
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
