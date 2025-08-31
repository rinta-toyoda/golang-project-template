package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authapi "example.com/gen/openapi/auth/go"
	"example.com/internal/domain/service"
	"example.com/internal/infrastructure/logger"
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

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req authapi.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid signup request", "error", err.Error())
		c.JSON(http.StatusBadRequest, authapi.Error{Message: "Invalid request format"})
		return
	}

	response, err := h.authService.SignUp(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create user", "error", err.Error(), "email", req.Email)

		switch err {
		case service.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, authapi.Error{Message: "User already exists"})
		default:
			c.JSON(http.StatusInternalServerError, authapi.Error{Message: "Internal server error"})
		}
		return
	}

	h.logger.Info("User created successfully", "user_id", response.User.Id, "email", response.User.Email)
	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authapi.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request", "error", err.Error())
		c.JSON(http.StatusBadRequest, authapi.Error{Message: "Invalid request format"})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("Failed login attempt", "error", err.Error(), "identifier", req.Email)

		switch err {
		case service.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, authapi.Error{Message: "Invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, authapi.Error{Message: "Internal server error"})
		}
		return
	}

	h.logger.Info("User logged in successfully", "user_id", response.User.Id)
	c.JSON(http.StatusOK, response)
}
