package transactionAccService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/accountRepository"
	"dompet-miniprojectalta/repository/transactionAccRepository"
	"errors"
	"fmt"
)

type TransactionAccService interface {
	DeleteTransactionAccount(id, userID uint) error
	CreateTransactionAccount(transAcc dto.TransactionAccount) error
}

type transactionAccService struct {
	transAccRepo transactionAccRepository.TransactionAccRepository
	accountRepo  accountRepository.AccountRepository
}

// DeleteTransactionAccount implements TransactionAccService
func (tas *transactionAccService) DeleteTransactionAccount(transAccID uint, userID uint) error {
	// get old transaction
	transAcc, err := tas.transAccRepo.GetTransactionAccountById(transAccID)
	if err != nil {
		return err
	}
	// check if user id in the transaction is the same as the user id in the token
	if transAcc.UserID != userID {
		return errors.New("you are not authorized to delete this transaction")
	}

	// get data account
	accountFrom, errFrom := tas.accountRepo.GetAccountById(transAcc.AccountFromID)
	accountTo, errTo := tas.accountRepo.GetAccountById(transAcc.AccountToID)
	if errFrom != nil || errTo != nil {
		return errors.New(fmt.Sprint(errFrom, errTo))
	}

	// check if balance account to is more than transaction amount
	if accountTo.Balance < transAcc.Amount {
		return errors.New("Not enough balance in account to")
	}

	// Update balance
	accountFrom.Balance += (transAcc.Amount + transAcc.AdminFee)
	accountTo.Balance -= transAcc.Amount

	// call repository to delete
	err = tas.transAccRepo.DeleteTransactionAccount(transAccID, accountFrom, accountTo)
	if err != nil {
		return err
	}
	return nil
}

// CreateTransactionAccount implements TransactionAccService
func (tas *transactionAccService) CreateTransactionAccount(transAcc dto.TransactionAccount) error {
	// check if accoount from and account to is same
	if transAcc.AccountFromID == transAcc.AccountToID {
		return errors.New("account from and account to is same")
	}

	// get data account
	accountFrom, errFrom := tas.accountRepo.GetAccountById(transAcc.AccountFromID)
	accountTo, errTo := tas.accountRepo.GetAccountById(transAcc.AccountToID)
	if errFrom != nil || errTo != nil {
		return errors.New(fmt.Sprint(errFrom, errTo))
	}
	// check if user id in the account is the same as the user id in the transaction
	if accountFrom.UserID != transAcc.UserID || accountTo.UserID != transAcc.UserID {
		return errors.New("you are not authorized to use this account")
	}

	// check if balance is enough
	if accountFrom.Balance < transAcc.Amount {
		return errors.New("Not enough balance")
	}

	// update balance
	accountFrom.Balance -= (transAcc.Amount + transAcc.AdminFee)
	accountTo.Balance += transAcc.Amount

	// call repository to create transaction account
	err := tas.transAccRepo.CreateTransactionAccount(transAcc, accountFrom, accountTo)
	if err != nil {
		return err
	}
	// success
	return nil
}

func NewTransAccService(transAccRepo transactionAccRepository.TransactionAccRepository, accountRepo accountRepository.AccountRepository) TransactionAccService {
	return &transactionAccService{
		transAccRepo: transAccRepo,
		accountRepo:  accountRepo,
	}
}
