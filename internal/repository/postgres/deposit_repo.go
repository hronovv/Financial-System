package postgres

import (
	"context"
	"errors"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
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

func (r *DepositRepo) TransferDepositToAccount(userID, fromDepositID, toAccountID int, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var fromUserID int
	var fromBalance float64
	var fromBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, balance, is_blocked FROM deposits WHERE id = $1 FOR UPDATE`,
		fromDepositID,
	).Scan(&fromUserID, &fromBalance, &fromBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	if fromUserID != userID {
		return domain.ErrForbidden
	}
	if fromBlocked {
		return domain.ErrDepositBlocked
	}
	if fromBalance < amount {
		return domain.ErrInsufficientFunds
	}

	var toUserID int
	var toBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, is_blocked FROM accounts WHERE id = $1 FOR UPDATE`,
		toAccountID,
	).Scan(&toUserID, &toBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrAccountNotFound
		}
		return err
	}
	if toUserID != userID {
		return domain.ErrForbidden
	}
	if toBlocked {
		return domain.ErrAccountBlocked
	}

	if _, err := tx.Exec(ctx, `UPDATE deposits SET balance = balance - $1 WHERE id = $2`, amount, fromDepositID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toAccountID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO transactions (from_deposit_id, to_account_id, amount, transaction_type) VALUES ($1, $2, $3, $4)`,
		fromDepositID, toAccountID, amount, "deposit_to_account",
	); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *DepositRepo) TransferDepositToDeposit(userID, fromDepositID, toDepositID int, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var fromUserID int
	var fromBalance float64
	var fromBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, balance, is_blocked FROM deposits WHERE id = $1 FOR UPDATE`,
		fromDepositID,
	).Scan(&fromUserID, &fromBalance, &fromBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	if fromUserID != userID {
		return domain.ErrForbidden
	}
	if fromBlocked {
		return domain.ErrDepositBlocked
	}
	if fromBalance < amount {
		return domain.ErrInsufficientFunds
	}

	var toUserID int
	var toBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, is_blocked FROM deposits WHERE id = $1 FOR UPDATE`,
		toDepositID,
	).Scan(&toUserID, &toBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrDepositNotFound
		}
		return err
	}
	if toUserID != userID {
		return domain.ErrForbidden
	}
	if toBlocked {
		return domain.ErrDepositBlocked
	}

	if _, err := tx.Exec(ctx, `UPDATE deposits SET balance = balance - $1 WHERE id = $2`, amount, fromDepositID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE deposits SET balance = balance + $1 WHERE id = $2`, amount, toDepositID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO transactions (from_deposit_id, to_deposit_id, amount, transaction_type) VALUES ($1, $2, $3, $4)`,
		fromDepositID, toDepositID, amount, "deposit_to_deposit",
	); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
