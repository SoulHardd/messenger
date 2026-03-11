package domain

type User struct {
	Id           int64
	Phone        string
	Login        string
	PasswordHash string
}

type UserRaw struct {
	Id       int64
	Phone    string
	Login    string
	Password string
}
