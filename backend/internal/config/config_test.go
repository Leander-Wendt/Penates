package config_test

import (
	"testing"
	"time"

	"github.com/Leander-Wendt/Penates/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		"ENV", "PORT", "DATABASE_URL", "JWT_SECRET", "JWT_EXPIRY",
		"CORS_ALLOWED_ORIGINS", "SEED_ADMIN_EMAIL", "SEED_ADMIN_PASSWORD",
		"SEED_ADMIN_ORG", "UPLOAD_DIR", "MAX_UPLOAD_SIZE_MB",
		"LOGIN_RATE_LIMIT_ATTEMPTS", "LOGIN_RATE_LIMIT_WINDOW",
	}
	for _, k := range keys {
		t.Setenv(k, "")
	}
}

func TestLoad_Defaults(t *testing.T) {
	clearEnv(t)
	t.Setenv("ENV", "test")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "uploads", cfg.UploadDir)
	assert.Equal(t, int64(10), cfg.MaxUploadSizeMB)
	assert.Equal(t, time.Hour*24, cfg.JWTExpiry)
	assert.Equal(t, 5, cfg.LoginRateLimitAttempts)
	assert.Equal(t, time.Minute, cfg.LoginRateLimitWindow)
}

func TestLoad_OverridesFromEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv("ENV", "test")
	t.Setenv("PORT", "9090")
	t.Setenv("JWT_EXPIRY", "2h")
	t.Setenv("MAX_UPLOAD_SIZE_MB", "25")
	t.Setenv("LOGIN_RATE_LIMIT_ATTEMPTS", "10")
	t.Setenv("LOGIN_RATE_LIMIT_WINDOW", "30s")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, 2*time.Hour, cfg.JWTExpiry)
	assert.Equal(t, int64(25), cfg.MaxUploadSizeMB)
	assert.Equal(t, 10, cfg.LoginRateLimitAttempts)
	assert.Equal(t, 30*time.Second, cfg.LoginRateLimitWindow)
}

func TestLoad_CORSOriginsParsed(t *testing.T) {
	clearEnv(t)
	t.Setenv("ENV", "test")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://a.com,http://b.com")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, []string{"http://a.com", "http://b.com"}, cfg.CORSAllowedOrigins)
}
