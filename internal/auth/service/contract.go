package service

import (
	"D/Go/messenger/internal/auth/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (int64, error)
	GetOneByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetOneByLogin(ctx context.Context, login string) (*domain.User, error)
	Delete(ctx context.Context, userId int64) error
	Update(ctx context.Context, user *domain.User) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error)
	DeleteByRefreshTokenHash(ctx context.Context, hash string) error
	DeleteByUserID(ctx context.Context, userId int64) error
	DeleteExpired(ctx context.Context) error
}

type JWTService interface {
	GenerateAccess(userId int64) (string, error)
	GenerateRefresh() (string, error)
}
