package app

import (
	"example.com/internal/domain/repository"
	"example.com/internal/domain/service"
	"example.com/internal/infrastructure/config"
	"example.com/internal/infrastructure/database"
	"example.com/internal/infrastructure/logger"
	"example.com/internal/interfaces/api"
	"example.com/pkg/security"
	"go.uber.org/dig"
	"gorm.io/gorm"
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
	if err := container.Provide(service.NewAuthService); err != nil {
		return nil, err
	}

	// Handlers
	if err := container.Provide(api.NewAuthHandler); err != nil {
		return nil, err
	}

	return container, nil
}
