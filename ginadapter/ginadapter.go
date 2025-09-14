// Package ginadapter provides a generic Gin HTTP API for users-core.package ginadapter

// Usage: import this package, inject your users.Service and Tokenizer, and call RegisterRoutes.
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

	// Admin/user management
	r.GET("/users", h.ListUsers)
	r.DELETE("/users/:id", h.DeleteUser)

	// Role management
	r.GET("/roles", h.ListRoles)
	r.POST("/users/:id/assign-role", h.AssignRole)
	r.POST("/users/:id/reset-password", h.ResetPassword)
}
