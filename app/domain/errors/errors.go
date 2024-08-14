package errors

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidParameter    = errors.New("invalid parameter")
	ErrUserNotFound        = errors.New("user not found")
	ErrPartnerError        = errors.New("partner error")
)
