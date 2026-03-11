package user

import "time"

type UserDB struct {
	Id           int64
	Phone        string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
