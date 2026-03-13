package http

import (
	"D/Go/messenger/internal/user/domain"
	"context"
)

type UserService interface {
	GetMe(ctx context.Context, userId int64) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, userId int64, update *domain.UpdateProfile) error
	GetByLogin(ctx context.Context, login string) (*domain.Profile, error)
	GetByPhone(ctx context.Context, phone string) (*domain.Profile, error)
	Search(ctx context.Context, query string, limit int) ([]domain.Profile, error)
}
