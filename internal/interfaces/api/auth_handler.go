package api

import (
	"net/http"

	"example.com/internal/domain/service"
	"example.com/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
	logger      logger.Logger
}

func NewAuthHandler(authService service.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"omitempty,min=3,max=15"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	LastLoginAt string `json:"last_login_at,omitempty"`
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid signup request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.authService.SignUp(c.Request.Context(), service.SignUpRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Error("Failed to create user", "error", err.Error(), "email", req.Email)

		switch err {
		case service.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.Info("User created successfully", "user_id", user.ID, "email", user.Email)

	userResp := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.UserName,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    userResp,
		"message": "Account created successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	user, err := h.authService.Login(c.Request.Context(), service.LoginRequest{
		Identifier: req.Email,
		Password:   req.Password,
	})
	if err != nil {
		h.logger.Warn("Failed login attempt", "error", err.Error(), "identifier", req.Email)

		switch err {
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.Info("User logged in successfully", "user_id", user.ID)

	userResp := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.UserName,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if user.LastLoginAt != nil {
		userResp.LastLoginAt = user.LastLoginAt.Format("2006-01-02T15:04:05Z07:00")
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    userResp,
		"message": "Login successful",
	})
}
