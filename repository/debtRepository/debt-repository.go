package debtRepository

import "dompet-miniprojectalta/models/dto"

type DebtRepostory interface {
	GetDebtById(id uint) (dto.Debt, error)
	CreateDebt(debt dto.Debt, transaction dto.TransactionDTO, account dto.AccountDTO) error
}