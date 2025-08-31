package service

import (
	"context"
	"errors"

	"example.com/internal/domain/entity"
	"example.com/internal/domain/repository"
	"example.com/pkg/security"
	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService interface {
	SignUp(ctx context.Context, req SignUpRequest) (*entity.User, error)
	Login(ctx context.Context, req LoginRequest) (*entity.User, error)
}

type SignUpRequest struct {
	Email    string
	Username string
	Password string
}

type LoginRequest struct {
	Identifier string // email or username
	Password   string
}

type authService struct {
	userRepo repository.UserRepository
	hasher   security.PasswordHasher
}

func NewAuthService(userRepo repository.UserRepository, hasher security.PasswordHasher) AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (s *authService) SignUp(ctx context.Context, req SignUpRequest) (*entity.User, error) {
	if _, err := s.userRepo.FindByUserName(ctx, req.Username); err == nil {
		return nil, ErrUserAlreadyExists
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

	return user, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*entity.User, error) {
	user, err := s.userRepo.FindByUserNameOrEmail(ctx, req.Identifier)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !s.hasher.Verify(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
