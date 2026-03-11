package http

import (
	"D/Go/messenger/internal/auth/domain"
	"D/Go/messenger/internal/auth/transport/http/dto"
	http2 "D/Go/messenger/internal/platform/http"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	service AuthService
}

func New(s AuthService) *AuthController {
	return &AuthController{
		service: s,
	}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.WriteErrorResponse(w, http2.ErrInvalidJSON)
		return
	}

	userRaw, err := dto.ToDomainLogin(&req)
	if err != nil {
		http2.WriteErrorResponse(w, http2.ErrInvalidJSON)
		return
	}

	var tokens *domain.Tokens
	tokens, err = c.service.Login(r.Context(), &userRaw, req.Field)
	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := dto.ToTokenResponse(tokens)

	http2.WriteResponse(w, http.StatusOK, response)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.WriteErrorResponse(w, http2.ErrInvalidJSON)
		return
	}

	userRaw := dto.ToDomainRegister(&req)

	tokens, err := c.service.Register(r.Context(), &userRaw)
	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := dto.ToTokenResponse(tokens)

	http2.WriteResponse(w, http.StatusCreated, response)
}

func (c *AuthController) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http2.WriteErrorResponse(w, http2.ErrInvalidJSON)
		return
	}

	token := dto.ToDomainToken(&req)
	tokens, err := c.service.RefreshTokens(r.Context(), token.RefreshToken)
	if err != nil {
		http2.WriteErrorResponse(w, err)
		return
	}

	response := dto.ToTokenResponse(tokens)

	http2.WriteResponse(w, http.StatusOK, response)
}
