package service

import (
	"D/Go/messenger/internal/auth/domain"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    UserRepository
	SessionRepo SessionRepository
	JWT         JWTService
}

func New(u UserRepository, s SessionRepository, jwt JWTService) *AuthService {
	return &AuthService{
		UserRepo:    u,
		SessionRepo: s,
		JWT:         jwt,
	}
}

func (s *AuthService) Login(ctx context.Context, ur *domain.UserRaw, field domain.UserField) (*domain.Tokens, error) {
	var user *domain.User
	var err error

	if ur.Password == "" || field == "" {
		return nil, domain.ErrMissingRequiredFields
	}
	if !validatePassword(ur.Password) {
		return nil, domain.ErrInvalidPassword
	}

	switch {
	case field == domain.UserFieldLogin:
		if ur.Login == "" {
			return nil, domain.ErrMissingRequiredFields
		}
		if !validateLogin(ur.Login) {
			return nil, domain.ErrInvalidLogin
		}
		user, err = s.UserRepo.GetOneByLogin(ctx, ur.Login)

	case field == domain.UserFieldPhone:
		if ur.Phone == "" {
			return nil, domain.ErrMissingRequiredFields
		}
		if !validatePhoneNumber(ur.Phone) {
			return nil, domain.ErrInvalidPhone
		}
		user, err = s.UserRepo.GetOneByPhone(ctx, ur.Phone)
	default:
		return nil, domain.ErrDatabase
	}

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(ur.Password),
	)
	if err != nil {
		return nil, domain.ErrIncorrectPassword
	}

	return s.generateTokens(ctx, user.Id)
}

func (s *AuthService) Register(ctx context.Context, ur *domain.UserRaw) (*domain.Tokens, error) {
	if ur.Phone == "" || ur.Login == "" || ur.Password == "" {
		return nil, domain.ErrMissingRequiredFields
	}
	if !validatePhoneNumber(ur.Phone) {
		return nil, domain.ErrInvalidPhone
	}
	if !validateLogin(ur.Login) {
		return nil, domain.ErrInvalidLogin
	}
	if !validatePassword(ur.Password) {
		return nil, domain.ErrInvalidPassword
	}

	pwHash, err := bcrypt.GenerateFromPassword(
		[]byte(ur.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Phone:        ur.Phone,
		Login:        ur.Login,
		PasswordHash: string(pwHash),
	}

	id, err := s.UserRepo.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, id)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*domain.Tokens, error) {

	hash := sha256.Sum256([]byte(refreshToken))
	hashStr := hex.EncodeToString(hash[:])

	session, err := s.SessionRepo.GetByRefreshTokenHash(ctx, hashStr)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	err = s.SessionRepo.DeleteByRefreshTokenHash(ctx, hashStr)
	if err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, session.UserId)
}

func (s *AuthService) InterruptAllSessions(ctx context.Context, userId int64) error {
	err := s.SessionRepo.DeleteByUserID(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) StartMonitorSessions(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.SessionRepo.DeleteExpired(ctx); err != nil {
					log.Printf("delete expired sessions: %v", err)
				}
			}
		}
	}()
}

func (s *AuthService) generateTokens(ctx context.Context, userId int64) (*domain.Tokens, error) {
	accessToken, err := s.JWT.GenerateAccess(userId)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.JWT.GenerateRefresh()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(refreshToken))

	session := &domain.Session{
		UserId:           userId,
		RefreshTokenHash: hex.EncodeToString(hash[:]),
	}

	err = s.SessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

var phoneRegex = regexp.MustCompile(`^\+\d{11}$`)
var passwordRegex = regexp.MustCompile(`^[A-Za-z\d@$!%*#?&]{8,}$`)
var loginRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3,30}$`)

func validatePhoneNumber(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func validatePassword(password string) bool {
	return passwordRegex.MatchString(password)
}

func validateLogin(login string) bool {
	return loginRegex.MatchString(login)
}
