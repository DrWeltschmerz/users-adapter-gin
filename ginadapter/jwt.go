package ginadapter

import (
	"net/http"
	"strings"

	core "github.com/DrWeltschmerz/users-core"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware returns a Gin middleware that validates JWTs using the provided Tokenizer.
func JWTMiddleware(tokenizer core.Tokenizer) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		userID, err := tokenizer.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
