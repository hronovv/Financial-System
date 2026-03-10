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
	TransferFromDeposit(userID, fromDepositID int, toAccountID, toDepositID *int, amount float64) error
}

type ManagerService interface {
	ApproveUser(userID int) error
	GetUserHistory(userID int) ([]domain.Transaction, error)
	BlockAccount(accountID int) error
	UnblockAccount(accountID int) error
	BlockDeposit(depositID int) error
	UnblockDeposit(depositID int) error
	GetEnterprisesWithEmployees() ([]domain.EnterpriseWithEmployees, error)
	AddEmployee(enterpriseID, userID int) error
	RemoveEmployee(enterpriseID, userID int) error
	ApproveSalaryApplication(applicationID int) error
}

type SalaryProjectService interface {
	ApplyForSalaryProject(userID, enterpriseID int, amount float64) (*domain.SalaryApplication, error)
	ReceiveSalary(userID, applicationID int, toAccountID, toDepositID *int) error
}

type AuditLogger interface {
	LogAction(userID *int, action string, details any) error
}

type AdminService interface {
	GetAllLogs() ([]domain.ActionLog, error)
	UndoAction(logID int, deps *repository.Repositories) error
	UndoAllActions(deps *repository.Repositories) error
}

type Services struct {
	Auth          AuthService
	Bank          BankService
	Enterprise    EnterpriseService
	Account       AccountService
	Deposit       DepositService
	Manager       ManagerService
	SalaryProject SalaryProjectService
	Audit         AuditLogger
	Admin         AdminService
	Repositories  *repository.Repositories
}

func NewServices(deps *repository.Repositories, jwtSecret string, jwtExpire time.Duration) *Services {
	return &Services{
		Auth:          NewAuthService(deps.User, jwtSecret, jwtExpire),
		Bank:          NewBankService(deps.Bank),
		Enterprise:    NewEnterpriseService(deps.Enterprise),
		Account:       NewAccountService(deps.Account),
		Deposit:       NewDepositService(deps.Deposit),
		Manager:       NewManagerService(deps.User, deps.Account, deps.Deposit, deps.Enterprise, deps.SalaryApplication),
		SalaryProject: NewSalaryProjectService(deps.Enterprise, deps.SalaryApplication),
		Audit:         NewAuditLogger(deps.ActionLog),
		Admin:         NewAdminService(deps.ActionLog),
		Repositories:  deps,
	}
}
