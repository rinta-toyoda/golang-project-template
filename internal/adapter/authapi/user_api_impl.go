// internal/adapter/authapi/user_api_impl.go
package authapi_adapter

import (
	"example.com/internal/usecase/userauth"
	"net/http"

	"github.com/gin-gonic/gin"

	authapi "example.com/gen/openapi/auth/go" // <- generated models/apis
)

// Implement your own struct that will be wired in routers.
// We keep a tiny surface: call use-cases and translate models.
type UserAPIImpl struct {
	SignupUC userauth.UserAuthSignupUseCase
	LoginUC  userauth.UserAuthLoginUseCase
}

func NewUserAPIImpl(signup userauth.UserAuthSignupUseCase, login userauth.UserAuthLoginUseCase) *UserAPIImpl {
	return &UserAPIImpl{SignupUC: signup, LoginUC: login}
}

// POST /auth/user/signup
// OpenAPI: SignupRequest { email, password, username } -> SignupResponse { user?, message }
func (h *UserAPIImpl) Signup(c *gin.Context) {
	var req authapi.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, authapi.Error{Error: "bad_request", Message: err.Error()})
		return
	}

	// Map transport -> domain (your UC signature: (ctx, username, email, password))
	if err := h.SignupUC.Call(c, req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, authapi.Error{Error: "signup_failed", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, authapi.SignupResponse{
		Message: "account created",
		// Optionally include a user object if your UC returns one
	})
}

// POST /auth/user/login
// OpenAPI: LoginRequest { email, password } -> LoginResponse { user?, message }
func (h *UserAPIImpl) Login(c *gin.Context) {
	var req authapi.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, authapi.Error{Error: "bad_request", Message: err.Error()})
		return
	}

	// Your existing UC used "identifier"; with OpenAPI we standardize to email.
	if err := h.LoginUC.Call(c, req.Email, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, authapi.Error{Error: "invalid_credentials", Message: "failed to login"})
		return
	}

	c.JSON(http.StatusOK, authapi.LoginResponse{
		Message: "logged in",
		// Optionally include a user object if available
	})
}
