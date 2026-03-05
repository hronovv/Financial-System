package postgres

import (
	"context"
	"errors"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SalaryApplicationRepo struct {
	db *pgxpool.Pool
}

func NewSalaryApplicationRepo(db *pgxpool.Pool) *SalaryApplicationRepo {
	return &SalaryApplicationRepo{db: db}
}

func (r *SalaryApplicationRepo) Create(app *domain.SalaryApplication) error {
	query := `
		INSERT INTO salary_applications (user_id, enterprise_id, amount, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(context.Background(), query,
		app.UserID, app.EnterpriseID, app.Amount, app.Status,
	).Scan(&app.ID, &app.CreatedAt, &app.UpdatedAt)
}

func (r *SalaryApplicationRepo) GetByID(id int) (*domain.SalaryApplication, error) {
	query := `
		SELECT id, user_id, enterprise_id, amount, status, created_at, updated_at, paid_at
		FROM salary_applications WHERE id = $1
	`
	var app domain.SalaryApplication
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&app.ID, &app.UserID, &app.EnterpriseID, &app.Amount, &app.Status,
		&app.CreatedAt, &app.UpdatedAt, &app.PaidAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

func (r *SalaryApplicationRepo) UpdateStatus(id int, status string) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE salary_applications SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, id)
	return err
}

func (r *SalaryApplicationRepo) RejectPendingByUserAndEnterprise(userID, enterpriseID int) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE salary_applications SET status = $1, updated_at = NOW()
		 WHERE user_id = $2 AND enterprise_id = $3 AND status = $4`,
		domain.SalaryApplicationStatusRejected, userID, enterpriseID, domain.SalaryApplicationStatusPending)
	return err
}

func (r *SalaryApplicationRepo) PaySalary(applicationID int, toAccountID *int, toDepositID *int) error {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var app domain.SalaryApplication
	err = tx.QueryRow(ctx,
		`SELECT id, user_id, enterprise_id, amount, status, paid_at
		 FROM salary_applications WHERE id = $1 FOR UPDATE`,
		applicationID,
	).Scan(&app.ID, &app.UserID, &app.EnterpriseID, &app.Amount, &app.Status, &app.PaidAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrNotFound
		}
		return err
	}
	if app.Status != domain.SalaryApplicationStatusApproved {
		return domain.ErrApplicationNotApproved
	}
	if app.PaidAt != nil {
		return domain.ErrApplicationAlreadyPaid
	}

	// Decrement enterprise balance
	res, err := tx.Exec(ctx,
		`UPDATE enterprises SET balance = balance - $1 WHERE id = $2 AND balance >= $1`,
		app.Amount, app.EnterpriseID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return domain.ErrInsufficientEnterpriseBalance
	}

	var toAccID, toDepID *int
	if toAccountID != nil {
		var accUserID int
		err = tx.QueryRow(ctx, `SELECT user_id FROM accounts WHERE id = $1 FOR UPDATE`, *toAccountID).Scan(&accUserID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrNotFound
			}
			return err
		}
		if accUserID != app.UserID {
			return domain.ErrForbidden
		}
		if _, err = tx.Exec(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2`, app.Amount, *toAccountID); err != nil {
			return err
		}
		toAccID = toAccountID
	} else if toDepositID != nil {
		var depUserID int
		err = tx.QueryRow(ctx, `SELECT user_id FROM deposits WHERE id = $1 FOR UPDATE`, *toDepositID).Scan(&depUserID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrNotFound
			}
			return err
		}
		if depUserID != app.UserID {
			return domain.ErrForbidden
		}
		if _, err = tx.Exec(ctx, `UPDATE deposits SET balance = balance + $1 WHERE id = $2`, app.Amount, *toDepositID); err != nil {
			return err
		}
		toDepID = toDepositID
	} else {
		return domain.ErrInvalidTransferTarget
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (from_account_id, from_deposit_id, to_account_id, to_deposit_id, amount, transaction_type)
		 VALUES (NULL, NULL, $1, $2, $3, $4)`,
		toAccID, toDepID, app.Amount, "salary")
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `UPDATE salary_applications SET paid_at = NOW() WHERE id = $1`, applicationID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
