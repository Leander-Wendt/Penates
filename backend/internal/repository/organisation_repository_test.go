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

func setupOrganisationTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))
	return db
}

func TestGormOrganisationRepository_CreateAndFindByID(t *testing.T) {
	db := setupOrganisationTestDB(t)
	repo := repository.NewGormOrganisationRepository(db)

	org := &models.Organisation{Name: "Acme", Description: "Acme Corp"}
	require.NoError(t, repo.Create(org))
	assert.NotZero(t, org.ID)

	found, err := repo.FindByID(org.ID)
	require.NoError(t, err)
	assert.Equal(t, "Acme", found.Name)
}

func TestGormOrganisationRepository_FindAll(t *testing.T) {
	db := setupOrganisationTestDB(t)
	repo := repository.NewGormOrganisationRepository(db)

	require.NoError(t, repo.Create(&models.Organisation{Name: "Org A"}))
	require.NoError(t, repo.Create(&models.Organisation{Name: "Org B"}))

	all, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestGormOrganisationRepository_Update(t *testing.T) {
	db := setupOrganisationTestDB(t)
	repo := repository.NewGormOrganisationRepository(db)

	org := &models.Organisation{Name: "Original"}
	require.NoError(t, repo.Create(org))

	org.Name = "Updated"
	require.NoError(t, repo.Update(org))

	found, err := repo.FindByID(org.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", found.Name)
}

func TestGormOrganisationRepository_Delete(t *testing.T) {
	db := setupOrganisationTestDB(t)
	repo := repository.NewGormOrganisationRepository(db)

	org := &models.Organisation{Name: "ToDelete"}
	require.NoError(t, repo.Create(org))

	require.NoError(t, repo.Delete(org.ID))

	_, err := repo.FindByID(org.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
