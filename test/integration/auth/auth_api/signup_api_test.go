package login_api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authapi "example.com/gen/openapi/auth/go"
	"example.com/internal/domain/entity"
	authservice "example.com/internal/domain/service/auth"
	authusecase "example.com/internal/domain/usecase/auth"
	"example.com/internal/infrastructure/logger"
	"example.com/internal/interfaces/api"
	"example.com/test/unit/mocks"
)

func setupSignupRouter() (*gin.Engine, *mocks.MockUserRepository, *mocks.MockPasswordHasher) {
	gin.SetMode(gin.TestMode)

	mockRepo := &mocks.MockUserRepository{}
	mockHasher := &mocks.MockPasswordHasher{}
	authSvc := authservice.NewService(mockRepo, mockHasher)
	signupUseCase := authusecase.NewSignupUseCase(authSvc)
	loginUseCase := authusecase.NewLoginUseCase(authSvc)
	testLogger := logger.New("test")

	authAPIHandler := api.NewAuthAPIHandler(signupUseCase, loginUseCase, testLogger)

	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/signup", authAPIHandler.UserSignup)
	}

	return router, mockRepo, mockHasher
}

func TestSignupAPI_Success(t *testing.T) {
	router, mockRepo, mockHasher := setupSignupRouter()

	signupReq := authapi.SignupRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock that user doesn't exist
	mockRepo.On("FindByUserName", mock.Anything, "testuser").Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("user not found"))

	// Mock successful password hashing
	mockHasher.On("Hash", "password123").Return("hashed_password", nil)

	// Mock successful user creation
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	body, _ := json.Marshal(signupReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response authapi.SignupResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", response.User.Email)
	assert.Equal(t, "testuser", response.User.Username)
	assert.Equal(t, "User created successfully", response.Message)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestSignupAPI_UserAlreadyExists_Email(t *testing.T) {
	router, mockRepo, _ := setupSignupRouter()

	existingUser := &entity.User{
		ID:       "existing-user-123",
		Email:    "test@example.com",
		UserName: "existinguser",
	}

	signupReq := authapi.SignupRequest{
		Email:    "test@example.com",
		Username: "newuser",
		Password: "password123",
	}

	// Mock that username doesn't exist but email already exists
	mockRepo.On("FindByUserName", mock.Anything, "newuser").Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

	body, _ := json.Marshal(signupReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "User already exists", errorResp.Message)

	mockRepo.AssertExpectations(t)
}

func TestSignupAPI_UserAlreadyExists_Username(t *testing.T) {
	router, mockRepo, _ := setupSignupRouter()

	existingUser := &entity.User{
		ID:       "existing-user-123",
		Email:    "existing@example.com",
		UserName: "testuser",
	}

	signupReq := authapi.SignupRequest{
		Email:    "new@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock that username already exists
	mockRepo.On("FindByUserName", mock.Anything, "testuser").Return(existingUser, nil)

	body, _ := json.Marshal(signupReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "User already exists", errorResp.Message)

	mockRepo.AssertExpectations(t)
}

func TestSignupAPI_InvalidJSON(t *testing.T) {
	router, _, _ := setupSignupRouter()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request format", errorResp.Message)
}

func TestSignupAPI_HashingError(t *testing.T) {
	router, mockRepo, mockHasher := setupSignupRouter()

	signupReq := authapi.SignupRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock that user doesn't exist
	mockRepo.On("FindByUserName", mock.Anything, "testuser").Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("user not found"))

	// Mock password hashing failure
	mockHasher.On("Hash", "password123").Return("", errors.New("hashing failed"))

	body, _ := json.Marshal(signupReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Internal server error", errorResp.Message)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}

func TestSignupAPI_DatabaseError(t *testing.T) {
	router, mockRepo, mockHasher := setupSignupRouter()

	signupReq := authapi.SignupRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	// Mock that user doesn't exist
	mockRepo.On("FindByUserName", mock.Anything, "testuser").Return(nil, errors.New("user not found"))
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("user not found"))

	// Mock successful password hashing
	mockHasher.On("Hash", "password123").Return("hashed_password", nil)

	// Mock database error during creation
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.New("database connection failed"))

	body, _ := json.Marshal(signupReq)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResp authapi.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Internal server error", errorResp.Message)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
}
