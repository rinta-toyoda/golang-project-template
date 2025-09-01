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
	"example.com/test/unit/mocks"
)

func TestAuthService_CheckUserExists_NoConflict(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"

	mockRepo.On("FindByUserName", ctx, username).Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found"))

	err := authSvc.CheckUserExists(ctx, email, username)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_CheckUserExists_EmailExists(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	email := "test@example.com"
	username := "testuser"

	existingUser := &entity.User{ID: "existing-user"}
	mockRepo.On("FindByUserName", ctx, username).Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", ctx, email).Return(existingUser, nil)

	err := authSvc.CheckUserExists(ctx, email, username)

	assert.Error(t, err)
	assert.Equal(t, authservice.ErrUserAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_CreateUser_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"
	hashedPassword := "hashed_password"

	mockHasher.On("Hash", password).Return(hashedPassword, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, err := authSvc.CreateUser(ctx, email, password, username)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, username, user.UserName)
	assert.Equal(t, hashedPassword, user.PasswordHash)
	assert.NotEmpty(t, user.ID)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestAuthService_CreateUser_HashError(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	username := "testuser"

	mockHasher.On("Hash", password).Return("", errors.New("hash failed"))

	user, err := authSvc.CreateUser(ctx, email, password, username)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "hash failed", err.Error())
	mockHasher.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

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

	mockRepo.On("FindByUserNameOrEmail", ctx, email).Return(existingUser, nil)
	mockHasher.On("Verify", password, "hashed_password").Return(true)

	user, err := authSvc.AuthenticateUser(ctx, email, password)

	assert.NoError(t, err)
	assert.Equal(t, existingUser, user)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_UserNotFound(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	email := "notfound@example.com"
	password := "password123"

	mockRepo.On("FindByUserNameOrEmail", ctx, email).Return(nil, errors.New("user not found"))

	user, err := authSvc.AuthenticateUser(ctx, email, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, authservice.ErrInvalidCredentials, err)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_AuthenticateUser_InvalidPassword(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	email := "test@example.com"
	password := "wrongpassword"

	existingUser := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	mockRepo.On("FindByUserNameOrEmail", ctx, email).Return(existingUser, nil)
	mockHasher.On("Verify", password, "hashed_password").Return(false)

	user, err := authSvc.AuthenticateUser(ctx, email, password)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, authservice.ErrInvalidCredentials, err)
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestAuthService_UpdateLastLogin_Success(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)

	ctx := context.Background()
	userID := "user-123"

	existingUser := &entity.User{
		ID:           userID,
		Email:        "test@example.com",
		UserName:     "testuser",
		PasswordHash: "hashed_password",
	}

	mockRepo.On("FindByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(user *entity.User) bool {
		return user.ID == userID && user.LastLoginAt != nil
	})).Return(nil)

	err := authSvc.UpdateLastLogin(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
