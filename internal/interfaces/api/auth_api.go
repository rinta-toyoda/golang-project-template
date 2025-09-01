package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authapi "example.com/gen/openapi/auth/go"
	authusecase "example.com/internal/domain/usecase/auth"
	"example.com/internal/infrastructure/logger"
)

// AuthAPIHandler extends the generated AuthUserAPI with actual business logic
type AuthAPIHandler struct {
	*authapi.AuthUserAPI
	signupUseCase authusecase.SignupUseCase
	loginUseCase  authusecase.LoginUseCase
	logger        logger.Logger
}

// NewAuthAPIHandler creates a new auth API handler that extends the generated API
func NewAuthAPIHandler(
	signupUseCase authusecase.SignupUseCase,
	loginUseCase authusecase.LoginUseCase,
	logger logger.Logger,
) *AuthAPIHandler {
	return &AuthAPIHandler{
		AuthUserAPI:   &authapi.AuthUserAPI{},
		signupUseCase: signupUseCase,
		loginUseCase:  loginUseCase,
		logger:        logger,
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

	user, err := h.loginUseCase.Call(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Failed login attempt", "error", err.Error(), "email", req.Email)

		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, authapi.Error{Message: "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, authapi.Error{Message: "Internal server error"})
		}
		return
	}

	// Convert domain model to API response
	apiUser := authapi.User{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.UserName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.LastLoginAt != nil {
		apiUser.LastLoginAt = *user.LastLoginAt
	}

	response := authapi.LoginResponse{
		User:    apiUser,
		Message: "Login successful",
	}

	h.logger.Info("User logged in successfully", "user_id", user.ID)
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

	user, err := h.signupUseCase.Call(c.Request.Context(), req.Email, req.Password, req.Username)
	if err != nil {
		h.logger.Error("Failed to create user", "error", err.Error(), "email", req.Email)

		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, authapi.Error{Message: "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, authapi.Error{Message: "Internal server error"})
		}
		return
	}

	// Convert domain model to API response
	apiUser := authapi.User{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.UserName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.LastLoginAt != nil {
		apiUser.LastLoginAt = *user.LastLoginAt
	}

	response := authapi.SignupResponse{
		User:    apiUser,
		Message: "User created successfully",
	}

	h.logger.Info("User created successfully", "user_id", user.ID, "email", user.Email)
	c.JSON(http.StatusCreated, response)
}
