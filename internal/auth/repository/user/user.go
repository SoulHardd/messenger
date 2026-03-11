package user

import (
	"D/Go/messenger/internal/auth/domain"
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const duplicateKeyCode = "23505"

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Create(ctx context.Context, user *domain.User) (int64, error) {
	userDB := UserDB{
		Phone:        user.Phone,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}
	query, args, err := squirrel.
		Insert("users").
		Columns("phone", "login", "password_hash").
		Values(userDB.Phone, userDB.Login, userDB.PasswordHash).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, domain.ErrDatabase
	}

	var id int64
	err = r.pool.QueryRow(ctx, query, args...).Scan(&id)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == duplicateKeyCode {
			switch pgError.ConstraintName {
			case "users_phone_key":
				return 0, domain.ErrPhoneExists
			case "users_login_key":
				return 0, domain.ErrLoginExists
			}
		}
		return 0, domain.ErrDatabase
	}
	return id, nil
}

func (r *Repository) GetOneByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query, args, err := squirrel.
		Select("id", "phone", "login", "password_hash").
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
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.User{
		Id:           user.Id,
		Phone:        user.Phone,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *Repository) GetOneByLogin(ctx context.Context, login string) (*domain.User, error) {
	query, args, err := squirrel.
		Select("id", "phone", "login", "password_hash").
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
		&user.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.User{
		Id:           user.Id,
		Phone:        user.Phone,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	userDB := UserDB{
		Id:           user.Id,
		Phone:        user.Phone,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}

	query, args, err := squirrel.
		Update("users").
		Set("phone", userDB.Phone).
		Set("login", userDB.Login).
		Set("password_hash", userDB.PasswordHash).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": userDB.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	result, err := r.pool.Exec(ctx, query, args...)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == duplicateKeyCode {
			switch pgError.ConstraintName {
			case "users_phone_key":
				return domain.ErrPhoneExists
			case "users_login_key":
				return domain.ErrLoginExists
			}
		}
		return domain.ErrDatabase
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, userId int64) error {

	query, args, err := squirrel.
		Delete("users").
		Where(squirrel.Eq{"id": userId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	result, err := r.pool.Exec(ctx, query, args...)

	if err != nil {
		return domain.ErrDatabase
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
