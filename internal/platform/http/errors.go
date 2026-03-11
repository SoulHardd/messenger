package http

import (
	"D/Go/messenger/internal/auth/domain"
	"errors"
	"fmt"
)

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrInvalidJSON    = HTTPError{Code: 400, Message: `{"error": "Invalid JSON"}`}
	ErrInternalServer = HTTPError{Code: 500, Message: `{"error": "Internal server error"}`}
)

func MapDomainError(err error) HTTPError {
	switch {
	case errors.Is(err, domain.ErrDatabase):
		return ErrInternalServer

	case errors.Is(err, domain.ErrUserNotFound):
		return HTTPError{Code: 404, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrUserNotFound)}
	case errors.Is(err, domain.ErrPhoneExists):
		return HTTPError{Code: 409, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrPhoneExists)}
	case errors.Is(err, domain.ErrLoginExists):
		return HTTPError{Code: 409, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrLoginExists)}
	case errors.Is(err, domain.ErrIncorrectPassword):
		return HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrIncorrectPassword)}
	case errors.Is(err, domain.ErrInvalidLogin):
		return HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidLogin)}
	case errors.Is(err, domain.ErrInvalidPassword):
		return HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidPassword)}
	case errors.Is(err, domain.ErrInvalidPhone):
		return HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidPhone)}
	case errors.Is(err, domain.ErrMissingRequiredFields):
		return HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrMissingRequiredFields)}

	case errors.Is(err, domain.ErrSessionNotFound):
		return HTTPError{Code: 404, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrSessionNotFound)}
	case errors.Is(err, domain.ErrInvalidToken):
		return HTTPError{Code: 400, Message: fmt.Sprintf(`{"error": "%s"}`, domain.ErrInvalidToken)}

	default:
		return ErrInternalServer
	}
}
