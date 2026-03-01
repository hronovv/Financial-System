package service

import (
	"financial_system/internal/repository"
)

type Services struct {
	// Auth AuthService 
}


func NewServices(deps *repository.Repositories) *Services {
	return &Services{
		// Auth: NewAuthService(deps.User),
	}
}