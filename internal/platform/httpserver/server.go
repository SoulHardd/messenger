package httpserver

import (
	httpAuth "D/Go/messenger/internal/auth/transport/http"
	"D/Go/messenger/internal/platform/config"
	authMW "D/Go/messenger/internal/platform/middleware/auth"
	httpUser "D/Go/messenger/internal/user/transport/http"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	HttpSrv *http.Server
	r       *chi.Mux
	cfg     *config.ServerConfig
	db      *pgxpool.Pool
	authMW  authMW.AuthMiddleware
}

func New(
	cfg *config.ServerConfig,
	pool *pgxpool.Pool,
	auth *httpAuth.AuthController,
	authMiddleware authMW.AuthMiddleware,
	user *httpUser.UserController,
) *Server {
	r := chi.NewRouter()

	s := &Server{
		HttpSrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: r,
		},
		r:      r,
		cfg:    cfg,
		db:     pool,
		authMW: authMiddleware,
	}

	s.setAuthRoutes(auth)
	s.setUserRoutes(user, authMiddleware)
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
