package routers

import (
	middlewares "example.com/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerApiRoutes(server *gin.Engine, db *gorm.DB) {

	v1 := server.Group("/api/v1")
	v1.Use(middlewares.Authenticate)
	// v1.GET("/api/v1/todo", todoHandler.GetAll)
	// v1.POST("/api/v1/todo", todoHandler.Create)
	// v1.GET("/api/v1/todo/:id", todoHandler.GetById)
	// v1.PUT("/api/v1/todo/:id", todoHandler.Update)
	// v1.DELETE("/api/v1/todo/:id", todoHandler.Delete)
}
