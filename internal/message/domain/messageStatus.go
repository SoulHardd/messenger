package domain

type MessageStatus int

var (
	StatusSent      MessageStatus = 1
	StatusDelivered MessageStatus = 2
	StatusRead      MessageStatus = 3
)
