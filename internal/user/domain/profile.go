package domain

type Profile struct {
	UserId    int64
	Nickname  string
	Bio       string
	AvatarURL string
	Login     string
	Phone     string
}

type UpdateProfile struct {
	Bio       *string
	Nickname  *string
	AvatarURL *string
}
