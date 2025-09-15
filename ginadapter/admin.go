package ginadapter

import (
	"net/http"

	core "github.com/DrWeltschmerz/users-core"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware checks if the authenticated user is an admin.
func AdminMiddleware(svc *core.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("userID")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
			return
		}
		user, err := svc.GetUserByID(c.Request.Context(), userID.(string))
		if err != nil || !svc.IsAdmin(user) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}
