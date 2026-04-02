package http

import (
	"D/Go/messenger/internal/message/domain"
	"D/Go/messenger/internal/platform/httpx"
	"errors"
	"fmt"
)

type ErrorMapper struct{}

func (ErrorMapper) MapDomainError(err error) httpx.HTTPError {
	switch {
	case errors.Is(err, domain.ErrDatabase):
		return httpx.ErrInternalServer

	case errors.Is(err, domain.ErrPermissionDenied):
		return httpx.HTTPError{
			Code:    403,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrPermissionDenied),
		}

	case errors.Is(err, domain.ErrInvalidInput):
		return httpx.HTTPError{
			Code:    400,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidInput),
		}

	case errors.Is(err, domain.ErrChatNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrChatNotFound),
		}

	case errors.Is(err, domain.ErrMessageNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrMessageNotFound),
		}

	default:
		return httpx.ErrInternalServer
	}
}
