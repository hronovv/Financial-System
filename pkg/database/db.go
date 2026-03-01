package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionInfo struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewPostgresClient(ctx context.Context, info ConnectionInfo) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		info.User, info.Password, info.Host, info.Port, info.DBName)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Ping database error: %w", err)
	}

	return pool, nil
}
