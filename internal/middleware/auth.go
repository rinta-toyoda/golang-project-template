package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SessionKeyUserID = "userId"
	ContextUserId    = "userId"
)

func Authenticate(c *gin.Context) {
	sess := sessions.Default(c)

	userId := sess.Get(SessionKeyUserID)
	if userId == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Authentication required",
		})
		return
	}

	c.Set(ContextUserId, userId)
	c.Next()
}
