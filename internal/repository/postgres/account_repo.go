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

func (r *AccountRepo) GetAccountByID(id int) (*domain.Account, error) {
	query := `
		SELECT id, user_id, bank_id, account_number, balance, is_blocked, created_at
		FROM accounts
		WHERE id = $1
	`

	var acc domain.Account

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.BankID,
		&acc.AccountNumber,
		&acc.Balance,
		&acc.IsBlocked,
		&acc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func (r *AccountRepo) SetAccountBlocked(id int, blocked bool) error {
	query := `
		UPDATE accounts
		SET is_blocked = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(context.Background(), query, blocked, id)
	return err
}

