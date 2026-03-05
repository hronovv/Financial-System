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
	ErrDepositAlreadyClosed     = errors.New("deposit is already closed")
	ErrDepositHasNonZeroBalance = errors.New("deposit has non-zero balance")
	ErrNotFound                 = errors.New("not found")
	ErrUserAlreadyActive        = errors.New("user is already active")
	ErrCanOnlyApproveClient     = errors.New("can only approve users with role client")
	ErrNotEmployee              = errors.New("user is not an employee of this enterprise")
	ErrApplicationNotPending    = errors.New("application is not pending")
	ErrApplicationNotApproved   = errors.New("application is not approved or already paid")
	ErrApplicationAlreadyPaid   = errors.New("salary for this application was already paid")
	ErrInsufficientEnterpriseBalance = errors.New("enterprise has insufficient balance")
)