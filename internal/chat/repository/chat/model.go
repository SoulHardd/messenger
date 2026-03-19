package chat

import "time"

type ChatDB struct {
	Id        int64
	Type      int
	Title     string
	CreatedBy int64
	CreatedAt time.Time
}

type ChatStateDB struct {
	LastReadMsgId *int64
	UnreadCount   *int
	LastMsgId     *int64
	LastMsgText   *string
	LastMsgAt     *time.Time
}
