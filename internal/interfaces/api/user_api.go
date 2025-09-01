package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	v1api "example.com/gen/openapi/v1/go"
	userservice "example.com/internal/domain/service/v1"
	"example.com/internal/infrastructure/logger"
)

// UserAPIHandler extends the generated UserLoginAPIAPI with actual business logic
type UserAPIHandler struct {
	*v1api.UserLoginAPIAPI
	userService userservice.Service
	logger      logger.Logger
}

// NewUserAPIHandler creates a new user API handler that extends the generated API
func NewUserAPIHandler(userService userservice.Service, logger logger.Logger) *UserAPIHandler {
	return &UserAPIHandler{
		UserLoginAPIAPI: &v1api.UserLoginAPIAPI{},
		userService:     userService,
		logger:          logger,
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

	response, err := h.userService.LookupUser(c.Request.Context(), email)
	if err != nil {
		h.logger.Warn("User lookup failed", "error", err.Error(), "email", email)

		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, v1api.Error{Message: "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, v1api.Error{Message: "Internal server error"})
		}
		return
	}

	h.logger.Info("User lookup successful", "email", email, "username", response.Username)
	c.JSON(http.StatusOK, response)
}
