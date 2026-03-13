package http

import (
	"D/Go/messenger/internal/platform/httpx"
	authMiddleware "D/Go/messenger/internal/platform/middleware/auth"
	"D/Go/messenger/internal/user/transport/http/dto"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	service UserService
}

func New(service UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (c *UserController) GetMe(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := c.service.GetMe(r.Context(), userId)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	resp := dto.ToProfileMeResponse(profile)

	httpx.WriteResponse(w, http.StatusOK, resp)
}

func (c *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	update := dto.ToDomainUpdateProfile(&req)

	err := c.service.UpdateProfile(r.Context(), userId, &update)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusOK, nil)
}

func (c *UserController) GetByLogin(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")

	profile, err := c.service.GetByLogin(r.Context(), login)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	resp := dto.ToProfileResponse(profile)

	httpx.WriteResponse(w, http.StatusOK, resp)
}

func (c *UserController) GetByPhone(w http.ResponseWriter, r *http.Request) {
	phone := chi.URLParam(r, "phone")

	profile, err := c.service.GetByPhone(r.Context(), phone)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	resp := dto.ToProfileResponse(profile)

	httpx.WriteResponse(w, http.StatusOK, resp)
}

func (c *UserController) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	profiles, err := c.service.Search(r.Context(), query, limit)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	resp := dto.ToSearchResponse(profiles)

	httpx.WriteResponse(w, http.StatusOK, resp)
}
