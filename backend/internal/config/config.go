package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the Penates backend, loaded from environment variables.
type Config struct {
	Env                    string
	Port                   string
	DatabaseURL            string
	JWTSecret              string
	JWTExpiry              time.Duration
	CORSAllowedOrigins     []string
	SeedAdminEmail         string
	SeedAdminPassword      string
	SeedAdminOrg           string
	UploadDir              string
	MaxUploadSizeMB        int64
	LoginRateLimitAttempts int
	LoginRateLimitWindow   time.Duration
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// Load reads configuration from environment variables, loading a .env file first when ENV is not "production".
func Load() (*Config, error) {
	env := getEnv("ENV", "development")
	if env != "production" {
		_ = godotenv.Load()
	}

	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		return nil, err
	}

	maxUploadSizeMB, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)
	if err != nil {
		return nil, err
	}

	loginRateLimitAttempts, err := strconv.Atoi(getEnv("LOGIN_RATE_LIMIT_ATTEMPTS", "5"))
	if err != nil {
		return nil, err
	}

	loginRateLimitWindow, err := time.ParseDuration(getEnv("LOGIN_RATE_LIMIT_WINDOW", "1m"))
	if err != nil {
		return nil, err
	}

	var origins []string
	if raw := getEnv("CORS_ALLOWED_ORIGINS", ""); raw != "" {
		for _, o := range strings.Split(raw, ",") {
			origins = append(origins, strings.TrimSpace(o))
		}
	}

	cfg := &Config{
		Env:                    env,
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/penates?sslmode=disable"),
		JWTSecret:              getEnv("JWT_SECRET", "change-me"),
		JWTExpiry:              jwtExpiry,
		CORSAllowedOrigins:     origins,
		SeedAdminEmail:         getEnv("SEED_ADMIN_EMAIL", "admin@penates.local"),
		SeedAdminPassword:      getEnv("SEED_ADMIN_PASSWORD", "change-me"),
		SeedAdminOrg:           getEnv("SEED_ADMIN_ORG", "Default"),
		UploadDir:              getEnv("UPLOAD_DIR", "uploads"),
		MaxUploadSizeMB:        maxUploadSizeMB,
		LoginRateLimitAttempts: loginRateLimitAttempts,
		LoginRateLimitWindow:   loginRateLimitWindow,
	}
	return cfg, nil
}
