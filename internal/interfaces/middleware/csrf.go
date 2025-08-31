package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func CSRF(secret string) gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: secret,
		ErrorFunc: func(c *gin.Context) {
			c.JSON(403, gin.H{"error": "CSRF token validation failed"})
			c.Abort()
		},
	})
}

func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.GetToken(c)
		c.JSON(200, gin.H{"token": token})
	}
}
