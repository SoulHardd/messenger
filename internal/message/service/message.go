package service

import (
	"D/Go/messenger/internal/message/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageService struct {
	MsgRepo       MessageRepository
	StatusRepo    MessageStatusRepository
	ChatStateRepo ChatStateRepository
	PartRepo      ParticipantRepository
	Hub           Hub
	pool          *pgxpool.Pool
}

func New(m MessageRepository, st MessageStatusRepository, c ChatStateRepository, p ParticipantRepository, h Hub, pool *pgxpool.Pool) *MessageService {
	return &MessageService{
		MsgRepo:       m,
		StatusRepo:    st,
		ChatStateRepo: c,
		PartRepo:      p,
		Hub:           h,
		pool:          pool,
	}
}

func (s *MessageService) SendMessage(ctx context.Context, msg domain.Message) (int64, error) {
	if msg.ChatId == 0 || msg.SenderId == 0 || msg.Content == "" {
		return 0, domain.ErrInvalidInput
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, domain.ErrDatabase
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ok, err := s.PartRepo.IsParticipant(ctx, msg.SenderId, msg.ChatId)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, domain.ErrPermissionDenied
	}

	msgRepoTx := s.MsgRepo.WithTx(tx)
	statusRepoTx := s.StatusRepo.WithTx(tx)
	chatStateRepoTx := s.ChatStateRepo.WithTx(tx)
	partRepoTx := s.PartRepo.WithTx(tx)

	id, err := msgRepoTx.Create(ctx, &msg)
	if err != nil {
		return 0, err
	}

	err = chatStateRepoTx.UpdateLastMessage(ctx, msg.ChatId, msg)
	if err != nil {
		return 0, err
	}

	err = chatStateRepoTx.IncrementUnread(ctx, msg.ChatId, msg.SenderId)
	if err != nil {
		return 0, err
	}

	participants, err := partRepoTx.GetParticipants(ctx, msg.ChatId)
	if err != nil {
		return 0, err
	}

	err = statusRepoTx.InitForMessage(ctx, msg.Id, msg.SenderId, participants)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, domain.ErrDatabase
	}

	if s.Hub != nil {
		s.Hub.SendToUsers(participants, msg)
	}

	return id, nil
}

func (s *MessageService) GetMessages(ctx context.Context, userId int64, chatId int64, limit int, cursor *domain.Cursor) ([]domain.Message, *domain.Cursor, error) {
	if userId == 0 || chatId == 0 || limit <= 0 || limit > 100 {
		return nil, nil, domain.ErrInvalidInput
	}

	ok, err := s.PartRepo.IsParticipant(ctx, userId, chatId)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, domain.ErrPermissionDenied
	}

	err = s.StatusRepo.MarkDelivered(ctx, chatId, userId)
	if err != nil {
		return nil, nil, err
	}

	messages, err := s.MsgRepo.GetChatMessages(ctx, chatId, userId, limit, cursor)
	if err != nil {
		return nil, nil, err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	var nextCursor *domain.Cursor
	if hasMore && len(messages) > 0 {
		last := messages[len(messages)-1]
		nextCursor = &domain.Cursor{
			Time: last.CreatedAt,
			Id:   last.Id,
		}
	}

	return messages, nextCursor, nil
}

func (s *MessageService) MarkAsRead(ctx context.Context, userId int64, chatId int64, messageId int64) error {
	if userId == 0 || chatId == 0 || messageId == 0 {
		return domain.ErrInvalidInput
	}

	ok, err := s.PartRepo.IsParticipant(ctx, userId, chatId)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrPermissionDenied
	}

	ok, err = s.MsgRepo.MessageExists(ctx, chatId, messageId)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrMessageNotFound
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.ErrDatabase
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	statusRepoTx := s.StatusRepo.WithTx(tx)
	chatStateRepoTx := s.ChatStateRepo.WithTx(tx)

	err = statusRepoTx.MarkReadUpTo(ctx, chatId, userId, messageId)
	if err != nil {
		return err
	}

	err = chatStateRepoTx.ResetUnread(ctx, chatId, userId)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.ErrDatabase
	}

	return nil
}
