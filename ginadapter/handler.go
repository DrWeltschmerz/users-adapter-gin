package ginadapter

import (
	"net/http"

	core "github.com/DrWeltschmerz/users-core"
	"github.com/gin-gonic/gin"
)

// UserHandlers holds dependencies for user endpoints.
type UserHandlers struct {
	Svc       *core.Service
	Tokenizer core.Tokenizer
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body core.UserRegisterInput true "User registration input"
// @Success 201 {object} core.User
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (h *UserHandlers) Register(c *gin.Context) {
	var input core.UserRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	user, err := h.Svc.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body core.UserLoginInput true "User login input"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *UserHandlers) Login(c *gin.Context) {
	var input core.UserLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	token, err := h.Svc.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} core.User
// @Failure 401 {object} map[string]string
// @Router /user/profile [get]
func (h *UserHandlers) GetProfile(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	user, err := h.Svc.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Update the profile of the authenticated user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body core.User true "User update input"
// @Success 200 {object} core.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/profile [put]
func (h *UserHandlers) UpdateProfile(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	var input core.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	input.ID = userID.(string)
	user, err := h.Svc.UpdateUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ListUsers godoc
// @Summary List all users
// @Description List all users
// @Tags admin
// @Produce json
// @Success 200 {array} core.User
// @Router /users [get]
func (h *UserHandlers) ListUsers(c *gin.Context) {
	users, err := h.Svc.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags admin
// @Param id path string true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandlers) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	err := h.Svc.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ListRoles godoc
// @Summary List all roles
// @Description List all roles
// @Tags admin
// @Produce json
// @Success 200 {array} core.Role
// @Router /roles [get]
func (h *UserHandlers) ListRoles(c *gin.Context) {
	roles, err := h.Svc.ListRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// AssignRole godoc
// @Summary Assign role to user
// @Description Assign a role to a user
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param input body struct{RoleID string `json:"role_id"`} true "Role assignment input"
// @Success 200 {object} core.User
// @Failure 400 {object} map[string]string
// @Router /users/{id}/assign-role [post]
func (h *UserHandlers) AssignRole(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		RoleID string `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RoleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	user, err := h.Svc.AssignRoleToUser(c.Request.Context(), userID, req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password for the authenticated user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body struct{OldPassword string `json:"old_password"`; NewPassword string `json:"new_password"`} true "Password change input"
// @Success 200 {object} core.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /user/change-password [post]
func (h *UserHandlers) ChangePassword(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	user, err := h.Svc.ChangePassword(c.Request.Context(), userID.(string), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset password for a user by ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param input body struct{NewPassword string `json:"new_password"`} true "Password reset input"
// @Success 200 {object} core.User
// @Failure 400 {object} map[string]string
// @Router /users/{id}/reset-password [post]
func (h *UserHandlers) ResetPassword(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	user, err := h.Svc.ResetPassword(c.Request.Context(), userID, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
