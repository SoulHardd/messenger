package http

import (
	"D/Go/messenger/internal/message/domain"
	"D/Go/messenger/internal/message/transport/http/dto"
	"D/Go/messenger/internal/platform/httpx"
	authMiddleware "D/Go/messenger/internal/platform/middleware/auth"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MessageController struct {
	service MessageService
}

func New(s MessageService) *MessageController {
	return &MessageController{service: s}
}

func (c *MessageController) SendMessage(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	msg := dto.ToDomainMessage(req, userId)

	id, err := c.service.SendMessage(r.Context(), msg)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusCreated, dto.IdResponse{Id: id})
}

func (c *MessageController) GetMessages(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	chatIdRaw := chi.URLParam(r, "chat_id")
	chatId, err := strconv.ParseInt(chatIdRaw, 10, 64)
	if err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidQuery, ErrorMapper{})
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		val, err := strconv.Atoi(l)
		if err != nil {
			httpx.WriteErrorResponse(w, httpx.ErrInvalidQuery, ErrorMapper{})
			return
		}
		limit = val
	}

	var cursor *domain.Cursor
	if cStr := r.URL.Query().Get("cursor"); cStr != "" {
		cParsed, err := dto.DecodeCursor(cStr)
		if err != nil {
			httpx.WriteErrorResponse(w, httpx.ErrInvalidQuery, ErrorMapper{})
			return
		}
		cursor = cParsed
	}

	messages, nextCursor, err := c.service.GetMessages(r.Context(), userId, chatId, limit, cursor)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}
	httpx.WriteResponse(w, http.StatusOK, dto.ToMessageListResponse(messages, nextCursor))
}

func (c *MessageController) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	var req dto.MarkAsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	err := c.service.MarkAsRead(r.Context(), userId, req.ChatId, req.MessageId)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusOK, nil)
}
