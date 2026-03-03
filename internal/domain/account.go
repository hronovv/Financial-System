package domain

import "time"

type Account struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	BankID        int       `json:"bank_id" db:"bank_id"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	Balance       float64   `json:"balance" db:"balance"`
	IsBlocked     bool      `json:"is_blocked" db:"is_blocked"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

