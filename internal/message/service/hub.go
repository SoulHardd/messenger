package service

type Hub interface {
	SendToUser(userId int64, event any)
	SendToUsers(userIds []int64, event any)
}
