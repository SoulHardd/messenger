package domain

import "errors"

var (
	ErrDatabase                  = errors.New("database error")
	ErrChatNotFound              = errors.New("chat not found")
	ErrParticipantNotFound       = errors.New("participant not found")
	ErrParticipantAlreadyExists  = errors.New("participant already exists")
	ErrStateAlreadyExists        = errors.New("chat state already exists")
	ErrChatOrParticipantNotFound = errors.New("chat or participant not found")
	ErrPermissionDenied          = errors.New("permission denied")
	ErrUserNotFound              = errors.New("user not found")
	ErrPrivateChatExists         = errors.New("private chat between these users already exists")
	ErrInvalidRole               = errors.New("invalid role")
	ErrInvalidInput              = errors.New("invalid input data")
)
