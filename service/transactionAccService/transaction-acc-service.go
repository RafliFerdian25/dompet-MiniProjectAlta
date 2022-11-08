package transactionAccService

import (
	"dompet-miniprojectalta/constant/constantError"
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/accountRepository"
	"dompet-miniprojectalta/repository/transactionAccRepository"
	"errors"
)

type TransactionAccService interface {
	GetTransactionAccount(userId uint, month int) ([]dto.GetTransactionAccountDTO, error)
	DeleteTransactionAccount(id, userID uint) error
	CreateTransactionAccount(transAcc dto.TransactionAccount) error
}

type transactionAccService struct {
	transAccRepo transactionAccRepository.TransactionAccRepository
	accountRepo  accountRepository.AccountRepository
}

// GetTransactionAccount implements TransactionAccService
func (tas *transactionAccService) GetTransactionAccount(userId uint, month int) ([]dto.GetTransactionAccountDTO, error) {
	// call repository to get transaction
	transactionAccounts, err := tas.transAccRepo.GetTransactionAccount(userId, month)
	if err != nil {
		return []dto.GetTransactionAccountDTO{}, err
	}
	return transactionAccounts, nil
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
		return errors.New(constantError.ErrorNotAuthorized)
	}

	// get data account
	accountFrom, errFrom := tas.accountRepo.GetAccountById(transAcc.AccountFromID)
	if errFrom != nil {
		return errFrom
	}
	accountTo, errTo := tas.accountRepo.GetAccountById(transAcc.AccountToID)
	if errTo != nil {
		return errTo
	}

	// check if balance account to is more than transaction amount
	if accountTo.Balance < transAcc.Amount {
		return errors.New(constantError.ErrorRecipientAccountNotEnoughBalance)
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
		return errors.New(constantError.ErrorSenderAndRecipientAccountIsSame)
	}

	// get data account
	accountFrom, errFrom := tas.accountRepo.GetAccountById(transAcc.AccountFromID)
	if errFrom != nil {
		return errFrom
	}
	accountTo, errTo := tas.accountRepo.GetAccountById(transAcc.AccountToID)
	if errTo != nil {
		return errTo
	}
	// check if user id in the account is the same as the user id in the transaction
	if accountFrom.UserID != transAcc.UserID || accountTo.UserID != transAcc.UserID {
		return errors.New(constantError.ErrorNotAuthorized)
	}

	// check if balance is enough
	if accountFrom.Balance < transAcc.Amount {
		return errors.New(constantError.ErrorAccountNotEnoughBalance)
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
