package transactionRepository

import (
	"dompet-miniprojectalta/models/dto"
)

type TransactionRepository interface {
	UpdateTransaction(newTransaction dto.TransactionDTO, oldTransaction dto.TransactionJoin, newAccount dto.AccountDTO) error
	GetTransactionById(id uint) (dto.TransactionJoin, error)
	GetSubCategoryById(id uint) (dto.SubCategoryDTO, error)
	GetAccountById(id uint) (dto.AccountDTO, error)
	CreateTransaction(transaction dto.TransactionDTO, categoryId uint, account dto.AccountDTO) error
}