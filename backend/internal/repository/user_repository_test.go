//go:build integration

package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/database"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil"
)

func setupUserTestDB(t *testing.T) (*gorm.DB, uint) {
	t.Helper()
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))
	org := models.Organisation{Name: "Test Org"}
	require.NoError(t, db.Create(&org).Error)
	return db, org.ID
}

func TestGormUserRepository_CreateAndFindByID(t *testing.T) {
	db, orgID := setupUserTestDB(t)
	repo := repository.NewGormUserRepository(db)

	user := &models.User{Email: "a@example.com", PasswordHash: "hash", Role: models.RoleStudent, OrganisationID: orgID}
	require.NoError(t, repo.Create(user))
	assert.NotZero(t, user.ID)

	found, err := repo.FindByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, "a@example.com", found.Email)
	assert.Equal(t, "Test Org", found.Organisation.Name)
}

func TestGormUserRepository_FindByEmail(t *testing.T) {
	db, orgID := setupUserTestDB(t)
	repo := repository.NewGormUserRepository(db)

	user := &models.User{Email: "b@example.com", PasswordHash: "hash", Role: models.RoleAdmin, OrganisationID: orgID}
	require.NoError(t, repo.Create(user))

	found, err := repo.FindByEmail("b@example.com")
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestGormUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, _ := setupUserTestDB(t)
	repo := repository.NewGormUserRepository(db)

	_, err := repo.FindByEmail("missing@example.com")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestGormUserRepository_Count(t *testing.T) {
	db, orgID := setupUserTestDB(t)
	repo := repository.NewGormUserRepository(db)

	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	require.NoError(t, repo.Create(&models.User{Email: "c@example.com", PasswordHash: "hash", Role: models.RoleStudent, OrganisationID: orgID}))

	count, err = repo.Count()
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestGormUserRepository_Delete(t *testing.T) {
	db, orgID := setupUserTestDB(t)
	repo := repository.NewGormUserRepository(db)

	user := &models.User{Email: "d@example.com", PasswordHash: "hash", Role: models.RoleStudent, OrganisationID: orgID}
	require.NoError(t, repo.Create(user))

	require.NoError(t, repo.Delete(user.ID))

	_, err := repo.FindByID(user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
