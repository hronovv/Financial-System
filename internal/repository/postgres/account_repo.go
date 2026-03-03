package postgres

import (
	"context"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) CreateAccount(account *domain.Account) error {
	query := `
		INSERT INTO accounts (user_id, bank_id, account_number, balance, is_blocked)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		account.UserID,
		account.BankID,
		account.AccountNumber,
		account.Balance,
		account.IsBlocked,
	).Scan(&account.ID, &account.CreatedAt)

	return err
}
