package service

import (
	"D/Go/messenger/internal/user/domain"
	"context"
)

type UserRepository interface {
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	GetByPhone(ctx context.Context, phone string) (*domain.User, error)
}

type ProfileRepository interface {
	GetByUserId(ctx context.Context, userId int64) (*domain.Profile, error)
	Update(ctx context.Context, userId int64, update *domain.UpdateProfile) error
	Search(ctx context.Context, query string, limit int) ([]domain.Profile, error)
}
