package main

import (
	"context"
	"financial_system/internal/config"
	"financial_system/internal/repository/postgres"
	"financial_system/pkg/database"
	"log"
)

func main() {
	cfg := config.MustLoad()
	dbPool, err := database.NewPostgresClient(context.Background(), database.ConnectionInfo{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
	})
	if err != nil {
		log.Fatalf("Failed to init db: %s", err)
	}
	defer dbPool.Close()

	log.Println("Connected to db!")
	repos := postgres.NewRepository(dbPool)
	
}
