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

type Repositories struct {
	User UserRepository
	Bank BankRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: postgres.NewUserRepo(db),
		Bank: postgres.NewBankRepo(db),
	}
}

type BankRepository interface {
	GetAllBanks() ([]domain.Bank, error)
}