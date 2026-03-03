package service

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Enterprise struct {
	repo repository.EnterpriseRepository
}

func NewEnterpriseService(repo repository.EnterpriseRepository) *Enterprise {
	return &Enterprise{repo: repo}
}

func (s *Enterprise) GetEnterprises() ([]domain.Enterprise, error) {
	return s.repo.GetAllEnterprises()
}

