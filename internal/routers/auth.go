package routers

import (
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"

	authapi "example.com/internal/adapter/authapi"
	"example.com/internal/repository"
	"example.com/internal/service"
	"example.com/internal/usecase/userauth"
)

func registerAuthRoutes(server *gin.Engine, db *gorm.DB) {
	handler := prepareUserAuthHandler(db)

	server.GET("/csrf-token", func(c *gin.Context) {
		token := csrf.GetToken(c)
		// align with OpenAPI: { "token": "..." }
		c.JSON(200, gin.H{"token": token})
	})

	authRoute := server.Group("/auth")
	{
		userAuth := authRoute.Group("/user")
		{
			userAuth.POST("/signup", handler.Signup)
			userAuth.POST("/login", handler.Login)
		}
	}
}

// DI
func prepareUserAuthHandler(db *gorm.DB) *authapi.UserAPIImpl {
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserAuthService(userRepository)

	signupUseCase := userauth.NewUserAuthSignupUseCase(userService)
	loginUseCase := userauth.NewUserAuthLoginUseCase(userService)

	return authapi.NewUserAPIImpl(signupUseCase, loginUseCase)
}
