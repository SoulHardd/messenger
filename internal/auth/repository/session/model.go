package session

import "time"

type SessionDB struct {
	Id               int64
	UserId           int64
	RefreshTokenHash string
	ExpiresAt        time.Time
}
