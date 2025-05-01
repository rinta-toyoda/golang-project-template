package unit

import (
	"example.com/internal/model"
	"example.com/tests/unit/utils"
	"github.com/google/uuid"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"example.com/internal/dto"
	"example.com/internal/service"
)

// Ensure the JWT secret is set before any tests run.
func TestMain(m *testing.M) {
	err := os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	if err != nil {
		return
	}
	os.Exit(m.Run())
}

func TestUserCanSignup(t *testing.T) {
	// Arrange
	repository := utils.NewFakeUserRepository()
	authService := service.NewAuthService(repository)

	request := dto.SignupRequest{
		Email:    "foo@example.com",
		Phone:    "+1234567890",
		Password: "password123",
	}
	token, err := authService.Signup(request.Email, request.Phone, request.Password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token, "expected a JWT token on signup success")

	// Verify that the user was stored with a hashed password
	stored, err := repository.FindByEmail(request.Email)
	assert.NoError(t, err)
	assert.Equal(t, request.Phone, stored.Phone)

	err = bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte(request.Password))
	assert.NoError(t, err, "password hash should match the original password")
}

func TestUserCannotSignupWithExistingEmail(t *testing.T) {
	// arrange
	repository := utils.NewFakeUserRepository()
	authService := service.NewAuthService(repository)

	existing := &model.User{ID: uuid.NewString(), Email: "dup@example.com", Phone: "+111"}
	err := repository.Create(existing)
	assert.NoError(t, err)

	// act
	// Preload an existing user with the same email
	token, err := authService.Signup(existing.Email, "+222", "pw")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "email already registered", err.Error())
}

func TestUserCannotSignupWithExistingPhone(t *testing.T) {
	// arrange
	repository := utils.NewFakeUserRepository()
	authService := service.NewAuthService(repository)

	existing := &model.User{ID: uuid.NewString(), Email: "example@example.com", Phone: "+111"}
	err := repository.Create(existing)
	assert.NoError(t, err)

	// act
	// Preload an existing user with the same email
	token, err := authService.Signup("owner@example.com", existing.Phone, "pw")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "phone already registered", err.Error())
}

func TestUserCanLogin(t *testing.T) {
	// arrange
	repository := utils.NewFakeUserRepository()
	// Create a user with a known password hash
	password := "secret123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{ID: uuid.NewString(), Email: "user@example.com", Phone: "+1010", PasswordHash: string(hash)}
	err := repository.Create(user)
	assert.NoError(t, err)

	// act
	authService := service.NewAuthService(repository)
	token, err := authService.Login(user.Email, password)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token, "expected a JWT token on successful login")
}

func TestUserCannotLoginWithInvalidCredentials(t *testing.T) {
	// arrange
	repository := utils.NewFakeUserRepository()
	// Create a user with a known password hash
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	user := &model.User{ID: uuid.NewString(), Email: "user@example.com", Phone: "+1010", PasswordHash: string(hash)}
	err := repository.Create(user)
	assert.NoError(t, err)
	authService := service.NewAuthService(repository)

	// act
	// Wrong password
	token, err := authService.Login("user@example.com", "wrong")

	// Non-existent user
	token2, err2 := authService.Login("nope@example.com", "pw")

	// assert
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Error(t, err2)
	assert.Empty(t, token2)
}
