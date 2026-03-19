package service

import (
	"D/Go/messenger/internal/chat/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *domain.Chat) (int64, error)
	GetUserChat(ctx context.Context, userId int64, chatId int64) (*domain.Chat, error)
	GetUserChats(ctx context.Context, userId int64, limit int, cursor *domain.Cursor) ([]domain.Chat, error)
	FindPrivateChat(ctx context.Context, userA int64, userB int64) (bool, error)
	CreateState(ctx context.Context, chatId int64, userId int64) error
	WithTx(tx pgx.Tx) ChatRepository
}

type ParticipantRepository interface {
	Add(ctx context.Context, p domain.Participant) error
	Remove(ctx context.Context, p domain.Participant) error
	ChangeParticipantRole(ctx context.Context, p domain.Participant) error
	GetParticipants(ctx context.Context, chatId int64) ([]domain.Participant, error)
	GetParticipant(ctx context.Context, chatId int64, userId int64) (*domain.Participant, error)
	WithTx(tx pgx.Tx) ParticipantRepository
}
