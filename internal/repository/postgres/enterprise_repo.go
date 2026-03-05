package postgres

import (
	"context"
	"errors"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5"
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

func (r *EnterpriseRepo) GetEnterpriseByID(id int) (*domain.Enterprise, error) {
	query := `SELECT id, name, balance FROM enterprises WHERE id = $1`
	var e domain.Enterprise
	err := r.db.QueryRow(context.Background(), query, id).Scan(&e.ID, &e.Name, &e.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *EnterpriseRepo) GetEnterprisesWithEmployees() ([]domain.EnterpriseWithEmployees, error) {
	enterprises, err := r.GetAllEnterprises()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(),
		`SELECT enterprise_id, user_id FROM enterprise_employees ORDER BY enterprise_id, user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	empMap := make(map[int][]int) // enterprise_id -> []user_id
	for rows.Next() {
		var eid, uid int
		if err := rows.Scan(&eid, &uid); err != nil {
			return nil, err
		}
		empMap[eid] = append(empMap[eid], uid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]domain.EnterpriseWithEmployees, 0, len(enterprises))
	for _, e := range enterprises {
		userIDs := empMap[e.ID]
		if userIDs == nil {
			userIDs = []int{}
		}
		result = append(result, domain.EnterpriseWithEmployees{Enterprise: e, EmployeeUserIDs: userIDs})
	}
	return result, nil
}

func (r *EnterpriseRepo) AddEmployee(enterpriseID, userID int) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO enterprise_employees (enterprise_id, user_id) VALUES ($1, $2) ON CONFLICT (enterprise_id, user_id) DO NOTHING`,
		enterpriseID, userID)
	return err
}

func (r *EnterpriseRepo) RemoveEmployee(enterpriseID, userID int) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM enterprise_employees WHERE enterprise_id = $1 AND user_id = $2`,
		enterpriseID, userID)
	return err
}

func (r *EnterpriseRepo) IsEmployee(enterpriseID, userID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM enterprise_employees WHERE enterprise_id = $1 AND user_id = $2)`,
		enterpriseID, userID).Scan(&exists)
	return exists, err
}

