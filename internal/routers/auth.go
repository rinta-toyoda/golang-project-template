package routers

import (
	"example.com/internal/handler/auth"
	"example.com/internal/repository"
	"example.com/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerAuthRoutes(server *gin.Engine, db *gorm.DB) {
	authRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(authRepository)
	authHandler := auth.NewAuthHandler(authService)

	authRoute := server.Group("/auth")
	{
		authRoute.POST("/login", authHandler.Login)
		authRoute.POST("/signup", authHandler.Signup)
	}
}
