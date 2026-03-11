package dto

import (
	"D/Go/messenger/internal/auth/domain"
	"fmt"
)

func ToDomainLogin(req *LoginRequest) (domain.UserRaw, error) {

	switch {
	case req.Field == domain.UserFieldLogin:
		return domain.UserRaw{
			Login:    req.Value,
			Password: req.Password,
		}, nil
	case req.Field == domain.UserFieldPhone:
		return domain.UserRaw{
			Phone:    req.Value,
			Password: req.Password,
		}, nil

	default:
		return domain.UserRaw{}, fmt.Errorf("invalid field")
	}
}

func ToDomainRegister(req *RegisterRequest) domain.UserRaw {
	return domain.UserRaw{
		Phone:    req.Phone,
		Login:    req.Login,
		Password: req.Password,
	}
}

func ToDomainToken(req *RefreshTokenRequest) domain.Tokens {
	return domain.Tokens{
		RefreshToken: req.RefreshToken,
	}
}

func ToTokenResponse(tokens *domain.Tokens) TokenResponse {
	return TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
