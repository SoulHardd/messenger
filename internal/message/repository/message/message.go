package message

import (
	"D/Go/messenger/internal/message/domain"
	"D/Go/messenger/internal/message/repository"
	"D/Go/messenger/internal/message/service"
	"context"
	"errors"
	"time"

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

func (r *Repository) WithTx(tx pgx.Tx) service.MessageRepository {
	return &Repository{
		q: tx,
	}
}

func (r *Repository) Create(ctx context.Context, msg *domain.Message) (int64, error) {
	query, args, err := squirrel.
		Insert("messages").
		Columns("chat_id", "sender_id", "content").
		Values(msg.ChatId, msg.SenderId, msg.Content).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, domain.ErrDatabase
	}

	var id int64
	var createdAt time.Time
	err = r.q.QueryRow(ctx, query, args...).Scan(&id, &createdAt)
	if err != nil {
		return 0, domain.ErrDatabase
	}

	msg.Id = id
	msg.CreatedAt = createdAt

	return id, nil
}

func (r *Repository) GetChatMessages(ctx context.Context, chatId int64, userId int64, limit int, cursor *domain.Cursor) ([]domain.Message, error) {
	builder := squirrel.
		Select(
			"m.id",
			"m.chat_id",
			"m.sender_id",
			"m.content",
			"m.is_edited",
			"m.created_at",
		).
		Column(
			"COALESCE(ms.status, CASE WHEN m.sender_id = ? THEN ? ELSE ? END) AS status",
			userId,
			int16(domain.StatusRead),
			int16(domain.StatusSent),
		).
		From("messages m").
		LeftJoin("messages_status ms ON ms.message_id = m.id AND ms.user_id = ?", userId).
		Where(squirrel.Eq{"m.chat_id": chatId}).
		OrderBy("m.created_at DESC", "m.id DESC").
		Limit(uint64(limit + 1)).
		PlaceholderFormat(squirrel.Dollar)

	if cursor != nil {
		builder = builder.Where(
			squirrel.Expr(`
				(m.created_at < ?)
				OR (
					m.created_at = ?
					AND m.id < ?
				)
			`, cursor.Time, cursor.Time, cursor.Id),
		)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, domain.ErrDatabase
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		return nil, domain.ErrDatabase
	}
	defer rows.Close()

	messages := make([]domain.Message, 0, limit+1)

	for rows.Next() {
		var msg MessageDB

		err := rows.Scan(
			&msg.Id,
			&msg.ChatId,
			&msg.SenderId,
			&msg.Content,
			&msg.IsEdited,
			&msg.CreatedAt,
			&msg.Status,
		)
		if err != nil {
			return nil, domain.ErrDatabase
		}

		messages = append(messages, domain.Message{
			Id:        msg.Id,
			ChatId:    msg.ChatId,
			SenderId:  msg.SenderId,
			Content:   msg.Content,
			Status:    domain.MessageStatus(msg.Status),
			IsEdited:  msg.IsEdited,
			CreatedAt: msg.CreatedAt,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, domain.ErrDatabase
	}

	return messages, nil
}

func (r *Repository) MessageExists(ctx context.Context, chatId int64, messageId int64) (bool, error) {
	query, args, err := squirrel.
		Select("id").
		From("messages").
		Where(squirrel.And{
			squirrel.Eq{"id": messageId},
			squirrel.Eq{"chat_id": chatId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return false, domain.ErrDatabase
	}

	var id int64
	err = r.q.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, domain.ErrDatabase
	}

	return true, nil
}
