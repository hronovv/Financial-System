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

// OpenAccount создаёт счёт в указанном банке.
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

// CloseAccount блокирует счёт. Баланс должен быть нулевым.
func (s *Account) CloseAccount(userID, accountID int) error {
	acc, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	if acc.UserID != userID {
		return domain.ErrForbidden
	}

	if acc.IsBlocked {
		return domain.ErrAccountAlreadyClosed
	}

	if acc.Balance != 0 {
		return domain.ErrAccountHasNonZeroBalance
	}

	return s.repo.SetAccountBlocked(accountID, true)
}

// TransferFromAccount переводит средства со счёта на счёт или вклад того же пользователя.
func (s *Account) TransferFromAccount(userID, fromAccountID int, toAccountID, toDepositID *int, amount float64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	hasToAccount := toAccountID != nil && *toAccountID > 0
	hasToDeposit := toDepositID != nil && *toDepositID > 0

	if hasToAccount == hasToDeposit {
		return domain.ErrInvalidTransferTarget
	}

	if hasToAccount {
		return s.repo.TransferAccountToAccount(userID, fromAccountID, *toAccountID, amount)
	}

	return s.repo.TransferAccountToDeposit(userID, fromAccountID, *toDepositID, amount)
}

// GetAccountHistory возвращает историю операций по счёту.
func (s *Account) GetAccountHistory(userID, accountID int) ([]domain.Transaction, error) {
	acc, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	if acc.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return s.repo.GetAccountHistory(accountID)
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

