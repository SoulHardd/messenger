package chat

import (
	"D/Go/messenger/internal/chat/domain"
	"D/Go/messenger/internal/chat/repository"
	"D/Go/messenger/internal/chat/service"
	"D/Go/messenger/internal/platform/pointers"
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const duplicateKeyCode = "23505"

type Repository struct {
	q repository.Queryer
}

func New(q repository.Queryer) *Repository {
	return &Repository{q: q}
}

func (r *Repository) WithTx(tx pgx.Tx) service.ChatRepository {
	return &Repository{
		q: tx,
	}
}

func (r *Repository) Create(ctx context.Context, chat *domain.Chat) (int64, error) {
	chatDB := ChatDB{
		Type:      int(chat.Type),
		Title:     chat.Title,
		CreatedBy: chat.OwnerId,
	}
	query, args, err := squirrel.
		Insert("chats").
		Columns("type", "created_by", "title").
		Values(chatDB.Type, chatDB.CreatedBy, chatDB.Title).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return 0, domain.ErrDatabase
	}

	var id int64
	err = r.q.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, domain.ErrDatabase
	}

	return id, nil
}

func (r *Repository) GetUserChat(ctx context.Context, userId int64, chatId int64) (*domain.Chat, error) {
	var chatDB ChatDB
	var state ChatStateDB

	query, args, err := squirrel.
		Select(
			"c.id",
			"c.type",
			"c.created_by",
			"c.title",
			"s.last_read_message_id",
			"s.unread_count",
			"s.last_message_id",
			"s.last_message_text",
			"s.last_message_at",
		).
		From("chat_state s").
		Join("chats c ON c.id = s.chat_id").
		Where(squirrel.And{
			squirrel.Eq{"s.user_id": userId},
			squirrel.Eq{"s.chat_id": chatId},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, domain.ErrDatabase
	}

	err = r.q.QueryRow(ctx, query, args...).Scan(
		&chatDB.Id,
		&chatDB.Type,
		&chatDB.CreatedBy,
		&chatDB.Title,
		&state.LastReadMsgId,
		&state.UnreadCount,
		&state.LastMsgId,
		&state.LastMsgText,
		&state.LastMsgAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrChatNotFound
		}
		return nil, domain.ErrDatabase
	}

	return &domain.Chat{
		Id:            chatDB.Id,
		Type:          domain.ChatType(chatDB.Type),
		Title:         chatDB.Title,
		OwnerId:       chatDB.CreatedBy,
		LastReadMsgId: pointers.ZeroIfNil(state.LastReadMsgId),
		UnreadCount:   pointers.ZeroIfNil(state.UnreadCount),
		LastMsgId:     pointers.ZeroIfNil(state.LastMsgId),
		LastMsgText:   pointers.ZeroIfNil(state.LastMsgText),
		LastMsgTime:   pointers.ZeroIfNil(state.LastMsgAt),
	}, nil
}

func (r *Repository) GetUserChats(ctx context.Context, userId int64, limit int, cursor *domain.Cursor) ([]domain.Chat, error) {
	builder := squirrel.
		Select(
			"c.id",
			"c.type",
			"c.created_by",
			"c.title",
			"s.last_read_message_id",
			"s.unread_count",
			"s.last_message_id",
			"s.last_message_text",
			"COALESCE(s.last_message_at, c.created_at)",
		).
		From("chat_state s").
		Join("chats c ON c.id = s.chat_id").
		Where(squirrel.Eq{"s.user_id": userId}).
		OrderBy("COALESCE(s.last_message_at, c.created_at) DESC", "c.id DESC").
		Limit(uint64(limit + 1)).
		PlaceholderFormat(squirrel.Dollar)

	if cursor != nil {
		builder = builder.Where(
			squirrel.Expr(`
				(COALESCE(s.last_message_at, c.created_at) < ?)
				OR (
					COALESCE(s.last_message_at, c.created_at) = ?
					AND c.id < ?
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

	chats := make([]domain.Chat, 0, limit+1)

	for rows.Next() {
		var chat ChatDB
		var state ChatStateDB

		err := rows.Scan(
			&chat.Id,
			&chat.Type,
			&chat.CreatedBy,
			&chat.Title,
			&state.LastReadMsgId,
			&state.UnreadCount,
			&state.LastMsgId,
			&state.LastMsgText,
			&state.LastMsgAt,
		)
		if err != nil {
			return nil, domain.ErrDatabase
		}

		chats = append(chats, domain.Chat{
			Id:            chat.Id,
			Type:          domain.ChatType(chat.Type),
			Title:         chat.Title,
			OwnerId:       chat.CreatedBy,
			LastReadMsgId: pointers.ZeroIfNil(state.LastReadMsgId),
			UnreadCount:   pointers.ZeroIfNil(state.UnreadCount),
			LastMsgId:     pointers.ZeroIfNil(state.LastMsgId),
			LastMsgText:   pointers.ZeroIfNil(state.LastMsgText),
			LastMsgTime:   pointers.ZeroIfNil(state.LastMsgAt),
		})
	}

	return chats, nil
}

func (r *Repository) FindPrivateChat(ctx context.Context, userA int64, userB int64) (bool, error) {
	if userA == userB {
		return false, nil
	}

	query, args, err := squirrel.
		Select("COUNT(*)").
		From("chats c").
		Join("chat_participants p1 ON c.id = p1.chat_id").
		Join("chat_participants p2 ON c.id = p2.chat_id").
		Where(squirrel.And{
			squirrel.Eq{"c.type": domain.ChatTypePrivate},
			squirrel.Eq{"p1.user_id": userA},
			squirrel.Eq{"p2.user_id": userB},
			squirrel.Expr("p1.user_id != p2.user_id"),
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return false, domain.ErrDatabase
	}

	var count int
	err = r.q.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, domain.ErrDatabase
	}
	return count > 0, nil
}

func (r *Repository) CreateState(ctx context.Context, chatId int64, userId int64) error {
	query, args, err := squirrel.
		Insert("chat_state").
		Columns("chat_id", "user_id").
		Values(chatId, userId).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return domain.ErrDatabase
	}

	_, err = r.q.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == duplicateKeyCode {
			return domain.ErrStateAlreadyExists
		}
		return domain.ErrDatabase
	}
	return nil
}
