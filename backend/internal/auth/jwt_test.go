package auth_test

import (
	"testing"
	"time"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	token, expiresAt, err := auth.GenerateToken(secret, 1, models.RoleAdmin, time.Hour)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().Add(time.Hour), expiresAt, 2*time.Second)

	claims, err := auth.ParseToken(secret, token)
	require.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, models.RoleAdmin, claims.Role)
}

func TestParseToken_InvalidSignature(t *testing.T) {
	token, _, err := auth.GenerateToken("secret-a", 1, models.RoleStudent, time.Hour)
	require.NoError(t, err)

	_, err = auth.ParseToken("secret-b", token)
	assert.Error(t, err)
}

func TestParseToken_Expired(t *testing.T) {
	token, _, err := auth.GenerateToken("test-secret", 1, models.RoleStudent, -time.Hour)
	require.NoError(t, err)

	_, err = auth.ParseToken("test-secret", token)
	assert.Error(t, err)
}

func TestParseToken_Malformed(t *testing.T) {
	_, err := auth.ParseToken("test-secret", "not-a-valid-token")
	assert.Error(t, err)
}
