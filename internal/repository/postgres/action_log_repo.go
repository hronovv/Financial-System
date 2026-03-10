package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ActionLogRepo struct {
	db *pgxpool.Pool
}

func NewActionLogRepo(db *pgxpool.Pool) *ActionLogRepo {
	return &ActionLogRepo{db: db}
}

func (r *ActionLogRepo) Create(log *domain.ActionLog) error {
	ctx := context.Background()

	var details any
	if len(log.Details) != 0 {
		if err := json.Unmarshal(log.Details, &details); err != nil {
			details = json.RawMessage(log.Details)
		}
	}

	query := `
		INSERT INTO action_logs (user_id, action_type, details, is_undone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query, log.UserID, log.Action, details, log.IsUndone).
		Scan(&log.ID, &log.CreatedAt)
	return err
}

func (r *ActionLogRepo) GetAll() ([]domain.ActionLog, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, action_type, details, is_undone, created_at
		FROM action_logs
		ORDER BY created_at DESC, id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.ActionLog

	for rows.Next() {
		var l domain.ActionLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.Details, &l.IsUndone, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *ActionLogRepo) GetByIDForUpdate(ctx context.Context, id int) (*domain.ActionLog, error) {
	var l domain.ActionLog
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, action_type, details, is_undone, created_at
		FROM action_logs
		WHERE id = $1
		FOR UPDATE
	`, id).Scan(&l.ID, &l.UserID, &l.Action, &l.Details, &l.IsUndone, &l.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *ActionLogRepo) MarkUndone(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `UPDATE action_logs SET is_undone = TRUE WHERE id = $1`, id)
	return err
}



