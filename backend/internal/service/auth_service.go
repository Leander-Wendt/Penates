package service

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	pkgauth "github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

// AuthService implements the login flow: verifying credentials and issuing JWTs.
type AuthService struct {
	repo      repository.UserRepository
	orgRepo   repository.OrganisationRepository
	jwtSecret string
	jwtExpiry time.Duration
	ssoConfig *SSOConfig
}

// SSOConfig holds SSO-specific configuration.
type SSOConfig struct {
	Enabled        bool
	DefaultRole    models.Role
	AllowedDomains []string
}

// JWTSecret returns the JWT secret for token generation.
func (s *AuthService) JWTSecret() string {
	return s.jwtSecret
}

// JWTExpiry returns the JWT expiry duration.
func (s *AuthService) JWTExpiry() time.Duration {
	return s.jwtExpiry
}

// NewAuthService builds an AuthService over the given UserRepository, JWT secret, and token expiry.
func NewAuthService(repo repository.UserRepository, orgRepo repository.OrganisationRepository, jwtSecret string, jwtExpiry time.Duration, ssoConfig *SSOConfig) *AuthService {
	return &AuthService{repo: repo, orgRepo: orgRepo, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry, ssoConfig: ssoConfig}
}

// Login verifies the given email/password pair and returns a signed JWT and its expiry.
// It returns ErrInvalidCredentials for both an unknown email and a wrong password, without
// distinguishing the two in the returned error.
func (s *AuthService) Login(email, password string) (string, time.Time, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", time.Time{}, ErrInvalidCredentials
		}
		return "", time.Time{}, err
	}
	if !pkgauth.CheckPasswordHash(password, user.PasswordHash) {
		return "", time.Time{}, ErrInvalidCredentials
	}
	token, expiresAt, err := pkgauth.GenerateToken(s.jwtSecret, user.ID, user.Role, s.jwtExpiry)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

// FindOrCreateSSOUser finds an existing user by SSO provider and subject, or creates a new one.
// If the user doesn't exist, it creates a new user with the email and name from SSO.
// The user's role is set to the SSO default role, and they're assigned to the default organisation.
func (s *AuthService) FindOrCreateSSOUser(provider, sub, email, name string) (*models.User, error) {
	// Try to find existing user by SSO
	user, err := s.repo.FindBySSOSub(provider, sub)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Try to find existing user by email
	user, err = s.repo.FindByEmail(email)
	if err == nil {
		// Link existing user to SSO provider
		user.SSOProvider = provider
		user.SSOSub = sub
		if err := s.repo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Validate email domain if configured
	if s.ssoConfig.Enabled && len(s.ssoConfig.AllowedDomains) > 0 {
		emailDomain := strings.Split(email, "@")
		if len(emailDomain) != 2 {
			return nil, errors.New("invalid email format")
		}
		allowed := false
		for _, domain := range s.ssoConfig.AllowedDomains {
			if emailDomain[1] == strings.TrimSpace(domain) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, errors.New("email domain not allowed")
		}
	}

	// Get default organisation
	orgs, err := s.orgRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if len(orgs) == 0 {
		return nil, errors.New("no organisations available")
	}

	// Create new user
	defaultRole := models.RoleStudent
	if s.ssoConfig != nil && s.ssoConfig.DefaultRole != "" {
		defaultRole = s.ssoConfig.DefaultRole
	}

	user = &models.User{
		Email:          email,
		PasswordHash:   "", // Empty for SSO users
		Role:           defaultRole,
		OrganisationID: orgs[0].ID,
		SSOProvider:    provider,
		SSOSub:         sub,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
