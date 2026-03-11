package httpserver

import (
	http2 "D/Go/messenger/internal/auth/transport/http"
	"D/Go/messenger/internal/platform/config"
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
}

func New(
	cfg *config.ServerConfig,
	pool *pgxpool.Pool,
	auth *http2.AuthController,
) *Server {
	r := chi.NewRouter()

	s := &Server{
		HttpSrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: r,
		},
		r:   r,
		cfg: cfg,
		db:  pool,
	}
	s.setAuthRoutes(auth)
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

func (s *Server) setAuthRoutes(auth *http2.AuthController) {
	s.r.Head("/healthcheck", HealthCheck)
	s.r.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.Login)
		r.Post("/register", auth.Register)
		r.Post("/refresh", auth.RefreshTokens)
	})
}
