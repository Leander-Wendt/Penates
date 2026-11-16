package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	pkgauth "github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
)

// OAuthConfig holds the OIDC configuration for a provider.
type OAuthConfig struct {
	ProviderURL         string
	ClientID            string
	ClientSecret        string
	RedirectURL         string // Backend callback URL for OAuth flow
	FrontendCallbackURL string // Frontend URL to redirect to after successful auth
	Scopes              []string
	AllowedDomains      []string
	DefaultRole         string
}

// OAuthHandler handles OIDC authentication flows.
type OAuthHandler struct {
	svc         *service.AuthService
	jwtSecret   string
	jwtExpiry   time.Duration
	oauthConfig OAuthConfig
	// stateStore stores random state strings for CSRF protection
	// In production, use a proper session store or cache
	stateStore map[string]bool
}

// NewOAuthHandler creates a new OAuthHandler.
func NewOAuthHandler(svc *service.AuthService, jwtSecret string, jwtExpiry time.Duration, cfg OAuthConfig) *OAuthHandler {
	return &OAuthHandler{
		svc:         svc,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
		oauthConfig: cfg,
		stateStore:  make(map[string]bool),
	}
}

// generateState generates a random state string for CSRF protection.
func (h *OAuthHandler) generateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// isAllowedDomain checks if the email domain is in the allowed list.
func (h *OAuthHandler) isAllowedDomain(email string) bool {
	if len(h.oauthConfig.AllowedDomains) == 0 {
		return true
	}

	emailDomain := strings.Split(email, "@")
	if len(emailDomain) != 2 {
		return false
	}

	for _, domain := range h.oauthConfig.AllowedDomains {
		if emailDomain[1] == strings.TrimSpace(domain) {
			return true
		}
	}
	return false
}

// Redirect redirects the user to the OIDC provider for authentication.
// GET /api/v1/auth/oauth/redirect
func (h *OAuthHandler) Redirect(c *gin.Context) {
	// Check if SSO is enabled
	if h.oauthConfig.ClientID == "" || h.oauthConfig.ClientSecret == "" || h.oauthConfig.ProviderURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SSO is not configured"})
		return
	}

	provider, err := oidc.NewProvider(context.Background(), h.oauthConfig.ProviderURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to configure OIDC provider"})
		return
	}

	conf := oauth2.Config{
		RedirectURL:  h.oauthConfig.RedirectURL,
		ClientID:     h.oauthConfig.ClientID,
		ClientSecret: h.oauthConfig.ClientSecret,
		Scopes:       h.oauthConfig.Scopes,
		Endpoint:     provider.Endpoint(),
	}

	// Generate and store state for CSRF protection
	state := h.generateState()
	h.stateStore[state] = true

	// Redirect to provider
	authURL := conf.AuthCodeURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// oidcUser represents the user info from OIDC provider.
type oidcUser struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Sub     string `json:"sub"`
	Picture string `json:"picture,omitempty"`
}

// Callback handles the OIDC callback from the provider.
// GET /api/v1/auth/oauth/callback
func (h *OAuthHandler) Callback(c *gin.Context) {
	// Check if SSO is enabled
	if h.oauthConfig.ClientID == "" || h.oauthConfig.ClientSecret == "" || h.oauthConfig.ProviderURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SSO is not configured"})
		return
	}

	// Verify state for CSRF protection
	state := c.Query("state")
	if state == "" || !h.stateStore[state] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}
	delete(h.stateStore, state)

	// Handle errors from provider
	if errMsg := c.Query("error"); errMsg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg, "error_description": c.Query("error_description")})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	provider, err := oidc.NewProvider(context.Background(), h.oauthConfig.ProviderURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to configure OIDC provider"})
		return
	}

	conf := oauth2.Config{
		RedirectURL:  h.oauthConfig.RedirectURL,
		ClientID:     h.oauthConfig.ClientID,
		ClientSecret: h.oauthConfig.ClientSecret,
		Scopes:       h.oauthConfig.Scopes,
		Endpoint:     provider.Endpoint(),
	}

	// Exchange code for token
	oauth2Token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	// Extract the ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No id_token field in oauth2 token"})
		return
	}

	// Verify the ID token
	idToken, err := provider.Verifier(&oidc.Config{ClientID: h.oauthConfig.ClientID}).Verify(context.Background(), rawIDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ID token: " + err.Error()})
		return
	}

	// Extract claims
	var claims oidcUser
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
		return
	}

	// Validate email
	if claims.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not provided by identity provider"})
		return
	}

	// Validate domain
	if !h.isAllowedDomain(claims.Email) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email domain not allowed"})
		return
	}

	// Find or create user
	providerName := "authentik" // Or extract from provider URL
	user, err := h.svc.FindOrCreateSSOUser(providerName, claims.Sub, claims.Email, claims.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Generate JWT for Penates
	token, expiresAt, err := pkgauth.GenerateToken(h.jwtSecret, user.ID, user.Role, h.jwtExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Redirect back to frontend with token in URL hash
	// The frontend will extract the token from the hash and complete the login
	frontendCallbackURL := h.oauthConfig.FrontendCallbackURL
	if frontendCallbackURL == "" {
		// Fallback: derive from RedirectURL by replacing the path
		frontendCallbackURL = strings.Replace(h.oauthConfig.RedirectURL, "/api/v1/auth/oauth/callback", "/login", 1)
	}
	separator := "#"
	if strings.Contains(frontendCallbackURL, "#") {
		separator = "&"
	}
	c.Redirect(http.StatusFound, frontendCallbackURL+separator+"token="+token+"&expiresAt="+expiresAt.Format(time.RFC3339))
}

// GetConfig returns the OIDC configuration for the frontend.
// GET /api/v1/auth/oauth/config
func (h *OAuthHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"enabled":        h.oauthConfig.ClientID != "" && h.oauthConfig.ClientSecret != "" && h.oauthConfig.ProviderURL != "",
		"providerURL":    h.oauthConfig.ProviderURL,
		"clientID":       h.oauthConfig.ClientID,
		"redirectURL":    h.oauthConfig.RedirectURL,
		"allowedDomains": h.oauthConfig.AllowedDomains,
	})
}
