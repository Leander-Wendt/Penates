package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func TestAuthService_Login_Success(t *testing.T) {
	hash, err := auth.HashPassword("correct-password")
	require.NoError(t, err)

	repo := new(mocks.UserRepository)
	repo.On("FindByEmail", "user@example.com").Return(&models.User{ID: 1, Email: "user@example.com", PasswordHash: hash, Role: models.RoleStudent}, nil)
	svc := service.NewAuthService(repo, "test-secret", time.Hour)

	token, expiresAt, err := svc.Login("user@example.com", "correct-password")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().Add(time.Hour), expiresAt, 2*time.Second)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	hash, err := auth.HashPassword("correct-password")
	require.NoError(t, err)

	repo := new(mocks.UserRepository)
	repo.On("FindByEmail", "user@example.com").Return(&models.User{ID: 1, Email: "user@example.com", PasswordHash: hash, Role: models.RoleStudent}, nil)
	svc := service.NewAuthService(repo, "test-secret", time.Hour)

	_, _, err = svc.Login("user@example.com", "wrong-password")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestAuthService_Login_UnknownEmail(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("FindByEmail", "missing@example.com").Return(nil, gorm.ErrRecordNotFound)
	svc := service.NewAuthService(repo, "test-secret", time.Hour)

	_, _, err := svc.Login("missing@example.com", "whatever")
	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}
