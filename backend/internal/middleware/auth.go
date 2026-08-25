package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
)

// bearerChallenge is the RFC 6750 §3 WWW-Authenticate header value sent on any
// 401 Unauthorized response caused by a missing, malformed, or invalid token.
const bearerChallenge = `Bearer realm="Penates API", error="invalid_token"`

// RequireAuth returns a gin.HandlerFunc that validates the Authorization
// Bearer token on every request using secret. If the header is missing,
// malformed, or the token is invalid/expired, it responds 401 Unauthorized
// with a WWW-Authenticate challenge and a generic error body. Otherwise it
// sets "userID" and "role" in the Gin context.
func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			unauthorized(c)
			return
		}

		claims, err := auth.ParseToken(secret, parts[1])
		if err != nil {
			unauthorized(c)
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", bearerChallenge)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}
