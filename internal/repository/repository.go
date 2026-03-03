package repository

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByEmail(email string) (*domain.User, error)
}

type BankRepository interface {
	GetAllBanks() ([]domain.Bank, error)
}

type EnterpriseRepository interface {
	GetAllEnterprises() ([]domain.Enterprise, error)
}

type AccountRepository interface {
	CreateAccount(account *domain.Account) error
	GetAccountByID(id int) (*domain.Account, error)
	SetAccountBlocked(id int, blocked bool) error
	TransferAccountToAccount(userID, fromAccountID, toAccountID int, amount float64) error
	TransferAccountToDeposit(userID, fromAccountID, toDepositID int, amount float64) error
}

type Repositories struct {
	User        UserRepository
	Bank        BankRepository
	Enterprise  EnterpriseRepository
	Account     AccountRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:       postgres.NewUserRepo(db),
		Bank:       postgres.NewBankRepo(db),
		Enterprise: postgres.NewEnterpriseRepo(db),
		Account:    postgres.NewAccountRepo(db),
	}
}