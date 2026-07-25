package middleware

import (
	"resource-service/features"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("JWT")
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}
		userID, role, err := features.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
		}
		c.Set("userID", userID)
		c.Set("userRole", role)
		c.Next()
	}
}
