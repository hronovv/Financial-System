package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	// "financial_system/internal/domain" 
)

type Repository struct {
	// User    domain.UserRepository    
	// Account domain.AccountRepository
}


func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		// User:    NewUserRepo(db),    
		// Account: NewAccountRepo(db), 
	}
}
