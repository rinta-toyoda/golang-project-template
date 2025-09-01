package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authapi "example.com/gen/openapi/auth/go"
	"example.com/internal/domain/entity"
	authservice "example.com/internal/domain/service/auth"
	"example.com/test/unit/mocks"
)

func TestAuthService_SignUp(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	req := authapi.SignupRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock repository calls
	mockRepo.On("FindByUserName", ctx, "testuser").Return(nil, assert.AnError)
	mockRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, assert.AnError)
	mockHasher.On("Hash", "password123").Return("hashed_password", nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	response, err := authSvc.SignUp(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "test@example.com", response.User.Email)
	assert.Equal(t, "testuser", response.User.Username)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	req := authapi.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedUser := &entity.User{
		ID:           "user-id",
		Email:        "test@example.com",
		UserName:     "testuser",
		PasswordHash: "hashed_password",
	}

	mockRepo.On("FindByUserNameOrEmail", ctx, "test@example.com").Return(expectedUser, nil)
	mockHasher.On("Verify", "password123", "hashed_password").Return(true)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	response, err := authSvc.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "user-id", response.User.Id)
	assert.Equal(t, "test@example.com", response.User.Email)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}