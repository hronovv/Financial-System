package service

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Deposit struct {
	repo repository.DepositRepository
}

func NewDepositService(repo repository.DepositRepository) *Deposit {
	return &Deposit{repo: repo}
}

func (s *Deposit) OpenDeposit(userID, bankID int, interestRate float64) (*domain.Deposit, error) {
	deposit := &domain.Deposit{
		UserID:       userID,
		BankID:       bankID,
		Balance:      0,
		InterestRate: interestRate,
		IsBlocked:    false,
	}
	if err := s.repo.CreateDeposit(deposit); err != nil {
		return nil, err
	}
	return deposit, nil
}

func (s *Deposit) CloseDeposit(userID, depositID int) error {
	d, err := s.repo.GetDepositByID(depositID)
	if err != nil {
		return err
	}
	if d.UserID != userID {
		return domain.ErrForbidden
	}
	if d.IsBlocked {
		return domain.ErrDepositAlreadyClosed
	}
	if d.Balance != 0 {
		return domain.ErrDepositHasNonZeroBalance
	}
	return s.repo.SetDepositBlocked(depositID, true)
}

// TransferFromDeposit переводит средства с вклада на счёт или другой вклад того же пользователя.
func (s *Deposit) TransferFromDeposit(userID, fromDepositID int, toAccountID, toDepositID *int, amount float64) error {
	hasAccount := toAccountID != nil && *toAccountID > 0
	hasDeposit := toDepositID != nil && *toDepositID > 0
	if hasAccount && hasDeposit {
		return domain.ErrInvalidTransferTarget
	}
	if !hasAccount && !hasDeposit {
		return domain.ErrInvalidTransferTarget
	}
	if hasAccount {
		return s.repo.TransferDepositToAccount(userID, fromDepositID, *toAccountID, amount)
	}
	return s.repo.TransferDepositToDeposit(userID, fromDepositID, *toDepositID, amount)
}
