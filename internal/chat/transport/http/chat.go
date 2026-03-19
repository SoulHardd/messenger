package http

import (
	"D/Go/messenger/internal/chat/domain"
	"D/Go/messenger/internal/chat/transport/http/dto"
	"D/Go/messenger/internal/platform/httpx"
	authMiddleware "D/Go/messenger/internal/platform/middleware/auth"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ChatController struct {
	service ChatService
}

func New(s ChatService) *ChatController {
	return &ChatController{service: s}
}

func (c *ChatController) CreatePrivateChat(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	var req dto.CreatePrivateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	chat := dto.ToDomainPrivateChat(req, userId)

	id, err := c.service.CreatePrivateChat(r.Context(), chat)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusCreated, dto.IdResponse{Id: id})
}

func (c *ChatController) CreateGroupChat(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	var req dto.CreateGroupChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	chat := dto.ToDomainGroupChat(req, userId)

	id, err := c.service.CreateGroupChat(r.Context(), chat)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusCreated, dto.IdResponse{Id: id})
}

func (c *ChatController) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
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

	chats, nextCursor, err := c.service.GetUserChats(r.Context(), userId, limit, cursor)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	response := dto.ToChatListResponse(chats, nextCursor)

	httpx.WriteResponse(w, http.StatusOK, response)
}

func (c *ChatController) GetUserChatById(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	chatIdStr := chi.URLParam(r, "chat_id")
	if chatIdStr == "" {
		httpx.WriteErrorResponse(w, errors.New("chat id is required"), ErrorMapper{})
		return
	}
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		httpx.WriteErrorResponse(w, errors.New("invalid chat id format"), ErrorMapper{})
		return
	}

	chat, err := c.service.GetUserChatById(r.Context(), userId, chatId)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusOK, dto.ToChatResponse(*chat))
}

func (c *ChatController) GetChatParticipants(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	idStr := r.URL.Query().Get("chat_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidQuery, ErrorMapper{})
		return
	}

	p, err := c.service.GetChatParticipants(r.Context(), userId, id)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusOK, dto.ToParticipantsResponse(p))
}

func (c *ChatController) AddParticipant(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	var req dto.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	domainPart, err := dto.ToDomainParticipant(req)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}
	err = c.service.AddParticipant(r.Context(), userId, *domainPart)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusOK, nil)
}

func (c *ChatController) DeleteParticipant(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	var req dto.RemovePartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	err := c.service.RemoveParticipant(r.Context(), userId, dto.ToDomainRemoveParticipant(req))
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusOK, nil)
}

func (c *ChatController) UpdateParticipantRole(w http.ResponseWriter, r *http.Request) {
	userId, ok := authMiddleware.UserId(r.Context())
	if !ok {
		httpx.WriteErrorResponse(w, httpx.ErrUnauthorized, ErrorMapper{})
		return
	}

	var req dto.ParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErrorResponse(w, httpx.ErrInvalidJSON, ErrorMapper{})
		return
	}

	domainPart, err := dto.ToDomainParticipant(req)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}
	err = c.service.ChangeParticipantRole(r.Context(), userId, *domainPart)
	if err != nil {
		httpx.WriteErrorResponse(w, err, ErrorMapper{})
		return
	}

	httpx.WriteResponse(w, http.StatusOK, nil)
}
