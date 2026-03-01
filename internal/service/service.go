package service

import (
	"financial_system/internal/repository"
)

type AuthService interface {
	SignUp(email, password string) error
	SignIn(email, password string) (string, error)
}

type Services struct {
	Auth AuthService
}

func NewServices(deps *repository.Repositories) *Services {
	return &Services{
		Auth: NewAuthService(deps.User),
	}
}
