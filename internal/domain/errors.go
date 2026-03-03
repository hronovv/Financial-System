package domain

import "errors"

var (
	ErrForbidden                = errors.New("forbidden")
	ErrAccountAlreadyClosed     = errors.New("account is already closed")
	ErrAccountHasNonZeroBalance = errors.New("account has non-zero balance")
)