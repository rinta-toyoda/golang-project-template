package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	authapi "example.com/gen/openapi/auth/go"
	"example.com/internal/domain/entity"
	"example.com/internal/domain/repository"
	"example.com/pkg/security"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service interface {
	SignUp(ctx context.Context, req authapi.SignupRequest) (*authapi.SignupResponse, error)
	Login(ctx context.Context, req authapi.LoginRequest) (*authapi.LoginResponse, error)
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

func (s *service) SignUp(ctx context.Context, req authapi.SignupRequest) (*authapi.SignupResponse, error) {
	if req.Username != "" {
		if _, err := s.userRepo.FindByUserName(ctx, req.Username); err == nil {
			return nil, ErrUserAlreadyExists
		}
	}

	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:           uuid.NewString(),
		UserName:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	apiUser := authapi.User{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.UserName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.LastLoginAt != nil {
		apiUser.LastLoginAt = *user.LastLoginAt
	}

	return &authapi.SignupResponse{
		User:    apiUser,
		Message: "User created successfully",
	}, nil
}

func (s *service) Login(ctx context.Context, req authapi.LoginRequest) (*authapi.LoginResponse, error) {
	user, err := s.userRepo.FindByUserNameOrEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.hasher.Verify(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		// Log error but don't fail login - this is non-critical
		_ = err
	}

	apiUser := authapi.User{
		Id:          user.ID,
		Email:       user.Email,
		Username:    user.UserName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		LastLoginAt: now,
	}

	return &authapi.LoginResponse{
		User:    apiUser,
		Message: "Login successful",
	}, nil
}
