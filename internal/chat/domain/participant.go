package domain

type Participant struct {
	ChatId int64
	UserId int64
	Role   ParticipantRole
}
