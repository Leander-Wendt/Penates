package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	pkgauth "github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

// AuthService implements the login flow: verifying credentials and issuing JWTs.
type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewAuthService builds an AuthService over the given UserRepository, JWT secret, and token expiry.
func NewAuthService(repo repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
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
