package utils

import (
	middlewares "example.com/internal/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SaveUserSession(c *gin.Context, userId string) error {
	sess := sessions.Default(c)
	sess.Set(middlewares.SessionKeyUserID, userId)
	return sess.Save()
}
