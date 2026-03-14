package service

import (
	"D/Go/messenger/internal/user/domain"
	"context"
)

type UserService struct {
	UserRepo    UserRepository
	ProfileRepo ProfileRepository
}

func New(ur UserRepository, pr ProfileRepository) *UserService {
	return &UserService{
		UserRepo:    ur,
		ProfileRepo: pr,
	}
}

func (s *UserService) GetMe(ctx context.Context, userId int64) (*domain.Profile, error) {
	var (
		user    *domain.User
		profile *domain.Profile
		err     error
	)
	user, err = s.UserRepo.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	profile, err = s.ProfileRepo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	profile.Login = user.Login
	profile.Phone = user.Phone

	return profile, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userId int64, update *domain.UpdateProfile) error {
	if update.Nickname == nil && update.Bio == nil && update.AvatarURL == nil {
		return domain.ErrEmptyUpdate
	}

	return s.ProfileRepo.Update(ctx, userId, update)
}

func (s *UserService) GetByLogin(ctx context.Context, login string) (*domain.Profile, error) {
	var (
		user    *domain.User
		profile *domain.Profile
		err     error
	)

	if login == "" {
		return nil, domain.ErrInvalidLogin
	}

	user, err = s.UserRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	profile, err = s.ProfileRepo.GetByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	profile.Login = user.Login

	return profile, nil
}

func (s *UserService) GetByPhone(ctx context.Context, phone string) (*domain.Profile, error) {
	var (
		user    *domain.User
		profile *domain.Profile
		err     error
	)
	user, err = s.UserRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	profile, err = s.ProfileRepo.GetByUserId(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	profile.Login = user.Login

	return profile, nil
}

func (s *UserService) Search(ctx context.Context, query string, limit int) ([]domain.Profile, error) {
	if query == "" {
		return nil, domain.ErrInvalidQuery
	}
	if limit <= 0 {
		return nil, domain.ErrInvalidLimit
	}
	if limit > 50 {
		limit = 20
	}

	profiles, err := s.ProfileRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}
