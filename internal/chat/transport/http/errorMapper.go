package http

import (
	"D/Go/messenger/internal/chat/domain"
	"D/Go/messenger/internal/platform/httpx"
	"errors"
	"fmt"
)

type ErrorMapper struct{}

func (ErrorMapper) MapDomainError(err error) httpx.HTTPError {
	switch {

	case errors.Is(err, domain.ErrDatabase):
		return httpx.ErrInternalServer

	case errors.Is(err, domain.ErrChatNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrChatNotFound),
		}

	case errors.Is(err, domain.ErrParticipantNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrParticipantNotFound),
		}

	case errors.Is(err, domain.ErrParticipantAlreadyExists):
		return httpx.HTTPError{
			Code:    409,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrParticipantAlreadyExists),
		}

	case errors.Is(err, domain.ErrChatOrParticipantNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrChatOrParticipantNotFound),
		}

	case errors.Is(err, domain.ErrPrivateChatExists):
		return httpx.HTTPError{
			Code:    409,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrPrivateChatExists),
		}

	case errors.Is(err, domain.ErrPermissionDenied):
		return httpx.HTTPError{
			Code:    403,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrPermissionDenied),
		}

	case errors.Is(err, domain.ErrUserNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrUserNotFound),
		}

	case errors.Is(err, domain.ErrInvalidInput):
		return httpx.HTTPError{
			Code:    400,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidInput),
		}

	case errors.Is(err, domain.ErrInvalidRole):
		return httpx.HTTPError{
			Code:    400,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidRole),
		}

	default:
		return httpx.ErrInternalServer
	}
}
