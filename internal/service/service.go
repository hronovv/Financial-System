package service

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type AuthService interface {
	SignUp(email, password string) error
	SignIn(email, password string) (string, error)
}

type BankService interface {
	GetBanks() ([]domain.Bank, error)
}

type EnterpriseService interface {
	GetEnterprises() ([]domain.Enterprise, error)
}

type AccountService interface {
	OpenAccount(userID, bankID int) (*domain.Account, error)
	CloseAccount(userID, accountID int) error
	TransferFromAccount(userID, fromAccountID int, toAccountID, toDepositID *int, amount float64) error
}

type Services struct {
	Auth        AuthService
	Bank        BankService
	Enterprise  EnterpriseService
	Account     AccountService
}

func NewServices(deps *repository.Repositories) *Services {
	return &Services{
		Auth:       NewAuthService(deps.User),
		Bank:       NewBankService(deps.Bank),
		Enterprise: NewEnterpriseService(deps.Enterprise),
		Account:    NewAccountService(deps.Account),
	}
}
