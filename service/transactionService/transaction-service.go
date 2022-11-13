package transactionService

import (
	"dompet-miniprojectalta/constant/constantCategory"
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/accountRepository"
	"dompet-miniprojectalta/repository/subCategoryRepository"
	"dompet-miniprojectalta/repository/transactionRepository"
	"errors"
	"time"
)

type TransactionService interface {
	GetTransaction(userId uint, month int) (map[string]interface{}, error)
	DeleteTransaction(id uint, userID uint) error
	UpdateTransaction(transaction dto.TransactionDTO, userId uint) error
	CreateTransaction(transaction dto.TransactionDTO) error
}

type transactionService struct {
	transactionRepo transactionRepository.TransactionRepository
	accountRepo     accountRepository.AccountRepository
	subCategoryRepo subCategoryRepository.SubCategoryRepository
}

// GetTransaction implements TransactionService
func (ts *transactionService) GetTransaction(userId uint, month int) (map[string]interface{}, error) {
	// call repository to get transaction
	expenseTransactions, err := ts.transactionRepo.GetTransaction(userId, constantCategory.ExpenseCategory, month)
	if err != nil {
		return map[string]interface{}{}, err
	}
	incomeTransactions, err := ts.transactionRepo.GetTransaction(userId, constantCategory.IncomeCategory, month)
	if err != nil {
		return map[string]interface{}{}, err
	}
	data := map[string]interface{}{
		"expense": expenseTransactions,
		"income": incomeTransactions,
		"month_transaction": time.Month(month).String(),
	}
	return data, nil
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
		return errors.New(constantError.ErrorNotAuthorized)
	}

	// get account of transaction
	account, errGetOldAccount := ts.accountRepo.GetAccountById(transaction.AccountID)
	if errGetOldAccount != nil {
		return errGetOldAccount
	}
	account.Balance -= transaction.Amount

	// call repository to delete
	err = ts.transactionRepo.DeleteTransaction(id, account)
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
		return errors.New(constantError.ErrorNotAuthorized)
	}

	// Get data subcategpry
	newSubCategory, err := ts.subCategoryRepo.GetSubCategoryById(newTransaction.SubCategoryID)
	if newTransaction.SubCategoryID != 0 {
		if err != nil {
			return err
		}
		// check if user id in the subcategory is the same as the user id in the transaction
		if *newSubCategory.UserID != userID && *newSubCategory.UserID != 0 {
			return errors.New(constantError.ErrorNotAuthorized)
		}
	}

	// check if category change
	if newSubCategory.CategoryID != oldTransaction.CategoryID {
		return errors.New(constantError.ErrorCannotChangeSubCategory)
	}

	// get data account
	oldAccount, errGetOldAccount := ts.accountRepo.GetAccountById(oldTransaction.AccountID)
	if errGetOldAccount != nil {
		return errGetOldAccount
	}
	var newAccount dto.AccountDTO
	if newTransaction.AccountID != 0 && newTransaction.AccountID != oldTransaction.AccountID {
		newAccount, err = ts.accountRepo.GetAccountById(newTransaction.AccountID)
		if err != nil {
			return err
		}
		// check if user id in the account is the same as the user id in the transaction
		if newAccount.UserID != userID {
			return errors.New(constantError.ErrorNotAuthorized)
		}
	} else {
		newAccount, err = ts.accountRepo.GetAccountById(oldTransaction.AccountID)
		if err != nil {
			return err
		}
		// update balance new account
		newAccount.Balance -= oldTransaction.Amount
	}

	// update balance old account
	oldAccount.Balance -= oldTransaction.Amount

	// check if category id is expense
	if newSubCategory.CategoryID == 2 {
		// check if balance is enough
		if oldTransaction.AccountID == newAccount.ID {
			if newTransaction.Amount > newAccount.Balance+(oldTransaction.Amount*-1) {
				return errors.New(constantError.ErrorOldAccountBalanceNotEnough)
			}
		} else {
			if newTransaction.Amount > newAccount.Balance {
				return errors.New(constantError.ErrorNewAccountBalanceNotEnough)
			}
		}
		newTransaction.Amount *= -1
	}

	// update balance new account
	newAccount.Balance += newTransaction.Amount

	// call repository to save transaction
	err = ts.transactionRepo.UpdateTransaction(newTransaction, oldAccount.ID, newAccount.Balance, oldAccount.Balance)
	if err != nil {
		return err
	}
	return nil
}

// CreateTransaction implements TransactionService
func (ts *transactionService) CreateTransaction(transaction dto.TransactionDTO) error {
	// Get data subcategpry
	subCategory, err := ts.subCategoryRepo.GetSubCategoryById(transaction.SubCategoryID)
	if err != nil {
		return err
	}
	// check if user id in the subcategory is the same as the user id in the transaction
	if *subCategory.UserID != transaction.UserID && *subCategory.UserID != 0 {
		return errors.New(constantError.ErrorNotAuthorized)
	}

	// get data account
	account, err := ts.accountRepo.GetAccountById(transaction.AccountID)
	if err != nil {
		return err
	}
	// check if user id in the account is the same as the user id in the transaction
	if account.UserID != transaction.UserID {
		return errors.New(constantError.ErrorNotAuthorized)
	}

	// check if category id is 1 (debt & loan) or 2 (expense)
	if subCategory.CategoryID == 1 {
		return errors.New(constantError.ErrorCannotUseCategory)
	} else if subCategory.CategoryID == 2 {
		// check if balance is enough
		if transaction.Amount > account.Balance {
			return errors.New(constantError.ErrorAccountNotEnoughBalance)
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

func NewTransactionService(transactionRepository transactionRepository.TransactionRepository,
	accountRepo accountRepository.AccountRepository,
	subCategoryRepo subCategoryRepository.SubCategoryRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepository,
		accountRepo:     accountRepo,
		subCategoryRepo: subCategoryRepo,
	}
}
