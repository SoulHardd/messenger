package http

import (
	"D/Go/messenger/internal/chat/domain"
	"context"
)

type ChatService interface {
	CreatePrivateChat(ctx context.Context, chat domain.PrivateChat) (int64, error)
	CreateGroupChat(ctx context.Context, chat domain.GroupChat) (int64, error)
	GetUserChats(ctx context.Context, userId int64, limit int, cursor *domain.Cursor) ([]domain.Chat, *domain.Cursor, error)
	GetUserChatById(ctx context.Context, userId int64, chatId int64) (*domain.Chat, error)
	GetChatParticipants(ctx context.Context, userId int64, chatId int64) ([]domain.Participant, error)
	AddParticipant(ctx context.Context, userId int64, p domain.Participant) error
	RemoveParticipant(ctx context.Context, userId int64, p domain.Participant) error
	ChangeParticipantRole(ctx context.Context, userId int64, p domain.Participant) error
}
