package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/internal/handler/auth"
	"example.com/internal/model"
	"example.com/internal/repository"
	"example.com/internal/service"
)

func setupTestContainer(t *testing.T) (context.Context, testcontainers.Container, *gorm.DB) {
	ctx := context.Background()
	request := testcontainers.ContainerRequest{
		Image: "postgres:15",
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("could not start container: %v", err)
	}

	host, err := container.Host(ctx)
	assert.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	assert.NoError(t, err)

	dbName := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=testdb sslmode=disable",
		host, port.Port(),
	)
	db, err := gorm.Open(postgres.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	// マイグレーション
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	return ctx, container, db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo)
	h := auth.NewAuthHandler(svc)

	grp := r.Group("/auth")
	{
		grp.POST("/signup", h.Signup)
		grp.POST("/login", h.Login)
	}
	return r
}

func TestUserCanSignup(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx, pgc, db := setupTestContainer(t)
	defer func() {
		_ = pgc.Terminate(ctx)
	}()

	router := setupRouter(db)

	// Act
	signupPayload := map[string]string{
		"email":    "foo@example.com",
		"phone":    "+1234567890",
		"password": "password123",
	}
	body, _ := json.Marshal(signupPayload)
	request := httptest.NewRequest("POST", "/auth/signup", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code, "expected 201 on signup")

	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	token, ok := resp["jwt_token"].(string)
	assert.True(t, ok && token != "", "signup response should contain token")
}

func TestUserCannotSignupWithInvalidData(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx, pgc, db := setupTestContainer(t)
	defer func() {
		_ = pgc.Terminate(ctx)
	}()

	router := setupRouter(db)

	// Act
	tests := []struct {
		name        string
		payload     map[string]string
		status      int
		description string
	}{
		{
			"Invalid Phone Number",
			map[string]string{
				"email":    "foo@example.com",
				"phone":    "some phone number", // Invalid phone number
				"password": "password123",
			},
			http.StatusBadRequest,
			"expected 400 on signup with invalid phone number",
		},
		{
			"Invalid Email",
			map[string]string{
				"email":    "", // Invalid phone number
				"phone":    "+1234567890",
				"password": "password123",
			},
			http.StatusBadRequest,
			"expected 400 on signup with invalid email",
		},
		{
			"Invalid Password",
			map[string]string{
				"email":    "",
				"phone":    "+1234567890",
				"password": "0", // Invalid password (too short)
			},
			http.StatusBadRequest,
			"expected 400 on signup with invalid password",
		},
		{
			"Missing Key",
			map[string]string{
				"phone":    "+1234567890",
				"password": "password1234", // Invalid password (too short)
			},
			http.StatusBadRequest,
			"expected 400 on signup with missing key",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			body, _ := json.Marshal(test.payload)
			request := httptest.NewRequest("POST", "/auth/signup", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code, "expected 400 on signup with invalid data")
		})
	}
}

func TestUserCanLogin(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx, container, db := setupTestContainer(t)
	defer func() {
		_ = container.Terminate(ctx)
	}()

	router := setupRouter(db)

	func() {
		signupPayload := map[string]string{
			"email":    "foo@example.com",
			"phone":    "+1234567890",
			"password": "password123",
		}
		body, _ := json.Marshal(signupPayload)
		request := httptest.NewRequest("POST", "/auth/signup", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		assert.Equal(t, http.StatusCreated, w.Code)
	}()

	// Act
	loginPayload := map[string]string{
		"email":    "foo@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(loginPayload)
	request := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, request)

	// Assert
	assert.Equal(t, http.StatusOK, w2.Code, "expected 200 on login")

	var resp2 map[string]interface{}
	assert.NoError(t, json.Unmarshal(w2.Body.Bytes(), &resp2))
	token, ok := resp2["jwt_token"].(string)
	assert.True(t, ok && token != "", "login response should contain token")
}

func TestUserCannotLoginWithInvalidData(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx, container, db := setupTestContainer(t)
	defer func() {
		_ = container.Terminate(ctx)
	}()

	router := setupRouter(db)

	// Create a user for login tests
	func() {
		signupPayload := map[string]string{
			"email":    "foo@example.com",
			"phone":    "+1234567890",
			"password": "password123",
		}
		body, _ := json.Marshal(signupPayload)
		request := httptest.NewRequest("POST", "/auth/signup", bytes.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, request)
		assert.Equal(t, http.StatusCreated, w.Code)
	}()

	// Act
	tests := []struct {
		name        string
		payload     map[string]string
		status      int
		description string
	}{
		{
			"Wrong Password",
			map[string]string{
				"email":    "foo@example.com",
				"password": "wrong password", // Invalid password
			},
			http.StatusBadRequest,
			"expected 400 on login with wrong password",
		},
		{
			"Wrong Email",
			map[string]string{
				"email":    "wrongemail@example.com", // Invalid Email
				"password": "password123",
			},
			http.StatusBadRequest,
			"expected 400 on login with wrong email",
		},
		{
			"Invalid Password",
			map[string]string{
				"password": "password123", // Missing Key
			},
			http.StatusBadRequest,
			"expected 400 on login with missing key",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			body, _ := json.Marshal(test.payload)
			request := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, request)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code, "expected 400 on login with invalid data")
		})
	}
}
