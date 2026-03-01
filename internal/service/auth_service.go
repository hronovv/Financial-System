package service

import (
	"errors"
	"financial_system/internal/domain"
	"financial_system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *Auth {
	return &Auth{repo: repo}
}

func (s *Auth) SignUp(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("не удалось обработать пароль")
	}

	user := domain.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         domain.RoleClient, 
		IsActive:     false,            
	}

	return s.repo.CreateUser(&user)
}

func (s *Auth) SignIn(email, password string) (string, error) {
	return "", nil
}
