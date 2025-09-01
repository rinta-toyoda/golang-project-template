package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1api "example.com/gen/openapi/v1/go"
	userusecase "example.com/internal/domain/usecase/v1"
	"example.com/internal/infrastructure/logger"
)

// UserAPIHandler extends the generated UserLoginAPIAPI with actual business logic
type UserAPIHandler struct {
	*v1api.UserLoginAPIAPI
	userLookupUseCase userusecase.UserLookupUseCase
	logger            logger.Logger
}

// NewUserAPIHandler creates a new user API handler that extends the generated API
func NewUserAPIHandler(userLookupUseCase userusecase.UserLookupUseCase, logger logger.Logger) *UserAPIHandler {
	return &UserAPIHandler{
		UserLoginAPIAPI:   &v1api.UserLoginAPIAPI{},
		userLookupUseCase: userLookupUseCase,
		logger:            logger,
	}
}

// UserLookup handles user lookup by email
func (h *UserAPIHandler) UserLookup(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		h.logger.Warn("Missing email parameter in lookup request")
		c.JSON(http.StatusBadRequest, v1api.Error{Message: "Email parameter is required"})
		return
	}

	user, err := h.userLookupUseCase.Call(c.Request.Context(), email)
	if err != nil {
		h.logger.Warn("User lookup failed", "error", err.Error(), "email", email)

		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, v1api.Error{Message: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, v1api.Error{Message: "Internal server error"})
		}
		return
	}

	// Convert domain model to API response
	response := v1api.UserLookupResponse{
		Username: user.UserName,
		Email:    user.Email,
	}

	h.logger.Info("User lookup successful", "email", email, "username", user.UserName)
	c.JSON(http.StatusOK, response)
}
