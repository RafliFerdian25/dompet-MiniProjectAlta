package debtService

import (
	"dompet-miniprojectalta/models/dto"
	"dompet-miniprojectalta/repository/accountRepository"
	"dompet-miniprojectalta/repository/debtRepository"
	"dompet-miniprojectalta/repository/subCategoryRepository"
	"errors"
	"fmt"
)

type DebtService interface {
	DeleteDebt(id, userId uint) error
	CreateDebt(debtTransaction dto.DebtTransactionDTO) error
}

type debtService struct {
	debtRepo        debtRepository.DebtRepostory
	accountRepo     accountRepository.AccountRepository
	subCategoryRepo subCategoryRepository.SubCategoryRepository
}

// DeleteDebt implements DebtService
func (ds *debtService) DeleteDebt(id uint, userID uint) error {
	// get old debt
	debt, err := ds.debtRepo.GetDebtById(id)
	if err != nil {
		return err
	}

	// check if user id in the debt is the same as the user id in the token
	if debt.Transactions[0].UserID != userID {
		return errors.New("you are not authorized to delete this debt")
	}

	// get data account
	account, err := ds.accountRepo.GetAccountById(debt.Transactions[0].AccountID)
	if err != nil {
		return err
	}
	if account.Balance - debt.Total < 0  {
		return errors.New("Not enough balance")
	}
	account.Balance -= debt.Total

	// call repository to delete
	err = ds.debtRepo.DeleteDebt(id, account)
	if err != nil {
		return err
	}
	return nil
}

// CreateDebt implements DebtService
func (ds *debtService) CreateDebt(debtTransaction dto.DebtTransactionDTO) error {
	// get data account
	account, err := ds.accountRepo.GetAccountById(debtTransaction.AccountID)
	if err != nil {
		return err
	}
	// check if user id in the account is the same as the user id in the transaction
	if account.UserID != debtTransaction.UserID {
		return errors.New("you are not authorized to use this account")
	}

	// check if the sub category is make expense
	if debtTransaction.SubCategoryID == 2 || debtTransaction.SubCategoryID == 3 {
		// check if balance is enough
		if debtTransaction.Amount > account.Balance {
			return errors.New("Not enough balance")
		}
		debtTransaction.Amount *= -1
	}

	// check if the debtTransaction is the new debt or not
	var debt dto.Debt
	if debtTransaction.DebtID == 0 {
		debt = dto.Debt{
			Name:       debtTransaction.Name,
			Total:      debtTransaction.Amount,
			Remaining:  debtTransaction.Amount,
			Note:       debtTransaction.Note,
			DebtStatus: "unpaid",
		}
		// check if debt or loan
		if debtTransaction.SubCategoryID == 1 {
			debt.Status = "debt"
		} else if debtTransaction.SubCategoryID == 3 {
			debt.Status = "loan"
		}
	} else {
		// get data debt
		debt, err = ds.debtRepo.GetDebtById(debtTransaction.DebtID)
		if err != nil {
			return err
		}
		//
		if debt.Status == "debt" {
			if debtTransaction.SubCategoryID == 3 || debtTransaction.SubCategoryID == 4 {
				return errors.New("you cannot change sub category")
			}
		} else if debt.Status == "loan" {
			if debtTransaction.SubCategoryID == 1 || debtTransaction.SubCategoryID == 2 {
				return errors.New("you cannot change sub category")
			}
		}
		// check if user id in the debt is the same as the user id in the transaction
		if debt.Transactions[0].UserID != debtTransaction.UserID {
			return errors.New("you are not authorized to use this debt")
		}
		// check if the sub category is make total increase or remaining decrease
		if debtTransaction.SubCategoryID == 1 || debtTransaction.SubCategoryID == 4 {
			debt.Total += debtTransaction.Amount
			debt.Remaining += debtTransaction.Amount
		} else if debtTransaction.SubCategoryID == 2 || debtTransaction.SubCategoryID == 3 {
			// check if amount is more than remaining
			if (debtTransaction.Amount * -1) > debt.Remaining {
				return errors.New(fmt.Sprint("Input amount is more than remaining debt. Unpaid amount is ", debt.Remaining))
			}
			debt.Remaining += debtTransaction.Amount
			// check if remaining is 0
			fmt.Println(debt.Remaining)
			if debt.Remaining == 0 {
				debt.DebtStatus = "paid"
			}
		}
	}
	// set account balance
	account.Balance += debtTransaction.Amount

	// set model transaction
	transaction := dto.TransactionDTO{
		UserID:        debtTransaction.UserID,
		DebtID:        debtTransaction.DebtID,
		SubCategoryID: debtTransaction.SubCategoryID,
		AccountID:     debtTransaction.AccountID,
		Amount:        debtTransaction.Amount,
		Note:          debtTransaction.Note,
	}
	// call repository to create debt
	err = ds.debtRepo.CreateDebt(debt, transaction, account)
	if err != nil {
		return err
	}
	return nil
}

// NewDebtService creates a new instance of debtService
func NewDebtService(debtRepo debtRepository.DebtRepostory, accountRepo accountRepository.AccountRepository, subCategoryRepo subCategoryRepository.SubCategoryRepository) DebtService {
	return &debtService{
		debtRepo:        debtRepo,
		accountRepo:     accountRepo,
		subCategoryRepo: subCategoryRepo,
	}
}