package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"example.com/internal/domain/entity"
	authservice "example.com/internal/domain/service/auth"
	authusecase "example.com/internal/domain/usecase/auth"
	"example.com/test/unit/mocks"
)

func TestSignupUseCase_Call_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewSignupUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"
	hashedPassword := "hashed_password"

	// Mock that user doesn't exist
	mockRepo.On("FindByUserName", ctx, username).Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found"))

	// Mock successful password hashing
	mockHasher.On("Hash", password).Return(hashedPassword, nil)

	// Mock successful user creation
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, err := useCase.Call(ctx, email, password, username)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, username, user.UserName)
	assert.Equal(t, hashedPassword, user.PasswordHash)
	assert.NotEmpty(t, user.ID)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestSignupUseCase_Call_UserAlreadyExists_Email(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewSignupUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"

	existingUser := &entity.User{
		ID:       "existing-user",
		Email:    email,
		UserName: "existinguser",
	}

	// Mock that username doesn't exist but email already exists
	mockRepo.On("FindByUserName", ctx, username).Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	user, err := useCase.Call(ctx, email, password, username)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, authservice.ErrUserAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestSignupUseCase_Call_UserAlreadyExists_Username(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewSignupUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"

	existingUser := &entity.User{
		ID:       "existing-user",
		Email:    "existing@example.com",
		UserName: username,
	}

	// Mock that username already exists
	mockRepo.On("FindByUserName", ctx, username).Return(existingUser, nil)

	user, err := useCase.Call(ctx, email, password, username)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, authservice.ErrUserAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestSignupUseCase_Call_HashingError(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewSignupUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"

	// Mock that user doesn't exist
	mockRepo.On("FindByUserName", ctx, username).Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found"))

	// Mock password hashing failure
	mockHasher.On("Hash", password).Return("", errors.New("hashing failed"))

	user, err := useCase.Call(ctx, email, password, username)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "hashing failed", err.Error())
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestSignupUseCase_Call_DatabaseError(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	useCase := authusecase.NewSignupUseCase(authSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"
	hashedPassword := "hashed_password"

	// Mock that user doesn't exist
	mockRepo.On("FindByUserName", ctx, username).Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found"))

	// Mock successful password hashing
	mockHasher.On("Hash", password).Return(hashedPassword, nil)

	// Mock database error during creation
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("database connection failed"))

	user, err := useCase.Call(ctx, email, password, username)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "database connection failed", err.Error())
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
