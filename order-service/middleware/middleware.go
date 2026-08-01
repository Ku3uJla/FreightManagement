package middleware

import "github.com/gin-gonic/gin"

type Headers struct {
	Authorization string `header:"Authorization"`
	UserID        string `header:"X-User-ID"`
	Role          string `header:"X-User-Role"`
}

func Middleawre() gin.HandlerFunc {
	return func(c *gin.Context) {
		var headers Headers
		if err := c.ShouldBindHeader(&headers); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid headers"})
			return
		}
		if headers.UserID == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing UserID"})
			return
		}
		c.Set("userID", headers.UserID)
		c.Set("userRole", headers.Role)
		c.Next()
	}
}
