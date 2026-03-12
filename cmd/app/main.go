package main

import (
	JWTService "D/Go/messenger/internal/auth/jwt"
	authSessionRepository "D/Go/messenger/internal/auth/repository/session"
	authUserRepository "D/Go/messenger/internal/auth/repository/user"
	authSrvc "D/Go/messenger/internal/auth/service"
	authController "D/Go/messenger/internal/auth/transport/http"
	"D/Go/messenger/internal/platform/config"
	"D/Go/messenger/internal/platform/database/connections"
	"D/Go/messenger/internal/platform/httpserver"
	"D/Go/messenger/internal/platform/jwt"
	authMW "D/Go/messenger/internal/platform/middleware/auth"
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	pool := connections.InitPool(appCtx, cfg.DatabaseCfg)

	authUserRepo := authUserRepository.New(pool)
	authSessionRepo := authSessionRepository.New(pool, cfg.JWTCfg.RefreshTokenTTL)
	jwtService := JWTService.New(cfg.JWTCfg.Secret, cfg.JWTCfg.AccessTokenTTL)
	authService := authSrvc.New(authUserRepo, authSessionRepo, jwtService)
	authCtrl := authController.New(authService)

	jwtVerifier := jwt.NewVerifier(cfg.JWTCfg.Secret)
	authMiddleware := authMW.Auth(jwtVerifier)

	srv := httpserver.New(cfg.ServerCfg, pool, authCtrl, authMiddleware)

	authService.StartMonitorSessions(appCtx, cfg.TimeCfg.TokensMonitorInterval)

	srv.ListenAndServe()

	<-shutdownCtx.Done()
	log.Printf("Shutting down service-courier")
	gsCtx, gsCancel := context.WithTimeout(context.Background(), cfg.TimeCfg.ShutdownTimeout)
	defer gsCancel()

	if err := srv.HttpSrv.Shutdown(gsCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	pool.Close()
	log.Println("Graceful shutdown complete")
}
