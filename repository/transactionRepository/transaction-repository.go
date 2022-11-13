package transactionRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type TransactionRepository interface {
	GetTransaction(userId, categoryID uint, month int) ([]dto.GetTransactionDTO, error)
	DeleteTransaction(id uint, account dto.AccountDTO) error
	UpdateTransaction(newTransaction dto.TransactionDTO, oldAccountID uint, newAccountBalance, oldAccountBalance float64) error
	GetTransactionById(id uint) (dto.TransactionJoin, error)
	CreateTransaction(transaction dto.TransactionDTO, categoryID uint, account dto.AccountDTO) error
}