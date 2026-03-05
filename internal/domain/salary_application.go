package domain

import "time"

const (
	SalaryApplicationStatusPending  = "pending"
	SalaryApplicationStatusApproved = "approved"
	SalaryApplicationStatusRejected  = "rejected"
)

type SalaryApplication struct {
	ID           int        `json:"id" db:"id"`
	UserID       int        `json:"user_id" db:"user_id"`
	EnterpriseID int        `json:"enterprise_id" db:"enterprise_id"`
	Amount       float64    `json:"amount" db:"amount"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	PaidAt       *time.Time `json:"paid_at,omitempty" db:"paid_at"`
}
