package admin

import (
	"net/http"

	"spoony/config"

	"github.com/gin-gonic/gin"
)

// Auth gates admin routes behind the X-Admin-Token header.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.App.AdminToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin panel is disabled (ADMIN_TOKEN not set)"})
			c.Abort()
			return
		}
		token := c.GetHeader("X-Admin-Token")
		if token != config.App.AdminToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid admin token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
