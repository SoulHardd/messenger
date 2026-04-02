package http

import (
	"D/Go/messenger/internal/message/domain"
	"context"
)

type MessageService interface {
	SendMessage(ctx context.Context, msg domain.Message) (int64, error)
	GetMessages(ctx context.Context, userId int64, chatId int64, limit int, cursor *domain.Cursor) ([]domain.Message, *domain.Cursor, error)
	MarkAsRead(ctx context.Context, userId int64, chatId int64, messageId int64) error
}
