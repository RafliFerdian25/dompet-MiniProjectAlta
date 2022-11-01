package accountRepository

import "dompet-miniprojectalta/models/dto"

type AccountRepository interface {
	// GetAccountByUser(userId uint) ([]dto.AccountDTO, error)
	// DeleteAccount(id uint) error
	GetAccountById(id uint) (dto.AccountDTO, error)
	UpdateAccount(account dto.AccountDTO) error
	CreateAccount(account dto.AccountDTO) error
}