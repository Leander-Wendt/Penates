package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Leander-Wendt/Penates/backend/internal/middleware"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

func newRBACTestRouter(role models.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/admin-only", func(c *gin.Context) {
		c.Set("role", role)
		c.Next()
	}, middleware.RequireRoles(models.RoleAdmin), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func TestRequireRoles_AllowedRole(t *testing.T) {
	r := newRBACTestRouter(models.RoleAdmin)
	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRoles_DisallowedRole(t *testing.T) {
	r := newRBACTestRouter(models.RoleStudent)
	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("WWW-Authenticate"))
}
