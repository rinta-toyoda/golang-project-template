package auth

import (
	"context"

	"example.com/internal/domain/entity"
	authservice "example.com/internal/domain/service/auth"
)

type SignupUseCase interface {
	Call(ctx context.Context, email, password, username string) (*entity.User, error)
}

type signupUseCase struct {
	authService authservice.Service
}

func NewSignupUseCase(authService authservice.Service) SignupUseCase {
	return &signupUseCase{
		authService: authService,
	}
}

func (uc *signupUseCase) Call(ctx context.Context, email, password, username string) (*entity.User, error) {
	// Check if user already exists
	if err := uc.authService.CheckUserExists(ctx, email, username); err != nil {
		return nil, err
	}

	// Create new user
	user, err := uc.authService.CreateUser(ctx, email, password, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
