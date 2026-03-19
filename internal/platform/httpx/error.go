package httpx

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrInvalidJSON    = HTTPError{Code: 400, Message: `{"error": "invalid JSON"}`}
	ErrInvalidQuery   = HTTPError{Code: 400, Message: `{"error": "invalid query parameter"}`}
	ErrInternalServer = HTTPError{Code: 500, Message: `{"error": "internal server error"}`}
	ErrUnauthorized   = HTTPError{Code: 401, Message: `{"error": "unauthorized"}`}
)
