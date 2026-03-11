package domain

type Session struct {
	Id               int64
	UserId           int64
	RefreshTokenHash string
}
