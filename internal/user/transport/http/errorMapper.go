package http

import (
	"D/Go/messenger/internal/platform/httpx"
	"D/Go/messenger/internal/user/domain"
	"errors"
	"fmt"
)

type ErrorMapper struct{}

func (ErrorMapper) MapDomainError(err error) httpx.HTTPError {

	switch {

	case errors.Is(err, domain.ErrDatabase):
		return httpx.ErrInternalServer

	case errors.Is(err, domain.ErrProfileNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrProfileNotFound),
		}

	case errors.Is(err, domain.ErrUserNotFound):
		return httpx.HTTPError{
			Code:    404,
			Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrUserNotFound),
		}

	default:
		return httpx.ErrInternalServer
	}
}
