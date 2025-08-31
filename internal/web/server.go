package web

import (
	"example.com/internal/routers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewServer(db *gorm.DB) *gin.Engine {
	server := gin.Default()

	// Attach middleware
	AttachSession(server)
	AttachCSRF(server)

	// Setup routes
	routers.SetupRouter(server, db)

	return server
}
