package transactionRepository

import "dompet-miniprojectalta/models/dto"

type TransactionRepository interface {
	GetSubCategoryById(id uint) (dto.SubCategoryDTO, error)
	GetAccountById(id uint) (dto.AccountDTO, error)
	CreateTransaction(transaction dto.TransactionDTO, categoryId uint, account dto.AccountDTO) error
}