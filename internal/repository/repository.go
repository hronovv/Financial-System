package repository

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	SetUserActive(id int, active bool) error
}

type BankRepository interface {
	GetAllBanks() ([]domain.Bank, error)
}

type EnterpriseRepository interface {
	GetAllEnterprises() ([]domain.Enterprise, error)
	GetEnterpriseByID(id int) (*domain.Enterprise, error)
	GetEnterprisesWithEmployees() ([]domain.EnterpriseWithEmployees, error)
	AddEmployee(enterpriseID, userID int) error
	RemoveEmployee(enterpriseID, userID int) error
	IsEmployee(enterpriseID, userID int) (bool, error)
}

type SalaryApplicationRepository interface {
	Create(app *domain.SalaryApplication) error
	GetByID(id int) (*domain.SalaryApplication, error)
	UpdateStatus(id int, status string) error
	RejectPendingByUserAndEnterprise(userID, enterpriseID int) error
	PaySalary(applicationID int, toAccountID *int, toDepositID *int) error
}

type AccountRepository interface {
	CreateAccount(account *domain.Account) error
	GetAccountByID(id int) (*domain.Account, error)
	GetAccountsByUserID(userID int) ([]domain.Account, error)
	SetAccountBlocked(id int, blocked bool) error
	TransferAccountToAccount(userID, fromAccountID, toAccountID int, amount float64) error
	TransferAccountToDeposit(userID, fromAccountID, toDepositID int, amount float64) error
	GetAccountHistory(accountID int) ([]domain.Transaction, error)
}

type DepositRepository interface {
	CreateDeposit(deposit *domain.Deposit) error
	GetDepositByID(id int) (*domain.Deposit, error)
	SetDepositBlocked(id int, blocked bool) error
	TransferDepositToAccount(userID, fromDepositID, toAccountID int, amount float64) error
	TransferDepositToDeposit(userID, fromDepositID, toDepositID int, amount float64) error
}

type Repositories struct {
	User               UserRepository
	Bank               BankRepository
	Enterprise         EnterpriseRepository
	Account            AccountRepository
	Deposit            DepositRepository
	SalaryApplication  SalaryApplicationRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:              postgres.NewUserRepo(db),
		Bank:              postgres.NewBankRepo(db),
		Enterprise:        postgres.NewEnterpriseRepo(db),
		Account:           postgres.NewAccountRepo(db),
		Deposit:           postgres.NewDepositRepo(db),
		SalaryApplication: postgres.NewSalaryApplicationRepo(db),
	}
}