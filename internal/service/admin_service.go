package service

import (
	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Admin struct {
	logRepo repository.ActionLogRepository
}

func NewAdminService(logRepo repository.ActionLogRepository) *Admin {
	return &Admin{logRepo: logRepo}
}

// GetAllLogs возвращает все записи action_logs в порядке убывания даты.
func (a *Admin) GetAllLogs() ([]domain.ActionLog, error) {
	return a.logRepo.GetAll()
}