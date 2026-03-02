package main

import (
	"context"
	"financial_system/internal/config"
	"financial_system/internal/repository"
	"financial_system/internal/service"
	"financial_system/internal/transport/rest"
	"financial_system/pkg/database"
	"log"
	"net/http"
)
// @title           Financial System API
// @version         1.0
// @description     API для управления финансами
// @host            localhost:8080
// @BasePath        /

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
	repos := repository.NewRepositories(dbPool)

	services := service.NewServices(repos)

	handlers := rest.NewHandler(services)


	srv := &http.Server{
		Addr:    ":" + cfg.HTTPServer.Address, 
		Handler: handlers.InitRoutes(), 
	}

	log.Printf("Сервер запущен на порту %s", cfg.HTTPServer.Address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Ошибка при запуске сервера: %s", err)
	}
}
