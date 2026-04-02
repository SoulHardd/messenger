package domain

import "time"

type Message struct {
	Id        int64
	ChatId    int64
	SenderId  int64
	Content   string
	Status    MessageStatus
	IsEdited  bool
	CreatedAt time.Time
}
