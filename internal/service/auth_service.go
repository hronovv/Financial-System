package service

import (
	"errors"
	"time"

	"financial_system/internal/domain"
	"financial_system/internal/repository"
	"financial_system/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	repo   repository.UserRepository
	secret string
	expire time.Duration
}

func NewAuthService(repo repository.UserRepository, secret string, expire time.Duration) *Auth {
	return &Auth{repo: repo, secret: secret, expire: expire}
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
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", domain.ErrUserNotActive
	}

	return jwt.NewToken(s.secret, user.ID, user.Role, s.expire)
}
