package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authapi "example.com/gen/openapi/auth/go"
	authservice "example.com/internal/domain/service/auth"
	"example.com/internal/infrastructure/logger"
)

// AuthAPIHandler extends the generated AuthUserAPI with actual business logic
type AuthAPIHandler struct {
	*authapi.AuthUserAPI
	authService authservice.Service
	logger      logger.Logger
}

// NewAuthAPIHandler creates a new auth API handler that extends the generated API
func NewAuthAPIHandler(authService authservice.Service, logger logger.Logger) *AuthAPIHandler {
	return &AuthAPIHandler{
		AuthUserAPI: &authapi.AuthUserAPI{},
		authService: authService,
		logger:      logger,
	}
}

// UserLogin handles user login with proper business logic
func (h *AuthAPIHandler) UserLogin(c *gin.Context) {
	var req authapi.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", "error", err.Error())
		c.JSON(http.StatusBadRequest, authapi.Error{Message: "Invalid request format"})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("Failed login attempt", "error", err.Error(), "email", req.Email)
		
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, authapi.Error{Message: "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, authapi.Error{Message: "Internal server error"})
		}
		return
	}

	h.logger.Info("User logged in successfully", "user_id", response.User.Id)
	c.JSON(http.StatusOK, response)
}

// UserSignup handles user registration with proper business logic
func (h *AuthAPIHandler) UserSignup(c *gin.Context) {
	var req authapi.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid signup request", "error", err.Error())
		c.JSON(http.StatusBadRequest, authapi.Error{Message: "Invalid request format"})
		return
	}

	response, err := h.authService.SignUp(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create user", "error", err.Error(), "email", req.Email)
		
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, authapi.Error{Message: "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, authapi.Error{Message: "Internal server error"})
		}
		return
	}

	h.logger.Info("User created successfully", "user_id", response.User.Id, "email", response.User.Email)
	c.JSON(http.StatusCreated, response)
}

