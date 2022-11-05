package transactionAccRepository

import "dompet-miniprojectalta/models/dto"

type TransactionAccRepository interface {
	DeleteTransactionAccount(id uint, accountFrom, accountTo dto.AccountDTO) error
	GetTransactionAccountById(id uint) (dto.TransactionAccount, error)
	CreateTransactionAccount(transAcc dto.TransactionAccount, accountFrom, accountTo dto.AccountDTO) error
}