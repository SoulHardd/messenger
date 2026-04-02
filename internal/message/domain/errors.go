package domain

import "errors"

var (
	ErrDatabase         = errors.New("database error")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidInput     = errors.New("invalid input data")
	ErrChatNotFound     = errors.New("chat not found")
	ErrMessageNotFound  = errors.New("message not found")
)
