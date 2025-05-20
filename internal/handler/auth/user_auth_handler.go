package auth

import (
	"example.com/internal/dto"
	"example.com/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserAuthHandler interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userAuthHandler struct {
	userAuthService service.UserAuthService
}

func NewUserAuthHandler(auth service.UserAuthService) UserAuthHandler {
	return &userAuthHandler{userAuthService: auth}
}

func (h userAuthHandler) Signup(c *gin.Context) {
	var request dto.SignupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userAuthService.Signup(request.Email, request.Phone, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := dto.SignupResponse{
		Message:  "signup success",
		JWTToken: token,
	}
	c.JSON(http.StatusCreated, response)
}

func (h userAuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userAuthService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := dto.LoginResponse{
		Message:  "login successful",
		JWTToken: token,
	}
	c.JSON(http.StatusOK, resp)
}
