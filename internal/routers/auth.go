package routers

import (
	"example.com/internal/handler/auth"
	"example.com/internal/repository"
	"example.com/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerAuthRoutes(server *gin.Engine, db *gorm.DB) {
	// Setup user components
	userRepository := repository.NewUserRepository(db)
	userAuthService := service.NewUserAuthService(userRepository)
	userAuthHandler := auth.NewUserAuthHandler(userAuthService)

	// Setup organization components
	organizationRepository := repository.NewOrganizationRepository(db)
	organizationAuthService := service.NewOrganizationAuthService(organizationRepository)
	organizationAuthHandler := auth.NewOrganizationAuthHandler(organizationAuthService)

	authRoute := server.Group("/auth")
	{
		userAuthRoute := authRoute.Group("/user")
		{
			userAuthRoute.POST("/signup", userAuthHandler.Signup)
			userAuthRoute.POST("/login", userAuthHandler.Login)
		}
		organizationAuthRoute := authRoute.Group("/organization")
		{
			organizationAuthRoute.POST("/signup", organizationAuthHandler.Signup)
			organizationAuthRoute.POST("/login", organizationAuthHandler.Login)
		}
	}
}
