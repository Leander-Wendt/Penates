package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgauth "github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/middleware"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

func newAuthTestRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", middleware.RequireAuth(secret), func(c *gin.Context) {
		userID, _ := c.Get("userID")
		role, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{"userID": userID, "role": role})
	})
	return r
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	r := newAuthTestRouter("secret")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `Bearer realm="Penates API", error="invalid_token"`, w.Header().Get("WWW-Authenticate"))
}

func TestRequireAuth_MalformedHeader(t *testing.T) {
	r := newAuthTestRouter("secret")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "NotBearer sometoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `Bearer realm="Penates API", error="invalid_token"`, w.Header().Get("WWW-Authenticate"))
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	r := newAuthTestRouter("secret")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `Bearer realm="Penates API", error="invalid_token"`, w.Header().Get("WWW-Authenticate"))
}

func TestRequireAuth_ExpiredToken(t *testing.T) {
	r := newAuthTestRouter("secret")
	token, _, err := pkgauth.GenerateToken("secret", 1, models.RoleStudent, -time.Hour)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `Bearer realm="Penates API", error="invalid_token"`, w.Header().Get("WWW-Authenticate"))
}

func TestRequireAuth_ValidToken(t *testing.T) {
	r := newAuthTestRouter("secret")
	token, _, err := pkgauth.GenerateToken("secret", 42, models.RoleAdmin, time.Hour)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("WWW-Authenticate"))
}
