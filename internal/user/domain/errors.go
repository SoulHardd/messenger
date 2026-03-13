package domain

import "errors"

var (
	ErrDatabase        = errors.New("database error")
	ErrProfileNotFound = errors.New("profile not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidLogin    = errors.New("invalid login")
	ErrInvalidQuery    = errors.New("invalid search query")
	ErrInvalidLimit    = errors.New("invalid limit")
	ErrEmptyUpdate     = errors.New("empty update")
)
