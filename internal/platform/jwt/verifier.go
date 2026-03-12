package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type Verifier struct {
	secret string
}

func NewVerifier(secret string) *Verifier {
	return &Verifier{
		secret: secret,
	}
}

func (v *Verifier) ParseAccess(tokenStr string) (int64, error) {
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(v.secret), nil
		},
	)

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || token.Valid {
		return 0, errors.New("invalid token")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid token")
	}

	return int64(userId), nil
}
