package chat

import (
	"D/Go/messenger/internal/message/domain"
	"D/Go/messenger/internal/message/repository"
	"D/Go/messenger/internal/message/service"
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q repository.Queryer
}

func New(p *pgxpool.Pool) *Repository {
	return &Repository{q: p}
}

func (r *Repository) WithTx(tx pgx.Tx) service.ChatStateRepository {
	return &Repository{
		q: tx,
	}
}

func (r *Repository) UpdateLastMessage(ctx context.Context, chatId int64, msg domain.Message) error {
	query, args, err := squirrel.
		Update("chat_state").
		Set("last_message_id", msg.Id).
		Set("last_message_text", msg.Content).
		Set("last_message_at", msg.CreatedAt).
		Where(squirrel.Eq{"chat_id": chatId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	result, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}
	if result.RowsAffected() == 0 {
		return domain.ErrChatNotFound
	}

	return nil
}

func (r *Repository) IncrementUnread(ctx context.Context, chatId int64, excludeUserId int64) error {
	query, args, err := squirrel.
		Update("chat_state").
		Set("unread_count", squirrel.Expr("unread_count + 1")).
		Where(squirrel.And{
			squirrel.Eq{"chat_id": chatId},
			squirrel.NotEq{"user_id": excludeUserId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	result, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}
	if result.RowsAffected() == 0 {
		return domain.ErrChatNotFound
	}

	return nil
}

func (r *Repository) ResetUnread(ctx context.Context, chatId int64, userId int64) error {
	query, args, err := squirrel.
		Update("chat_state").
		Set("unread_count", 0).
		Set("last_read_message_id", squirrel.Expr("last_message_id")).
		Where(squirrel.And{
			squirrel.Eq{"chat_id": chatId},
			squirrel.Eq{"user_id": userId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	result, err := r.q.Exec(ctx, query, args...)
	if err != nil {
		return domain.ErrDatabase
	}
	if result.RowsAffected() == 0 {
		return domain.ErrChatNotFound
	}

	return nil
}
