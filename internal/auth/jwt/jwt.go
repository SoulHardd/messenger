package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	secret    string
	accessTTL time.Duration
}

func New(secret string, accessTTL time.Duration) *JWTService {
	return &JWTService{
		secret:    secret,
		accessTTL: accessTTL,
	}
}

func (j *JWTService) GenerateAccess(userId int64) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(j.accessTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secret))
}

func (j *JWTService) GenerateRefresh() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
