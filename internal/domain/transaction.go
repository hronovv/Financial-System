package domain

import "time"

type Transaction struct {
	ID            int        `json:"id" db:"id"`
	FromAccountID *int       `json:"from_account_id" db:"from_account_id"`
	FromDepositID *int       `json:"from_deposit_id" db:"from_deposit_id"`
	ToAccountID   *int       `json:"to_account_id" db:"to_account_id"`
	ToDepositID   *int       `json:"to_deposit_id" db:"to_deposit_id"`
	Amount        float64    `json:"amount" db:"amount"`
	Type          string     `json:"transaction_type" db:"transaction_type"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

