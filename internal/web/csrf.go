package web

import (
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
	"os"
)

func AttachCSRF(r *gin.Engine) {
	r.Use(csrf.Middleware(csrf.Options{
		Secret: os.Getenv("CSRF_AUTH_KEY"),
		TokenGetter: func(c *gin.Context) string {
			return c.GetHeader("X-XSRF-TOKEN")
		},
		ErrorFunc: func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
		},
	}))
}
