package auth

import (
	"context"

	"example.com/internal/domain/entity"
	authservice "example.com/internal/domain/service/auth"
)

type LoginUseCase interface {
	Call(ctx context.Context, email, password string) (*entity.User, error)
}

type loginUseCase struct {
	authService authservice.Service
}

func NewLoginUseCase(authService authservice.Service) LoginUseCase {
	return &loginUseCase{
		authService: authService,
	}
}

func (uc *loginUseCase) Call(ctx context.Context, email, password string) (*entity.User, error) {
	// Authenticate user
	user, err := uc.authService.AuthenticateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	// Update last login time (non-critical operation)
	if err := uc.authService.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail login - this is non-critical
		_ = err
	}

	return user, nil
}
