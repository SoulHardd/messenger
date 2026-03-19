package domain

import "time"

type Chat struct {
	Id            int64
	Type          ChatType
	Title         string
	OwnerId       int64
	LastReadMsgId int64
	UnreadCount   int
	LastMsgId     int64
	LastMsgText   string
	LastMsgTime   time.Time
}

type PrivateChat struct {
	FirstUId  int64
	SecondUId int64
}

type GroupChat struct {
	OwnerId int64
	Title   string
	Users   []int64
}
