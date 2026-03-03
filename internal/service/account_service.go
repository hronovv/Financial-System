package service

import (
	"crypto/rand"
	"fmt"

	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Account struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *Account {
	return &Account{repo: repo}
}

// OpenAccount создает новый счет для пользователя в банке.
// TODO: userID должен браться из JWT, а не из тела запроса.
func (s *Account) OpenAccount(userID, bankID int) (*domain.Account, error) {
	accountNumber, err := generateAccountNumber()
	if err != nil {
		return nil, err
	}

	account := &domain.Account{
		UserID:        userID,
		BankID:        bankID,
		AccountNumber: accountNumber,
		Balance:       0,
		IsBlocked:     false,
	}

	if err := s.repo.CreateAccount(account); err != nil {
		return nil, err
	}

	return account, nil
}

func generateAccountNumber() (string, error) {
	const length = 16
	const digits = "0123456789"

	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	for i := range buf {
		buf[i] = digits[int(buf[i])%len(digits)]
	}

	return fmt.Sprintf("%s", string(buf)), nil
}

