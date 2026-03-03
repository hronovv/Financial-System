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

type Services struct {
	Auth AuthService
	Bank BankService
}

func NewServices(deps *repository.Repositories) *Services {
	return &Services{
		Auth: NewAuthService(deps.User),
		Bank: NewBankService(deps.Bank),
	}
}
