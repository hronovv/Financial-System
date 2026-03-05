package postgres

import (
	"context"
	"errors"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5"
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

func (r *AccountRepo) GetAccountsByUserID(userID int) ([]domain.Account, error) {
	query := `
		SELECT id, user_id, bank_id, account_number, balance, is_blocked, created_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY id
	`
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.Account
	for rows.Next() {
		var acc domain.Account
		if err := rows.Scan(
			&acc.ID,
			&acc.UserID,
			&acc.BankID,
			&acc.AccountNumber,
			&acc.Balance,
			&acc.IsBlocked,
			&acc.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
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

func (r *AccountRepo) TransferAccountToAccount(userID, fromAccountID, toAccountID int, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Lock source account
	var fromUserID int
	var fromBalance float64
	var fromBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, balance, is_blocked FROM accounts WHERE id = $1 FOR UPDATE`,
		fromAccountID,
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
		return domain.ErrAccountBlocked
	}
	if fromBalance < amount {
		return domain.ErrInsufficientFunds
	}

	// Lock destination account
	var toUserID int
	var toBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, is_blocked FROM accounts WHERE id = $1 FOR UPDATE`,
		toAccountID,
	).Scan(&toUserID, &toBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	if toUserID != userID {
		return domain.ErrForbidden
	}
	if toBlocked {
		return domain.ErrAccountBlocked
	}

	// Update balances
	if _, err := tx.Exec(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromAccountID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toAccountID); err != nil {
		return err
	}

	// Log transaction
	if _, err := tx.Exec(ctx,
		`INSERT INTO transactions (from_account_id, to_account_id, amount, transaction_type) VALUES ($1, $2, $3, $4)`,
		fromAccountID, toAccountID, amount, "account_to_account",
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *AccountRepo) TransferAccountToDeposit(userID, fromAccountID, toDepositID int, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Lock source account
	var fromUserID int
	var fromBalance float64
	var fromBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, balance, is_blocked FROM accounts WHERE id = $1 FOR UPDATE`,
		fromAccountID,
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
		return domain.ErrAccountBlocked
	}
	if fromBalance < amount {
		return domain.ErrInsufficientFunds
	}

	// Lock destination deposit
	var toUserID int
	var toBlocked bool
	err = tx.QueryRow(ctx,
		`SELECT user_id, is_blocked FROM deposits WHERE id = $1 FOR UPDATE`,
		toDepositID,
	).Scan(&toUserID, &toBlocked)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	if toUserID != userID {
		return domain.ErrForbidden
	}
	if toBlocked {
		return domain.ErrDepositBlocked
	}

	// Update balances
	if _, err := tx.Exec(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromAccountID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE deposits SET balance = balance + $1 WHERE id = $2`, amount, toDepositID); err != nil {
		return err
	}

	// Log transaction
	if _, err := tx.Exec(ctx,
		`INSERT INTO transactions (from_account_id, to_deposit_id, amount, transaction_type) VALUES ($1, $2, $3, $4)`,
		fromAccountID, toDepositID, amount, "account_to_deposit",
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *AccountRepo) GetAccountHistory(accountID int) ([]domain.Transaction, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx, `
		SELECT id, from_account_id, from_deposit_id, to_account_id, to_deposit_id, amount, transaction_type, created_at
		FROM transactions
		WHERE from_account_id = $1 OR to_account_id = $1
		ORDER BY created_at DESC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []domain.Transaction

	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.FromAccountID,
			&t.FromDepositID,
			&t.ToAccountID,
			&t.ToDepositID,
			&t.Amount,
			&t.Type,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		history = append(history, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}


