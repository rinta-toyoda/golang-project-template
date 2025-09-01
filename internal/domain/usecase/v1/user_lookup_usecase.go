package user

import (
	"context"

	"example.com/internal/domain/entity"
	userservice "example.com/internal/domain/service/v1"
)

type UserLookupUseCase interface {
	Call(ctx context.Context, email string) (*entity.User, error)
}

type userLookupUseCase struct {
	userService userservice.Service
}

func NewUserLookupUseCase(userService userservice.Service) UserLookupUseCase {
	return &userLookupUseCase{
		userService: userService,
	}
}

func (uc *userLookupUseCase) Call(ctx context.Context, email string) (*entity.User, error) {
	// Find user by email
	user, err := uc.userService.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
