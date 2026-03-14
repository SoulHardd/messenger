package user

import (
	"D/Go/messenger/internal/user/domain"
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetById(ctx context.Context, id int64) (*domain.User, error) {
	query, args, err := squirrel.
		Select("id", "phone", "login").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	var user UserDB
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Phone,
		&user.Login,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.User{
		Id:    user.Id,
		Phone: user.Phone,
		Login: user.Login,
	}, nil
}

func (r *Repository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	query, args, err := squirrel.
		Select("id", "phone", "login").
		From("users").
		Where(squirrel.Eq{"login": login}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	var user UserDB
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Phone,
		&user.Login,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.User{
		Id:    user.Id,
		Phone: user.Phone,
		Login: user.Login,
	}, nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query, args, err := squirrel.
		Select("id", "phone", "login").
		From("users").
		Where(squirrel.Eq{"phone": phone}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	var user UserDB
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Phone,
		&user.Login,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.User{
		Id:    user.Id,
		Phone: user.Phone,
		Login: user.Login,
	}, nil
}
