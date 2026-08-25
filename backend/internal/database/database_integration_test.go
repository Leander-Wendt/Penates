//go:build integration

package database_test

import (
	"testing"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/database"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAutoMigrate_CreatesAllTables(t *testing.T) {
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))

	for _, table := range []string{"organisations", "users", "items", "loan_requests", "loan_request_items"} {
		assert.True(t, db.Migrator().HasTable(table), "expected table %q to exist", table)
	}
}

func TestBootstrapAdmin_CreatesAdminWhenUsersEmpty(t *testing.T) {
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))

	require.NoError(t, database.BootstrapAdmin(db, "admin@example.com", "s3cr3t-password", "Default"))

	var user models.User
	require.NoError(t, db.Where("email = ?", "admin@example.com").First(&user).Error)
	assert.Equal(t, models.RoleAdmin, user.Role)
	assert.True(t, auth.CheckPasswordHash("s3cr3t-password", user.PasswordHash))

	var org models.Organisation
	require.NoError(t, db.First(&org, user.OrganisationID).Error)
	assert.Equal(t, "Default", org.Name)
}

func TestBootstrapAdmin_NoOpWhenUsersExist(t *testing.T) {
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))

	require.NoError(t, database.BootstrapAdmin(db, "first@example.com", "password-one", "Default"))
	require.NoError(t, database.BootstrapAdmin(db, "second@example.com", "password-two", "Default"))

	var count int64
	require.NoError(t, db.Model(&models.User{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)

	var user models.User
	require.NoError(t, db.First(&user).Error)
	assert.Equal(t, "first@example.com", user.Email)
}

func TestBootstrapAdmin_ReusesExistingOrganisation(t *testing.T) {
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))

	org := models.Organisation{Name: "Default"}
	require.NoError(t, db.Create(&org).Error)

	require.NoError(t, database.BootstrapAdmin(db, "admin@example.com", "s3cr3t-password", "Default"))

	var count int64
	require.NoError(t, db.Model(&models.Organisation{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}
