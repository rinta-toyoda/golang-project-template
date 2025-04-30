package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(server *gin.Engine, db *gorm.DB) {
	registerAuthRoutes(server, db)
	registerApiRoutes(server, db)
}
