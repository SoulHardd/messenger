package auth

type JWTVerifier interface {
	ParseAccess(token string) (int64, error)
}
