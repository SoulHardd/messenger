package session

import (
	"D/Go/messenger/internal/auth/domain"
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool       *pgxpool.Pool
	refreshTTL time.Duration
}

func New(pool *pgxpool.Pool, refreshTTL time.Duration) *Repository {
	return &Repository{
		pool:       pool,
		refreshTTL: refreshTTL,
	}
}

func (r *Repository) Create(ctx context.Context, session *domain.Session) error {
	sessionDB := SessionDB{
		UserId:           session.UserId,
		RefreshTokenHash: session.RefreshTokenHash,
		ExpiresAt:        time.Now().UTC().Add(r.refreshTTL),
	}

	query, args, err := squirrel.
		Insert("sessions").
		Columns("user_id", "refresh_token_hash", "expires_at").
		Values(sessionDB.UserId, sessionDB.RefreshTokenHash, sessionDB.ExpiresAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	_, err = r.pool.Exec(ctx, query, args...)

	if err != nil {
		return domain.ErrDatabase
	}
	return nil
}

func (r *Repository) GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error) {
	query, args, err := squirrel.
		Select("id", "user_id", "refresh_token_hash", "expires_at").
		From("sessions").
		Where(squirrel.Eq{"refresh_token_hash": hash}).
		Where("expires_at > NOW()").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	var session SessionDB
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&session.Id,
		&session.UserId,
		&session.RefreshTokenHash,
		&session.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.Session{
		Id:               session.Id,
		UserId:           session.UserId,
		RefreshTokenHash: session.RefreshTokenHash,
	}, nil
}

func (r *Repository) DeleteByRefreshTokenHash(ctx context.Context, hash string) error {

	query, args, err := squirrel.
		Delete("sessions").
		Where(squirrel.Eq{"refresh_token_hash": hash}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}

	if result.RowsAffected() == 0 {
		return domain.ErrSessionNotFound
	}

	return nil
}

func (r *Repository) DeleteByUserID(ctx context.Context, userId int64) error {

	query, args, err := squirrel.
		Delete("sessions").
		Where(squirrel.Eq{"user_id": userId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}

	return nil
}

func (r *Repository) DeleteExpired(ctx context.Context) error {

	query, args, err := squirrel.
		Delete("sessions").
		Where("expires_at < NOW()").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}

	return nil
}
