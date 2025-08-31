package service_test

import (
	"context"
	"testing"

	"example.com/internal/domain/entity"
	"example.com/internal/domain/service"
	"example.com/test/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_SignUp(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authService := service.NewAuthService(mockRepo, mockHasher)

	ctx := context.Background()
	req := service.SignUpRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock repository calls
	mockRepo.On("FindByUserName", ctx, "testuser").Return(nil, assert.AnError)
	mockRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, assert.AnError)
	mockHasher.On("Hash", "password123").Return("hashed_password", nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, err := authService.SignUp(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.UserName)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authService := service.NewAuthService(mockRepo, mockHasher)

	ctx := context.Background()
	req := service.LoginRequest{
		Identifier: "test@example.com",
		Password:   "password123",
	}

	expectedUser := &entity.User{
		ID:           "user-id",
		Email:        "test@example.com",
		UserName:     "testuser",
		PasswordHash: "hashed_password",
	}

	mockRepo.On("FindByUserNameOrEmail", ctx, "test@example.com").Return(expectedUser, nil)
	mockHasher.On("Verify", "password123", "hashed_password").Return(true)

	user, err := authService.Login(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
