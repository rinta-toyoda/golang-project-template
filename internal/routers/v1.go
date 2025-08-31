package routers

import (
	middlewares "example.com/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerApiRoutes(server *gin.Engine, db *gorm.DB) {

	v1Route := server.Group("/api/v1")
	v1Route.Use(middlewares.Authenticate)

}
