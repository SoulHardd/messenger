package httpx

type DomainErrorMapper interface {
	MapDomainError(err error) HTTPError
}
