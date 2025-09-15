package ginadapter

import (
	core "github.com/DrWeltschmerz/users-core"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all user and role routes on the given Gin router.
// Requires a users-core Service and Tokenizer (for JWT auth).
func RegisterRoutes(r *gin.Engine, svc *core.Service, tokenizer core.Tokenizer) {
	h := &UserHandlers{Svc: svc, Tokenizer: tokenizer}

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// Authenticated user routes
	auth := r.Group("/user")
	auth.Use(JWTMiddleware(tokenizer))
	auth.GET("/profile", h.GetProfile)
	auth.PUT("/profile", h.UpdateProfile)
	auth.POST("/change-password", h.ChangePassword)

	// Admin routes
	admin := r.Group("")
	admin.Use(JWTMiddleware(tokenizer), AdminMiddleware(svc))
	admin.GET("/users", h.ListUsers)
	admin.DELETE("/users/:id", h.DeleteUser)
	admin.GET("/roles", h.ListRoles)
	admin.POST("/users/:id/assign-role", h.AssignRole)
	admin.POST("/users/:id/reset-password", h.ResetPassword)
}
