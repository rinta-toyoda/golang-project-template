package login_api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authapi "example.com/gen/openapi/auth/go"
	"example.com/internal/domain/entity"
	authservice "example.com/internal/domain/service/auth"
	"example.com/internal/infrastructure/logger"
	"example.com/internal/interfaces/api"
	"example.com/test/unit/mocks"
)

func setupLoginRouter() (*gin.Engine, *mocks.MockUserRepository, *mocks.MockPasswordHasher) {
	gin.SetMode(gin.TestMode)
	
	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	testLogger := logger.New("test")
	
	authAPIHandler := api.NewAuthAPIHandler(authSvc, testLogger)
	
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/login", authAPIHandler.UserLogin)
	}
	
	return router, mockRepo, mockHasher
}

func TestLoginAPI_Success(t *testing.T) {
	router, mockRepo, mockHasher := setupLoginRouter()
	
	now := time.Now()
	user := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		UserName:     "testuser",
		PasswordHash: "hashed_password",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	loginReq := authapi.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	mockRepo.On("FindByUserNameOrEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockHasher.On("Verify", "password123", "hashed_password").Return(true)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
	
	body, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response authapi.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", response.User.Id)
	assert.Equal(t, "test@example.com", response.User.Email)
	assert.Equal(t, "testuser", response.User.Username)
	assert.Equal(t, "Login successful", response.Message)
	
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLoginAPI_InvalidCredentials_WrongPassword(t *testing.T) {
	router, mockRepo, mockHasher := setupLoginRouter()
	
	user := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	
	loginReq := authapi.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	
	mockRepo.On("FindByUserNameOrEmail", mock.Anything, "test@example.com").Return(user, nil)
	mockHasher.On("Verify", "wrongpassword", "hashed_password").Return(false)
	
	body, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid credentials", errorResp.Message)
	
	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestLoginAPI_InvalidCredentials_UserNotFound(t *testing.T) {
	router, mockRepo, _ := setupLoginRouter()
	
	loginReq := authapi.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}
	
	mockRepo.On("FindByUserNameOrEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("user not found"))
	
	body, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid credentials", errorResp.Message)
	
	mockRepo.AssertExpectations(t)
}

func TestLoginAPI_InvalidJSON(t *testing.T) {
	router, _, _ := setupLoginRouter()
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request format", errorResp.Message)
}