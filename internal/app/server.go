package app

import (
	"fmt"

	"example.com/internal/infrastructure/config"
	"example.com/internal/infrastructure/logger"
	"example.com/internal/interfaces/api"
	"example.com/internal/interfaces/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type Server struct {
	engine *gin.Engine
	config *config.Config
	logger logger.Logger
}

func NewServer(container *dig.Container) (*Server, error) {
	var cfg *config.Config
	var log logger.Logger
	var authHandler *api.AuthHandler

	if err := container.Invoke(func(
		c *config.Config,
		l logger.Logger,
		ah *api.AuthHandler,
	) {
		cfg = c
		log = l
		authHandler = ah
	}); err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	// Middleware
	engine.Use(middleware.Session(cfg.JWT.SecretKey))
	engine.Use(middleware.CSRF(cfg.JWT.SecretKey))

	// Routes
	setupRoutes(engine, authHandler)

	return &Server{
		engine: engine,
		config: cfg,
		logger: log,
	}, nil
}

func (s *Server) Run() error {
	addr := ":" + s.config.Server.Port
	s.logger.Info("Starting server", "address", addr, "env", s.config.Server.Env)

	return s.engine.Run(addr)
}

func setupRoutes(engine *gin.Engine, authHandler *api.AuthHandler) {
	// CSRF token endpoint
	engine.GET("/csrf-token", middleware.CSRFToken())

	// API routes
	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/signup", authHandler.SignUp)
				auth.POST("/login", authHandler.Login)
			}
		}
	}

	// Legacy auth routes for backward compatibility
	engine.POST("/auth/user/signup", authHandler.SignUp)
	engine.POST("/auth/user/login", authHandler.Login)
}
