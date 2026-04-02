package status

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

func (r *Repository) WithTx(tx pgx.Tx) service.MessageStatusRepository {
	return &Repository{
		q: tx,
	}
}

func (r *Repository) InitForMessage(ctx context.Context, messageId int64, senderId int64, participantIds []int64) error {
	if len(participantIds) == 0 {
		return nil
	}

	builder := squirrel.
		Insert("messages_status").
		Columns("message_id", "user_id", "status")

	for _, userId := range participantIds {
		status := int16(domain.StatusSent)
		if userId == senderId {
			status = int16(domain.StatusRead)
		}
		builder = builder.Values(messageId, userId, status)
	}

	query, args, err := builder.
		Suffix("ON CONFLICT (message_id, user_id) DO UPDATE SET status = EXCLUDED.status, updated_at = now()").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.ErrDatabase
	}

	if _, err = r.q.Exec(ctx, query, args...); err != nil {
		return domain.ErrDatabase
	}

	return nil
}

func (r *Repository) MarkDelivered(ctx context.Context, chatId int64, userId int64) error {
	query := `
		UPDATE messages_status AS ms
		SET status = $1, updated_at = now()
		FROM messages AS m
		WHERE ms.message_id = m.id
		  AND ms.user_id = $2
		  AND m.chat_id = $3
		  AND m.sender_id <> $2
		  AND ms.status < $1
	`
	if _, err := r.q.Exec(ctx, query, int16(domain.StatusDelivered), userId, chatId); err != nil {
		return domain.ErrDatabase
	}

	return nil
}

func (r *Repository) MarkReadUpTo(ctx context.Context, chatId int64, userId int64, messageId int64) error {
	query := `
		UPDATE messages_status AS ms
		SET status = $1, updated_at = now()
		FROM messages AS m
		WHERE ms.message_id = m.id
		  AND ms.user_id = $2
		  AND m.chat_id = $3
		  AND m.id <= $4
		  AND ms.status < $1
	`
	if _, err := r.q.Exec(ctx, query, int16(domain.StatusRead), userId, chatId, messageId); err != nil {
		return domain.ErrDatabase
	}

	return nil
}
