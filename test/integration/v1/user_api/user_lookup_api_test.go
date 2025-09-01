package user_lookup_api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	v1api "example.com/gen/openapi/v1/go"
	"example.com/internal/domain/entity"
	userservice "example.com/internal/domain/service/v1"
	"example.com/internal/infrastructure/logger"
	"example.com/internal/interfaces/api"
	"example.com/test/unit/mocks"
)

func setupUserLookupRouter() (*gin.Engine, *mocks.MockUserRepository) {
	gin.SetMode(gin.TestMode)
	
	mockRepo := &mocks.MockUserRepository{}
	userSvc := userservice.NewService(mockRepo)
	testLogger := logger.New("test")
	
	userAPIHandler := api.NewUserAPIHandler(userSvc, testLogger)
	
	router := gin.New()
	user := router.Group("/user")
	{
		user.GET("/lookup", userAPIHandler.UserLookup)
	}
	
	return router, mockRepo
}

func TestUserLookupAPI_Success(t *testing.T) {
	router, mockRepo := setupUserLookupRouter()
	
	user := &entity.User{
		ID:        "user-123",
		Email:     "test@example.com",
		UserName:  "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user/lookup?email=test@example.com", nil)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response v1api.UserLookupResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", response.Username)
	assert.Equal(t, "test@example.com", response.Email)
	
	mockRepo.AssertExpectations(t)
}

func TestUserLookupAPI_UserNotFound(t *testing.T) {
	router, mockRepo := setupUserLookupRouter()
	
	mockRepo.On("FindByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("user not found"))
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user/lookup?email=notfound@example.com", nil)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var errorResp v1api.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", errorResp.Message)
	
	mockRepo.AssertExpectations(t)
}

func TestUserLookupAPI_MissingEmailParameter(t *testing.T) {
	router, _ := setupUserLookupRouter()
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user/lookup", nil)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var errorResp v1api.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Email parameter is required", errorResp.Message)
}

func TestUserLookupAPI_DatabaseError(t *testing.T) {
	router, mockRepo := setupUserLookupRouter()
	
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("database connection failed"))
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/user/lookup?email=test@example.com", nil)
	
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var errorResp v1api.Error
	err := json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.NoError(t, err)
	assert.Equal(t, "Internal server error", errorResp.Message)
	
	mockRepo.AssertExpectations(t)
}