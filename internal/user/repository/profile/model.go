package profile

type ProfileDB struct {
	UserId    int64
	Login     *string
	Nickname  *string
	Bio       *string
	AvatarURL *string
}
