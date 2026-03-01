package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	// User domain.UserRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		// User: postgres.NewUserRepo(db),
	}
}
