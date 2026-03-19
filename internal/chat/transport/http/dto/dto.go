package dto

import "time"

type CreatePrivateChatRequest struct {
	ParticipantId int64 `json:"user_id"`
}

type CreateGroupChatRequest struct {
	Title string  `json:"title"`
	Users []int64 `json:"users"`
}

type ParticipantRequest struct {
	ChatId int64  `json:"chat_id"`
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
}

type RemovePartRequest struct {
	ChatId int64 `json:"chat_id"`
	UserId int64 `json:"user_id"`
}

type IdResponse struct {
	Id int64 `json:"id"`
}

type ChatResponse struct {
	Id          int64     `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	OwnerId     int64     `json:"owner_id"`
	LastMsgText string    `json:"last_msg_text"`
	LastMsgTime time.Time `json:"last_msg_time"`
	UnreadCount int       `json:"unread_count"`
}

type ChatListResponse struct {
	Chats      []ChatResponse `json:"chats"`
	NextCursor *string        `json:"next_cursor,omitempty"`
}

type ParticipantResponse struct {
	ChatId int64  `json:"chat_id"`
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
}
