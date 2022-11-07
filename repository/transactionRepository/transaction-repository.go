package transactionRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type TransactionRepository interface {
	GetTransaction(month int, userId, categpryID uint) ([]dto.GetTransactionDTO, error)
	DeleteTransaction(id uint, account dto.AccountDTO) error
	UpdateTransaction(newTransaction dto.TransactionDTO, oldTransaction dto.TransactionJoin, newAccount, oldAccount dto.AccountDTO) error
	GetTransactionById(id uint) (dto.TransactionJoin, error)
	CreateTransaction(transaction dto.TransactionDTO, categoryID uint, account dto.AccountDTO) error
}