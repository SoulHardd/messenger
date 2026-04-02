package dto

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type NewMessageEvent struct {
	Id        int64  `json:"id"`
	ChatId    int64  `json:"chat_id"`
	SenderId  int64  `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type ReadEvent struct {
	ChatId    int64 `json:"chat_id"`
	UserId    int64 `json:"user_id"`
	MessageId int64 `json:"message_id"`
}
