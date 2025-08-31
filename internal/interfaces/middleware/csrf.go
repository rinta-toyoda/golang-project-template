package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func CSRF(secret string) gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: secret,
		ErrorFunc: func(c *gin.Context) {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token validation failed"})
			c.Abort()
		},
	})
}

func CSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.GetToken(c)
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func RequireXSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerToken := c.GetHeader("X-XSRF-TOKEN")
		if headerToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "X-XSRF-TOKEN header is required"})
			c.Abort()
			return
		}

		cookieToken := csrf.GetToken(c)
		if headerToken != cookieToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "XSRF token mismatch"})
			c.Abort()
			return
		}

		c.Next()
	}
}
