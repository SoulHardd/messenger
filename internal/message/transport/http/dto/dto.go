package dto

type SendMessageRequest struct {
	ChatId  int64  `json:"chat_id"`
	Content string `json:"content"`
}

type MarkAsReadRequest struct {
	ChatId    int64 `json:"chat_id"`
	MessageId int64 `json:"message_id"`
}

type IdResponse struct {
	Id int64 `json:"id"`
}

type MessageResponse struct {
	Id        int64  `json:"id"`
	ChatId    int64  `json:"chat_id"`
	SenderId  int64  `json:"sender_id"`
	Content   string `json:"content"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type MessageListResponse struct {
	Messages   []MessageResponse `json:"messages"`
	NextCursor *string           `json:"next_cursor,omitempty"`
}
