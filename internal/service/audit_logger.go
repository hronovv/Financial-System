package service

import (
	"encoding/json"

	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Audit struct {
	repo repository.ActionLogRepository
}

func NewAuditLogger(repo repository.ActionLogRepository) *Audit {
	return &Audit{repo: repo}
}

// LogAction writes an audit record to action_logs. details is serialized to JSON.
func (a *Audit) LogAction(userID *int, action string, details any) error {
	var raw json.RawMessage
	if details != nil {
		b, err := json.Marshal(details)
		if err != nil {
			return err
		}
		raw = b
	}

	log := &domain.ActionLog{
		UserID:   userID,
		Action:   action,
		Details:  raw,
		IsUndone: false,
	}
	return a.repo.Create(log)
}

