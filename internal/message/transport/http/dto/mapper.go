package dto

import (
	"D/Go/messenger/internal/message/domain"
	"time"
)

const dateLayout = time.RFC3339Nano

func ToDomainMessage(r SendMessageRequest, senderId int64) domain.Message {
	return domain.Message{
		ChatId:   r.ChatId,
		SenderId: senderId,
		Content:  r.Content,
	}
}

func ToMessageResponse(m domain.Message) MessageResponse {
	return MessageResponse{
		Id:        m.Id,
		ChatId:    m.ChatId,
		SenderId:  m.SenderId,
		Content:   m.Content,
		Status:    int(m.Status),
		CreatedAt: m.CreatedAt.Format(dateLayout),
	}
}

func ToMessageListResponse(messages []domain.Message, cursor *domain.Cursor) MessageListResponse {
	resp := MessageListResponse{
		Messages: make([]MessageResponse, 0, len(messages)),
	}

	for _, m := range messages {
		resp.Messages = append(resp.Messages, ToMessageResponse(m))
	}

	if cursor != nil {
		encoded := EncodeCursor(cursor)
		resp.NextCursor = &encoded
	}

	return resp
}
