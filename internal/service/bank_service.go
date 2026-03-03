package service

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Bank struct {
	repo repository.BankRepository
}

func NewBankService(repo repository.BankRepository) *Bank {
	return &Bank{repo: repo}
}

func (s *Bank) GetBanks() ([]domain.Bank, error) {
	return s.repo.GetAllBanks()
}