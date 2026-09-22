package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// staleWindowSweepInterval is how often LoginRateLimiter clears out expired
// per-IP windows, so long-running processes don't accumulate one entry per
// distinct IP forever.
const staleWindowSweepInterval = 5 * time.Minute

// loginWindow tracks how many login attempts a single client IP has made
// within the current fixed window.
type loginWindow struct {
	count     int
	expiresAt time.Time
}

// LoginRateLimiter throttles login attempts per client IP using a fixed
// window: at most maxAttempts requests within window, after which further
// requests are rejected until the window resets.
type LoginRateLimiter struct {
	mu          sync.Mutex
	windows     map[string]*loginWindow
	maxAttempts int
	window      time.Duration
	lastSweptAt time.Time
}

// NewLoginRateLimiter builds a LoginRateLimiter allowing at most maxAttempts
// requests per client IP within window.
func NewLoginRateLimiter(maxAttempts int, window time.Duration) *LoginRateLimiter {
	return &LoginRateLimiter{
		windows:     make(map[string]*loginWindow),
		maxAttempts: maxAttempts,
		window:      window,
		lastSweptAt: time.Now(),
	}
}

// Middleware returns a gin.HandlerFunc that rejects requests once the
// calling IP (per gin.Context.ClientIP) has exceeded the configured attempt
// limit within the current window, responding 429 Too Many Requests with a
// Retry-After header giving the number of seconds until the window resets.
func (l *LoginRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		retryAfter, allowed := l.allow(c.ClientIP(), time.Now())
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many login attempts, please try again later"})
			return
		}
		c.Next()
	}
}

// allow records an attempt for key at time now, reporting whether it's
// within the limit. When it isn't, it also returns how long the caller must
// wait before the window resets.
func (l *LoginRateLimiter) allow(key string, now time.Time) (time.Duration, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.sweepExpired(now)

	w, ok := l.windows[key]
	if !ok || now.After(w.expiresAt) {
		l.windows[key] = &loginWindow{count: 1, expiresAt: now.Add(l.window)}
		return 0, true
	}
	if w.count >= l.maxAttempts {
		return w.expiresAt.Sub(now), false
	}
	w.count++
	return 0, true
}

// sweepExpired removes windows that expired before now, if enough time has
// passed since the last sweep to make another one worthwhile. Callers must
// hold l.mu.
func (l *LoginRateLimiter) sweepExpired(now time.Time) {
	if now.Sub(l.lastSweptAt) < staleWindowSweepInterval {
		return
	}
	for key, w := range l.windows {
		if now.After(w.expiresAt) {
			delete(l.windows, key)
		}
	}
	l.lastSweptAt = now
}
