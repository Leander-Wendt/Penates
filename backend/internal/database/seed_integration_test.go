//go:build integration

package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Leander-Wendt/Penates/backend/internal/database"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil"
)

func TestSeedDummyData_CreatesExpectedRecords(t *testing.T) {
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))

	require.NoError(t, database.SeedDummyData(db))

	var orgCount, userCount, itemCount, loanRequestCount, joinCount int64
	require.NoError(t, db.Model(&models.Organisation{}).Count(&orgCount).Error)
	require.NoError(t, db.Model(&models.User{}).Count(&userCount).Error)
	require.NoError(t, db.Model(&models.Item{}).Count(&itemCount).Error)
	require.NoError(t, db.Model(&models.LoanRequest{}).Count(&loanRequestCount).Error)
	require.NoError(t, db.Model(&models.LoanRequestItem{}).Count(&joinCount).Error)

	assert.Equal(t, int64(2), orgCount)
	assert.Equal(t, int64(3), userCount)
	assert.Equal(t, int64(6), itemCount)
	assert.Equal(t, int64(3), loanRequestCount)
	assert.Greater(t, joinCount, int64(0))

	var student models.User
	require.NoError(t, db.Where("email = ?", "student1@penates.local").First(&student).Error)
	assert.Equal(t, models.RoleStudent, student.Role)
}

func TestSeedDummyData_IsIdempotent(t *testing.T) {
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))

	require.NoError(t, database.SeedDummyData(db))
	require.NoError(t, database.SeedDummyData(db))

	var orgCount, userCount, itemCount, loanRequestCount int64
	require.NoError(t, db.Model(&models.Organisation{}).Count(&orgCount).Error)
	require.NoError(t, db.Model(&models.User{}).Count(&userCount).Error)
	require.NoError(t, db.Model(&models.Item{}).Count(&itemCount).Error)
	require.NoError(t, db.Model(&models.LoanRequest{}).Count(&loanRequestCount).Error)

	assert.Equal(t, int64(2), orgCount)
	assert.Equal(t, int64(3), userCount)
	assert.Equal(t, int64(6), itemCount)
	assert.Equal(t, int64(3), loanRequestCount)
}

func TestSeedDummyData_WorksAlongsideBootstrapAdmin(t *testing.T) {
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))

	require.NoError(t, database.BootstrapAdmin(db, "admin@penates.local", "admin-pw-123", "Default"))
	require.NoError(t, database.SeedDummyData(db))

	var userCount int64
	require.NoError(t, db.Model(&models.User{}).Count(&userCount).Error)
	assert.Equal(t, int64(4), userCount)

	var orgCount int64
	require.NoError(t, db.Model(&models.Organisation{}).Count(&orgCount).Error)
	assert.Equal(t, int64(3), orgCount)
}
