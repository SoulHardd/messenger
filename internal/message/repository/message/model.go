package message

import "time"

type MessageDB struct {
	Id        int64
	ChatId    int64
	SenderId  int64
	Content   string
	Status    int16
	IsEdited  bool
	CreatedAt time.Time
}
