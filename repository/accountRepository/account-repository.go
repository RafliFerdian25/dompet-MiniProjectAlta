package accountRepository

import "dompet-miniprojectalta/models/dto"

type AccountRepository interface {
	DeleteAccount(id uint) error
	GetAccountByUser(userId uint) ([]dto.AccountDTO, error)
	GetAccountById(id uint) (dto.AccountDTO, error)
	UpdateAccount(account dto.AccountDTO) error
	CreateAccount(account dto.AccountDTO) error
}