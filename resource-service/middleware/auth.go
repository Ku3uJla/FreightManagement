package middleware

import (
	"user-service/features"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("JWT")
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}
		userID, role, err := features.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
		}
		ctx.Set("userID", userID)
		ctx.Set("userRole", role)
		ctx.Next()
	}
}
