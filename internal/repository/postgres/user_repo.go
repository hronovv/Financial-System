package postgres

import (
	"context"
	"financial_system/internal/domain"

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

// Заглушка, чтобы интерфейс не ругался
func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}
