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
	// TLSCertFile and TLSKeyFile enable HTTPS when both are set; the server
	// falls back to plain HTTP when either is empty (e.g. in tests).
	TLSCertFile string
	TLSKeyFile  string
	// SSO configuration - all OIDC settings are optional to allow disabling SSO
	SSOEnabled             bool
	SSOProviderURL         string   // OIDC provider discovery URL (e.g., https://authentik.example.com/application/o/penates/.well-known/openid-configuration)
	SSOClientID            string   // OIDC client ID
	SSOClientSecret        string   // OIDC client secret
	SSORedirectURL         string   // Backend callback URL (e.g., https://penates.example.com/api/v1/auth/oauth/callback)
	SSOFrontendCallbackURL string   // Frontend URL to redirect to after auth (e.g., https://penates.example.com/login)
	SSOScopes              []string // Requested scopes (e.g., openid,email,profile)
	SSOAllowedDomains      []string // Allowed email domains for SSO users
	SSODefaultRole         string   // Default role for new SSO users (admin, logistics, student)
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

	// SSO configuration
	ssoEnabled, err := strconv.ParseBool(getEnv("SSO_ENABLED", "false"))
	if err != nil {
		return nil, err
	}

	var ssoScopes []string
	if raw := getEnv("SSO_SCOPES", "openid,email,profile"); raw != "" {
		ssoScopes = strings.Split(raw, ",")
		for i, scope := range ssoScopes {
			ssoScopes[i] = strings.TrimSpace(scope)
		}
	}

	var ssoAllowedDomains []string
	if raw := getEnv("SSO_ALLOWED_DOMAINS", ""); raw != "" {
		for _, d := range strings.Split(raw, ",") {
			ssoAllowedDomains = append(ssoAllowedDomains, strings.TrimSpace(d))
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
		TLSCertFile:            getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:             getEnv("TLS_KEY_FILE", ""),
		SSOEnabled:             ssoEnabled,
		SSOProviderURL:         getEnv("SSO_PROVIDER_URL", ""),
		SSOClientID:            getEnv("SSO_CLIENT_ID", ""),
		SSOClientSecret:        getEnv("SSO_CLIENT_SECRET", ""),
		SSORedirectURL:         getEnv("SSO_REDIRECT_URL", ""),
		SSOFrontendCallbackURL: getEnv("SSO_FRONTEND_CALLBACK_URL", ""),
		SSOScopes:              ssoScopes,
		SSOAllowedDomains:      ssoAllowedDomains,
		SSODefaultRole:         getEnv("SSO_DEFAULT_ROLE", "student"),
	}
	return cfg, nil
}
