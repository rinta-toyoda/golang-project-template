package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"example.com/internal/domain/entity"
	userservice "example.com/internal/domain/service/v1"
	userusecase "example.com/internal/domain/usecase/v1"
	"example.com/test/unit/mocks"
)

func TestUserLookupUseCase_Call_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	userSvc := userservice.NewService(mockRepo)
	useCase := userusecase.NewUserLookupUseCase(userSvc)

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

	user, err := useCase.Call(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserLookupUseCase_Call_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	userSvc := userservice.NewService(mockRepo)
	useCase := userusecase.NewUserLookupUseCase(userSvc)

	ctx := context.Background()
	email := "notfound@example.com"

	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found"))

	user, err := useCase.Call(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUserLookupUseCase_Call_DatabaseError(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	userSvc := userservice.NewService(mockRepo)
	useCase := userusecase.NewUserLookupUseCase(userSvc)

	ctx := context.Background()
	email := "test@example.com"

	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("database connection failed"))

	user, err := useCase.Call(ctx, email)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "database connection failed", err.Error())
	mockRepo.AssertExpectations(t)
}
