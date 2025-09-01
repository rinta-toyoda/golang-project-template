package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"example.com/internal/domain/entity"
	"example.com/internal/domain/repository"
	"example.com/pkg/security"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	CheckUserExists(ctx context.Context, email, username string) error
	CreateUser(ctx context.Context, email, password, username string) (*entity.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*entity.User, error)
	UpdateLastLogin(ctx context.Context, userID string) error
}

type service struct {
	userRepo repository.UserRepository
	hasher   security.PasswordHasher
}

func NewService(userRepo repository.UserRepository, hasher security.PasswordHasher) Service {
	return &service{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (s *service) CheckUserExists(ctx context.Context, email, username string) error {
	if username != "" {
		if _, err := s.userRepo.FindByUserName(ctx, username); err == nil {
			return ErrUserAlreadyExists
		}
	}

	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return ErrUserAlreadyExists
	}

	return nil
}

func (s *service) CreateUser(ctx context.Context, email, password, username string) (*entity.User, error) {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:           uuid.NewString(),
		UserName:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) AuthenticateUser(ctx context.Context, email, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByUserNameOrEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.hasher.Verify(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *service) UpdateLastLogin(ctx context.Context, userID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.LastLoginAt = &now
	return s.userRepo.Update(ctx, user)
}
