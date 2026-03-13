package dto

import "D/Go/messenger/internal/user/domain"

func ToProfileMeResponse(p *domain.Profile) ProfileMeResponse {
	return ProfileMeResponse{
		Login:     p.Login,
		Phone:     p.Phone,
		Nickname:  p.Nickname,
		Bio:       p.Bio,
		AvatarURL: p.AvatarURL,
	}
}

func ToProfileResponse(p *domain.Profile) ProfileResponse {
	return ProfileResponse{
		Login:     p.Login,
		Nickname:  p.Nickname,
		Bio:       p.Bio,
		AvatarURL: p.AvatarURL,
	}
}

func ToDomainUpdateProfile(req *UpdateProfileRequest) domain.UpdateProfile {
	return domain.UpdateProfile{
		Nickname:  req.Nickname,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
	}
}

func ToSearchResponse(profiles []domain.Profile) SearchResponse {
	resp := SearchResponse{
		Users: make([]ProfileResponse, 0, len(profiles)),
	}
	for _, p := range profiles {
		resp.Users = append(resp.Users, ProfileResponse{
			Login:     p.Login,
			Nickname:  p.Nickname,
			Bio:       p.Bio,
			AvatarURL: p.AvatarURL,
		})
	}

	return resp
}
