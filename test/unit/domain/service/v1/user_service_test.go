package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"example.com/internal/domain/entity"
	userservice "example.com/internal/domain/service/v1"
	"example.com/test/unit/mocks"
)

func TestUserService_LookupUser_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	userSvc := userservice.NewService(mockRepo)

	ctx := context.Background()
	email := "test@example.com"

	expectedUser := &entity.User{
		ID:        "user-123",
		Email:     "test@example.com",
		UserName:  "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", ctx, email).Return(expectedUser, nil)

	response, err := userSvc.LookupUser(ctx, email)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "test@example.com", response.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_LookupUser_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	userSvc := userservice.NewService(mockRepo)

	ctx := context.Background()
	email := "notfound@example.com"

	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found"))

	response, err := userSvc.LookupUser(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserService_LookupUser_DatabaseError(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	userSvc := userservice.NewService(mockRepo)

	ctx := context.Background()
	email := "test@example.com"

	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("database connection failed"))

	response, err := userSvc.LookupUser(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "database connection failed", err.Error())
	mockRepo.AssertExpectations(t)
}