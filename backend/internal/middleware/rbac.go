package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// RequireRoles returns a gin.HandlerFunc that allows the request to proceed
// only if the "role" set in the Gin context by RequireAuth is one of roles.
// Otherwise it responds 403 Forbidden with no WWW-Authenticate header, since
// this rejects valid credentials with insufficient permission rather than
// missing/invalid credentials.
func RequireRoles(roles ...models.Role) gin.HandlerFunc {
	allowed := make(map[models.Role]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		userRole, ok := role.(models.Role)
		if !ok || !allowed[userRole] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
