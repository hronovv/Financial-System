package domain

import (
	"encoding/json"
	"time"
)

// ActionLog описывает запись аудита делового действия.
type ActionLog struct {
	ID        int             `json:"id" db:"id"`
	UserID    *int            `json:"user_id,omitempty" db:"user_id"`
	Action    string          `json:"action_type" db:"action_type"`
	Details   json.RawMessage `json:"details,omitempty" db:"details"`
	IsUndone  bool            `json:"is_undone" db:"is_undone"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

