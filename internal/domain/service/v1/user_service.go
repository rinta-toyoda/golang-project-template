package user

import (
	"context"
	"errors"

	"example.com/internal/domain/entity"
	"example.com/internal/domain/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Service interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
