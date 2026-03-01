package domain

const (
	RoleClient  = "client"
	RoleManager = "manager"
	RoleAdmin   = "admin"
)

type User struct {
	ID           int    `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"` 
	Role         string `json:"role" db:"role"`
	IsActive     bool   `json:"is_active" db:"is_active"` 
}