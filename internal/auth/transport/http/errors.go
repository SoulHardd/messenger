package http

import (
	"D/Go/messenger/internal/auth/domain"
	"D/Go/messenger/internal/platform/httpx"
	"errors"
	"fmt"
)

type ErrorMapper struct{}

func (ErrorMapper) MapDomainError(err error) httpx.HTTPError {
	switch {
	case errors.Is(err, domain.ErrDatabase):
		return httpx.ErrInternalServer

	case errors.Is(err, domain.ErrUserNotFound):
		return httpx.HTTPError{Code: 404, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrUserNotFound)}
	case errors.Is(err, domain.ErrPhoneExists):
		return httpx.HTTPError{Code: 409, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrPhoneExists)}
	case errors.Is(err, domain.ErrLoginExists):
		return httpx.HTTPError{Code: 409, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrLoginExists)}
	case errors.Is(err, domain.ErrIncorrectPassword):
		return httpx.HTTPError{Code: 401, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrIncorrectPassword)}
	case errors.Is(err, domain.ErrInvalidLogin):
		return httpx.HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidLogin)}
	case errors.Is(err, domain.ErrInvalidPassword):
		return httpx.HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidPassword)}
	case errors.Is(err, domain.ErrInvalidPhone):
		return httpx.HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidPhone)}
	case errors.Is(err, domain.ErrMissingRequiredFields):
		return httpx.HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrMissingRequiredFields)}

	case errors.Is(err, domain.ErrSessionNotFound):
		return httpx.HTTPError{Code: 404, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrSessionNotFound)}
	case errors.Is(err, domain.ErrInvalidToken):
		return httpx.HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidToken)}

	default:
		return httpx.ErrInternalServer
	}
}
