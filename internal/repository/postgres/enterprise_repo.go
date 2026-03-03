package postgres

import (
	"context"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EnterpriseRepo struct {
	db *pgxpool.Pool
}

func NewEnterpriseRepo(db *pgxpool.Pool) *EnterpriseRepo {
	return &EnterpriseRepo{db: db}
}

func (r *EnterpriseRepo) GetAllEnterprises() ([]domain.Enterprise, error) {
	query := `
		SELECT id, name, balance
		FROM enterprises
		ORDER BY id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enterprises []domain.Enterprise

	for rows.Next() {
		var e domain.Enterprise
		if err := rows.Scan(&e.ID, &e.Name, &e.Balance); err != nil {
			return nil, err
		}
		enterprises = append(enterprises, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enterprises, nil
}

