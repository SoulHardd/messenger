package domain

import "errors"

var (
	ErrDatabase              = errors.New("database error")
	ErrUserNotFound          = errors.New("user not found")
	ErrPhoneExists           = errors.New("phone is already exists")
	ErrLoginExists           = errors.New("login is already exists")
	ErrIncorrectPassword     = errors.New("incorrect password")
	ErrSessionNotFound       = errors.New("session not found")
	ErrInvalidToken          = errors.New("invalid token")
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrInvalidPhone          = errors.New("invalid phone")
	ErrInvalidLogin          = errors.New("invalid login")
	ErrInvalidPassword       = errors.New("invalid password")
)
