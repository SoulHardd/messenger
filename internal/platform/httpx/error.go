package httpx

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrInvalidJSON    = HTTPError{Code: 400, Message: `{"error": "Invalid JSON"}`}
	ErrInternalServer = HTTPError{Code: 500, Message: `{"error": "Internal server error"}`}
)
