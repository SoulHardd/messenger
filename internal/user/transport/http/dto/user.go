package dto

type ProfileMeResponse struct {
	Login     string `json:"login"`
	Phone     string `json:"phone"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type ProfileResponse struct {
	Login     string `json:"login"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateProfileRequest struct {
	Nickname  *string `json:"nickname"`
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}

type SearchResponse struct {
	Users []ProfileResponse `json:"users"`
}
