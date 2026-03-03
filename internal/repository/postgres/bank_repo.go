package postgres

import (
	"context"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BankRepo struct {
	db *pgxpool.Pool
}

func NewBankRepo(db *pgxpool.Pool) *BankRepo {
	return &BankRepo{db: db}
}

func (r *BankRepo) GetAllBanks() ([]domain.Bank, error) {
	query := `
		SELECT id, name
		FROM banks
		ORDER BY id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banks []domain.Bank

	for rows.Next() {
		var b domain.Bank
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return nil, err
		}
		banks = append(banks, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

