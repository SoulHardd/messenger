package http

import (
	"D/Go/messenger/internal/auth/domain"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, ur *domain.UserRaw, field domain.UserField) (*domain.Tokens, error)
	Register(ctx context.Context, ur *domain.UserRaw) (*domain.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.Tokens, error)
}
