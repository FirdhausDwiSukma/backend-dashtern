package domain

import "errors"

var (
	ErrUserNotFound    = errors.New("USER_NOT_FOUND")
	ErrInvalidPassword = errors.New("INVALID_PASSWORD")
)
