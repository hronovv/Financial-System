package postgres

import (
	"context"
	"errors"
	"financial_system/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, role, is_active) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id
	`

	err := r.db.QueryRow(context.Background(), query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
	).Scan(&user.ID)

	return err
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active
		FROM users
		WHERE email = $1
	`

	var u domain.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetUserByID(id int) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, role, is_active
		FROM users
		WHERE id = $1
	`
	var u domain.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) SetUserActive(id int, active bool) error {
	query := `UPDATE users SET is_active = $1 WHERE id = $2`
	_, err := r.db.Exec(context.Background(), query, active, id)
	return err
}
