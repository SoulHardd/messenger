package httpserver

import (
	"D/Go/messenger/internal/platform/config"
	authMW "D/Go/messenger/internal/platform/middleware/auth"
	"errors"
	"fmt"
	"log"
	"net/http"

	httpAuth "D/Go/messenger/internal/auth/transport/http"
	httpChat "D/Go/messenger/internal/chat/transport/http"
	httpMsg "D/Go/messenger/internal/message/transport/http"
	httpUser "D/Go/messenger/internal/user/transport/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	HttpSrv *http.Server
	r       *chi.Mux
	cfg     *config.ServerConfig
}

func New(
	cfg *config.ServerConfig,
	auth *httpAuth.AuthController,
	authMiddleware authMW.AuthMiddleware,
	user *httpUser.UserController,
	chat *httpChat.ChatController,
	msg *httpMsg.MessageController,
) *Server {
	r := chi.NewRouter()

	s := &Server{
		HttpSrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: r,
		},
		r:   r,
		cfg: cfg,
	}

	s.setAuthRoutes(auth)
	s.setUserRoutes(user, authMiddleware)
	s.setChatRoutes(chat, authMiddleware)
	s.setMessageRoutes(msg, authMiddleware)
	return s
}

func (s *Server) ListenAndServe() {
	go func() {
		log.Printf("Server is listening on %d", s.cfg.Port)
		if err := s.HttpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server error: %v", err)
		}
	}()
}

func (s *Server) setAuthRoutes(auth *httpAuth.AuthController) {
	s.r.Route("/api/v1", func(r chi.Router) {
		r.Head("/healthcheck", HealthCheck)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", auth.Login)
			r.Post("/register", auth.Register)
			r.Post("/refresh", auth.RefreshTokens)
		})
	})
}

func (s *Server) setUserRoutes(user *httpUser.UserController, authMiddleware authMW.AuthMiddleware) {
	s.r.Route("/api/v1/user", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/me", user.GetMe)
		r.Patch("/me", user.UpdateProfile)
		r.Get("/{login}", user.GetByLogin)
		r.Get("/search", user.Search)
	})
}

func (s *Server) setChatRoutes(chat *httpChat.ChatController, authMiddleware authMW.AuthMiddleware) {
	s.r.Route("/api/v1/chat", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/private", chat.CreatePrivateChat)
		r.Post("/group", chat.CreateGroupChat)
		r.Get("/chats", chat.GetUserChats)
		r.Get("/{chat_id}", chat.GetUserChatById)
		r.Get("/{chat_id}/participants", chat.GetChatParticipants)
		r.Post("/participant", chat.AddParticipant)
		r.Patch("/participant", chat.UpdateParticipantRole)
		r.Delete("/participant", chat.DeleteParticipant)
	})
}

func (s *Server) setMessageRoutes(msg *httpMsg.MessageController, authMiddleware authMW.AuthMiddleware) {
	s.r.Route("/api/v1/message", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Post("/", msg.SendMessage)
		r.Get("/{chat_id}", msg.GetMessages)
		r.Post("/mark-as-read", msg.MarkAsRead)
	})
}
