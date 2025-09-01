package app

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/dig"

	authapi "example.com/gen/openapi/auth/go"
	"example.com/internal/infrastructure/config"
	"example.com/internal/infrastructure/logger"
	"example.com/internal/interfaces/middleware"
)

type Server struct {
	engine *gin.Engine
	config *config.Config
	logger logger.Logger
}

func NewServer(container *dig.Container) (*Server, error) {
	var cfg *config.Config
	var log logger.Logger
	var authUserAPI *authapi.AuthUserAPI

	if err := container.Invoke(func(
		c *config.Config,
		l logger.Logger,
		aua *authapi.AuthUserAPI,
	) {
		cfg = c
		log = l
		authUserAPI = aua
	}); err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	// CORS middleware for Swagger UI
	if cfg.Server.Env != "production" {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowHeaders = []string{
			"Origin", "Content-Type", "Accept", "X-XSRF-TOKEN",
			"X-Requested-With", "X-CSRF-Token", "Authorization",
		}
		corsConfig.AllowCredentials = true
		corsConfig.ExposeHeaders = []string{"Set-Cookie", "X-CSRF-Token"}
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		engine.Use(cors.New(corsConfig))
	}

	// Middleware
	engine.Use(middleware.Session(cfg.Security.SessionSecret))
	engine.Use(middleware.CSRF(cfg.Security.CSRFSecret))

	// Routes
	setupRoutes(engine, authUserAPI)

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

func setupRoutes(engine *gin.Engine, authUserAPI *authapi.AuthUserAPI) {
	// Serve OpenAPI specs first
	engine.Static("/api/auth", "./api/auth")
	engine.Static("/api/v1", "./api/v1")

	// Swagger UI endpoints
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/auth/openapi.yaml")))

	// CSRF token endpoint
	engine.GET("/csrf-token", middleware.CSRFToken())

	// API routes with XSRF protection
	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			auth.Use(middleware.RequireXSRF())
			{
				auth.POST("/signup", authUserAPI.UserSignup)
				auth.POST("/login", authUserAPI.UserLogin)
				auth.GET("/user/lookup", authUserAPI.UserLookup)
			}
		}
	}

	// Legacy auth routes for backward compatibility
	legacyAuth := engine.Group("/auth/user")
	legacyAuth.Use(middleware.RequireXSRF())
	{
		legacyAuth.POST("/signup", authUserAPI.UserSignup)
		legacyAuth.POST("/login", authUserAPI.UserLogin)
		legacyAuth.GET("/lookup", authUserAPI.UserLookup)
	}
}
