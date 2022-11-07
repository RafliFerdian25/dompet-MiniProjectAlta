package debtRepository

import "dompet-miniprojectalta/models/dto"

type DebtRepostory interface {
	GetDebt(userId uint, subCategory int, debtStatus string) ([]dto.GetDebtTransactionResponse, error)
	DeleteDebt(id uint, account dto.AccountDTO) error
	GetDebtById(id uint) (dto.Debt, error)
	CreateDebt(debt dto.Debt, transaction dto.TransactionDTO, account dto.AccountDTO) error
}