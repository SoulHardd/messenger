package profile

import (
	"D/Go/messenger/internal/platform/pointers"
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

func (r *Repository) GetByUserId(ctx context.Context, userId int64) (*domain.Profile, error) {
	query, args, err := squirrel.
		Select("user_id", "nickname", "bio", "avatar_url").
		From("user_profiles").
		Where(squirrel.Eq{"user_id": userId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	var p ProfileDB

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&p.UserId,
		&p.Nickname,
		&p.Bio,
		&p.AvatarURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProfileNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.Profile{
		UserId:    p.UserId,
		Nickname:  pointers.ZeroIfNil(p.Nickname),
		Bio:       pointers.ZeroIfNil(p.Bio),
		AvatarURL: pointers.ZeroIfNil(p.AvatarURL),
	}, nil
}

func (r *Repository) Update(ctx context.Context, userId int64, up *domain.UpdateProfile) error {
	builder := squirrel.
		Update("user_profiles").
		Where(squirrel.Eq{"user_id": userId})

	if up.Nickname != nil {
		builder = builder.Set("nickname", *up.Nickname)
	}
	if up.Bio != nil {
		builder = builder.Set("bio", *up.Bio)
	}
	if up.AvatarURL != nil {
		builder = builder.Set("avatar_url", *up.AvatarURL)
	}

	query, args, err := builder.
		Set("updated_at", squirrel.Expr("now()")).
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

func (r *Repository) Search(ctx context.Context, queryStr string, limit int) ([]domain.Profile, error) {
	query, args, err := squirrel.
		Select(
			"u.id",
			"p.nickname",
			"p.bio",
			"p.avatar_url",
			"u.login",
		).
		From("users u").
		LeftJoin("user_profiles p ON p.user_id = u.id").
		Where("u.login ILIKE ?", "%"+queryStr+"%").
		OrderBy("u.login").
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, domain.ErrDatabase
	}
	defer rows.Close()

	var profiles []domain.Profile

	for rows.Next() {
		var p ProfileDB
		err := rows.Scan(
			&p.UserId,
			&p.Nickname,
			&p.Bio,
			&p.AvatarURL,
			&p.Login,
		)
		if err != nil {
			return nil, domain.ErrDatabase
		}
		domainProfile := domain.Profile{
			UserId:    p.UserId,
			Nickname:  pointers.ZeroIfNil(p.Nickname),
			Bio:       pointers.ZeroIfNil(p.Bio),
			AvatarURL: pointers.ZeroIfNil(p.AvatarURL),
			Login:     pointers.ZeroIfNil(p.Login),
		}
		profiles = append(profiles, domainProfile)
	}

	return profiles, nil
}
