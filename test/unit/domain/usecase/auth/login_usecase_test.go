package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"example.com/internal/domain/entity"
	authservice "example.com/internal/domain/service/auth"
	authusecase "example.com/internal/domain/usecase/auth"
	"example.com/test/unit/mocks"
)

func TestLoginUseCase_Call_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewLoginUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	existingUser := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		UserName:     "testuser",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock authentication flow
	mockRepo.On("FindByUserNameOrEmail", ctx, email).Return(existingUser, nil)
	mockHasher.On("Verify", password, "hashed_password").Return(true)

	// Mock last login update
	mockRepo.On("FindByID", ctx, "user-123").Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, err := useCase.Call(ctx, email, password)

	assert.NoError(t, err)
	assert.Equal(t, existingUser, user)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLoginUseCase_Call_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewLoginUseCase(authSvc)

	ctx := context.Background()
	email := "notfound@example.com"
	password := "password123"

	// Mock user not found
	mockRepo.On("FindByUserNameOrEmail", ctx, email).Return(nil, errors.New("user not found"))

	user, err := useCase.Call(ctx, email, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, authservice.ErrInvalidCredentials, err)
	mockRepo.AssertExpectations(t)
}

func TestLoginUseCase_Call_InvalidPassword(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewLoginUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "wrongpassword"

	existingUser := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		UserName:     "testuser",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock authentication flow with wrong password
	mockRepo.On("FindByUserNameOrEmail", ctx, email).Return(existingUser, nil)
	mockHasher.On("Verify", password, "hashed_password").Return(false)

	user, err := useCase.Call(ctx, email, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, authservice.ErrInvalidCredentials, err)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLoginUseCase_Call_UpdateLastLoginError(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewLoginUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	existingUser := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		UserName:     "testuser",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Mock authentication flow
	mockRepo.On("FindByUserNameOrEmail", ctx, email).Return(existingUser, nil)
	mockHasher.On("Verify", password, "hashed_password").Return(true)

	// Mock last login update failure (non-critical)
	mockRepo.On("FindByID", ctx, "user-123").Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("update failed"))

	user, err := useCase.Call(ctx, email, password)

	// UpdateLastLogin error should not fail login (non-critical operation)
	assert.NoError(t, err)
	assert.Equal(t, existingUser, user)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
