package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Leander-Wendt/Penates/backend/internal/middleware"
)

func newRateLimitTestRouter(limiter *middleware.LoginRateLimiter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", limiter.Middleware(), func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func doLoginRequest(r *gin.Engine, ip string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = ip + ":12345"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestLoginRateLimiter_AllowsUpToLimit(t *testing.T) {
	limiter := middleware.NewLoginRateLimiter(3, time.Minute)
	r := newRateLimitTestRouter(limiter)

	for i := 0; i < 3; i++ {
		w := doLoginRequest(r, "1.2.3.4")
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestLoginRateLimiter_BlocksOverLimit(t *testing.T) {
	limiter := middleware.NewLoginRateLimiter(3, time.Minute)
	r := newRateLimitTestRouter(limiter)

	for i := 0; i < 3; i++ {
		doLoginRequest(r, "1.2.3.4")
	}
	w := doLoginRequest(r, "1.2.3.4")

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.NotEmpty(t, w.Header().Get("Retry-After"))
}

func TestLoginRateLimiter_TracksIPsIndependently(t *testing.T) {
	limiter := middleware.NewLoginRateLimiter(1, time.Minute)
	r := newRateLimitTestRouter(limiter)

	w1 := doLoginRequest(r, "1.2.3.4")
	w2 := doLoginRequest(r, "5.6.7.8")

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestLoginRateLimiter_ResetsAfterWindow(t *testing.T) {
	limiter := middleware.NewLoginRateLimiter(1, 30*time.Millisecond)
	r := newRateLimitTestRouter(limiter)

	doLoginRequest(r, "1.2.3.4")
	blocked := doLoginRequest(r, "1.2.3.4")
	assert.Equal(t, http.StatusTooManyRequests, blocked.Code)

	time.Sleep(50 * time.Millisecond)

	allowed := doLoginRequest(r, "1.2.3.4")
	assert.Equal(t, http.StatusOK, allowed.Code)
}
