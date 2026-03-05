package postgres

import (
	"context"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DepositRepo struct {
	db *pgxpool.Pool
}

func NewDepositRepo(db *pgxpool.Pool) *DepositRepo {
	return &DepositRepo{db: db}
}

func (r *DepositRepo) CreateDeposit(deposit *domain.Deposit) error {
	query := `
		INSERT INTO deposits (user_id, bank_id, balance, interest_rate, is_blocked)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		deposit.UserID,
		deposit.BankID,
		deposit.Balance,
		deposit.InterestRate,
		deposit.IsBlocked,
	).Scan(&deposit.ID, &deposit.CreatedAt)

	return err
}

func (r *DepositRepo) GetDepositByID(id int) (*domain.Deposit, error) {
	query := `
		SELECT id, user_id, bank_id, balance, interest_rate, is_blocked, created_at
		FROM deposits
		WHERE id = $1
	`

	var d domain.Deposit
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&d.ID,
		&d.UserID,
		&d.BankID,
		&d.Balance,
		&d.InterestRate,
		&d.IsBlocked,
		&d.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DepositRepo) SetDepositBlocked(id int, blocked bool) error {
	query := `
		UPDATE deposits
		SET is_blocked = $1
		WHERE id = $2
	`
	_, err := r.db.Exec(context.Background(), query, blocked, id)
	return err
}
