package transactionRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type TransactionRepository interface {
	DeleteTransaction(id uint, account dto.AccountDTO) error
	UpdateTransaction(newTransaction dto.TransactionDTO, oldTransaction dto.TransactionJoin, newAccount, oldAccount dto.AccountDTO) error
	GetTransactionById(id uint) (dto.TransactionJoin, error)
	CreateTransaction(transaction dto.TransactionDTO, categoryId uint, account dto.AccountDTO) error
}