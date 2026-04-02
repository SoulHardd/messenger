package service

import (
	"D/Go/messenger/internal/message/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *domain.Message) (int64, error)
	GetChatMessages(ctx context.Context, chatId int64, userId int64, limit int, cursor *domain.Cursor) ([]domain.Message, error)
	MessageExists(ctx context.Context, chatId int64, messageId int64) (bool, error)
	WithTx(tx pgx.Tx) MessageRepository
}

type ChatStateRepository interface {
	UpdateLastMessage(ctx context.Context, chatId int64, msg domain.Message) error
	IncrementUnread(ctx context.Context, chatId int64, excludeUserId int64) error
	ResetUnread(ctx context.Context, chatId int64, userId int64) error
	WithTx(tx pgx.Tx) ChatStateRepository
}

type ParticipantRepository interface {
	GetParticipants(ctx context.Context, chatId int64) ([]int64, error)
	IsParticipant(ctx context.Context, userId int64, chatId int64) (bool, error)
	WithTx(tx pgx.Tx) ParticipantRepository
}

type MessageStatusRepository interface {
	InitForMessage(ctx context.Context, messageId int64, senderId int64, participantIds []int64) error
	MarkDelivered(ctx context.Context, chatId int64, userId int64) error
	MarkReadUpTo(ctx context.Context, chatId int64, userId int64, messageId int64) error
	WithTx(tx pgx.Tx) MessageStatusRepository
}
