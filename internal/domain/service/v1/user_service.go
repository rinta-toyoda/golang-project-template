package user

import (
	"context"
	"errors"

	authapi "example.com/gen/openapi/auth/go"
	"example.com/internal/domain/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Service interface {
	LookupUser(ctx context.Context, email string) (*authapi.UserLookupResponse, error)
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) LookupUser(ctx context.Context, email string) (*authapi.UserLookupResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &authapi.UserLookupResponse{
		Username: user.UserName,
		Email:    user.Email,
	}, nil
}