package user

import (
	"context"
	"errors"

	v1api "example.com/gen/openapi/v1/go"
	"example.com/internal/domain/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Service interface {
	LookupUser(ctx context.Context, email string) (*v1api.UserLookupResponse, error)
}

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) LookupUser(ctx context.Context, email string) (*v1api.UserLookupResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &v1api.UserLookupResponse{
		Username: user.UserName,
		Email:    user.Email,
	}, nil
}
