package app

import (
	"go.uber.org/dig"
	"gorm.io/gorm"

	"example.com/internal/domain/repository"
	authservice "example.com/internal/domain/service/auth"
	userservice "example.com/internal/domain/service/v1"
	authusecase "example.com/internal/domain/usecase/auth"
	userusecase "example.com/internal/domain/usecase/v1"
	"example.com/internal/infrastructure/config"
	"example.com/internal/infrastructure/database"
	"example.com/internal/infrastructure/logger"
	"example.com/internal/interfaces/api"
	"example.com/pkg/security"
)

func BuildContainer() (*dig.Container, error) {
	container := dig.New()

	// Configuration
	if err := container.Provide(config.Load); err != nil {
		return nil, err
	}

	// Logger
	if err := container.Provide(func(cfg *config.Config) logger.Logger {
		return logger.New(cfg.Server.Env)
	}); err != nil {
		return nil, err
	}

	// Database
	if err := container.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		if cfg.Database.URL != "" {
			return database.ConnectFromURL(cfg.Database.URL)
		}

		dbCfg := database.Config{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			DBName:   cfg.Database.DBName,
			SSLMode:  cfg.Database.SSLMode,
		}
		return database.Connect(dbCfg)
	}); err != nil {
		return nil, err
	}

	// Security
	if err := container.Provide(security.NewBcryptHasher); err != nil {
		return nil, err
	}

	// Repositories
	if err := container.Provide(func(db *gorm.DB) repository.UserRepository {
		return database.NewUserRepository(db)
	}); err != nil {
		return nil, err
	}

	// Services
	if err := container.Provide(authservice.NewService); err != nil {
		return nil, err
	}
	if err := container.Provide(userservice.NewService); err != nil {
		return nil, err
	}

	// Use Cases
	if err := container.Provide(authusecase.NewSignupUseCase); err != nil {
		return nil, err
	}
	if err := container.Provide(authusecase.NewLoginUseCase); err != nil {
		return nil, err
	}
	if err := container.Provide(userusecase.NewUserLookupUseCase); err != nil {
		return nil, err
	}

	// API Handlers
	if err := container.Provide(api.NewAuthAPIHandler); err != nil {
		return nil, err
	}
	if err := container.Provide(api.NewUserAPIHandler); err != nil {
		return nil, err
	}

	return container, nil
}
