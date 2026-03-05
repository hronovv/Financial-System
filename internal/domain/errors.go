package domain

import "errors"

var (
	ErrInvalidCredentials       = errors.New("invalid email or password")
	ErrUserNotActive            = errors.New("user not active")
	ErrForbidden                = errors.New("forbidden")
	ErrAccountAlreadyClosed     = errors.New("account is already closed")
	ErrAccountHasNonZeroBalance = errors.New("account has non-zero balance")
	ErrInvalidAmount            = errors.New("invalid amount")
	ErrInvalidTransferTarget    = errors.New("invalid transfer target")
	ErrInsufficientFunds        = errors.New("insufficient funds")
	ErrAccountBlocked           = errors.New("account is blocked")
	ErrDepositBlocked           = errors.New("deposit is blocked")
	ErrNotFound                 = errors.New("not found")
)