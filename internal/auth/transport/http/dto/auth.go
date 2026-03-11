package dto

import "D/Go/messenger/internal/auth/domain"

type LoginRequest struct {
	Field    domain.UserField `json:"field"`
	Value    string           `json:"value"`
	Password string           `json:"password"`
}

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
