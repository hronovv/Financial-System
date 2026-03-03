package domain

type Enterprise struct {
	ID      int     `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	Balance float64 `json:"balance" db:"balance"`
}

