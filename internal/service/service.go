package service

import (
	"time"

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
	GetAccountHistory(userID, accountID int) ([]domain.Transaction, error)
}

type DepositService interface {
	OpenDeposit(userID, bankID int, interestRate float64) (*domain.Deposit, error)
	CloseDeposit(userID, depositID int) error
}

type ManagerService interface {
	ApproveUser(userID int) error
	GetUserHistory(userID int) ([]domain.Transaction, error)
}

type Services struct {
	Auth        AuthService
	Bank        BankService
	Enterprise  EnterpriseService
	Account     AccountService
	Deposit     DepositService
	Manager     ManagerService
}

func NewServices(deps *repository.Repositories, jwtSecret string, jwtExpire time.Duration) *Services {
	return &Services{
		Auth:       NewAuthService(deps.User, jwtSecret, jwtExpire),
		Bank:       NewBankService(deps.Bank),
		Enterprise: NewEnterpriseService(deps.Enterprise),
		Account:    NewAccountService(deps.Account),
		Deposit:    NewDepositService(deps.Deposit),
		Manager:    NewManagerService(deps.User, deps.Account),
	}
}
