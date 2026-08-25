package auth_test

import (
	"testing"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := auth.HashPassword("s3cr3t-password")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "s3cr3t-password", hash)
}

func TestCheckPasswordHash_Success(t *testing.T) {
	hash, err := auth.HashPassword("correct-horse")
	require.NoError(t, err)
	assert.True(t, auth.CheckPasswordHash("correct-horse", hash))
}

func TestCheckPasswordHash_Failure(t *testing.T) {
	hash, err := auth.HashPassword("correct-horse")
	require.NoError(t, err)
	assert.False(t, auth.CheckPasswordHash("wrong-password", hash))
}
