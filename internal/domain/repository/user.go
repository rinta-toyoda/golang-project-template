package repository

import (
	"context"

	"example.com/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUserName(ctx context.Context, userName string) (*entity.User, error)
	FindByUserNameOrEmail(ctx context.Context, identifier string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}
