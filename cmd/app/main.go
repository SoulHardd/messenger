package main

import (
	"D/Go/messenger/internal/platform/config"
	"D/Go/messenger/internal/platform/database/connections"
	"D/Go/messenger/internal/platform/httpserver"
	"D/Go/messenger/internal/platform/jwt"
	"context"
	"log"
	"os/signal"
	"syscall"

	JWTService "D/Go/messenger/internal/auth/jwt"
	authSessionRepository "D/Go/messenger/internal/auth/repository/session"
	authUserRepository "D/Go/messenger/internal/auth/repository/user"
	authSrvc "D/Go/messenger/internal/auth/service"
	authCtrl "D/Go/messenger/internal/auth/transport/http"

	authMW "D/Go/messenger/internal/platform/middleware/auth"

	profileRepository "D/Go/messenger/internal/user/repository/profile"
	userRepository "D/Go/messenger/internal/user/repository/user"
	userSrvc "D/Go/messenger/internal/user/service"
	userCtrl "D/Go/messenger/internal/user/transport/http"

	chatRepository "D/Go/messenger/internal/chat/repository/chat"
	partRepository "D/Go/messenger/internal/chat/repository/participant"
	chatSrvc "D/Go/messenger/internal/chat/service"
	chatCtrl "D/Go/messenger/internal/chat/transport/http"
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
	authController := authCtrl.New(authService)

	userRepo := userRepository.New(pool)
	profileRepo := profileRepository.New(pool)
	userService := userSrvc.New(userRepo, profileRepo)
	userController := userCtrl.New(userService)

	partRepo := partRepository.New(pool)
	chatRepo := chatRepository.New(pool)
	chatService := chatSrvc.New(chatRepo, partRepo, pool)
	chatController := chatCtrl.New(chatService)

	jwtVerifier := jwt.NewVerifier(cfg.JWTCfg.Secret)
	authMiddleware := authMW.Auth(jwtVerifier)

	srv := httpserver.New(cfg.ServerCfg, authController, authMiddleware, userController, chatController)

	authService.StartMonitorSessions(appCtx, cfg.TimeCfg.TokensMonitorInterval)

	srv.ListenAndServe()

	<-shutdownCtx.Done()
	log.Printf("Shutting down...")
	gsCtx, gsCancel := context.WithTimeout(context.Background(), cfg.TimeCfg.ShutdownTimeout)
	defer gsCancel()

	if err := srv.HttpSrv.Shutdown(gsCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	pool.Close()
	log.Println("Graceful shutdown complete")
}
