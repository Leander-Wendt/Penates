package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// currentUserID returns the authenticated user's ID from the Gin context, as
// set by middleware.RequireAuth.
func currentUserID(c *gin.Context) uint {
	v, _ := c.Get("userID")
	id, _ := v.(uint)
	return id
}

// currentRole returns the authenticated user's Role from the Gin context, as
// set by middleware.RequireAuth.
func currentRole(c *gin.Context) models.Role {
	v, _ := c.Get("role")
	role, _ := v.(models.Role)
	return role
}
