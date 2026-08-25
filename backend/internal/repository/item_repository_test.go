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

func setupItemTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))
	return db
}

func TestGormItemRepository_CreateAndFindByInventoryNumber(t *testing.T) {
	db := setupItemTestDB(t)
	repo := repository.NewGormItemRepository(db)

	item := &models.Item{InventoryNumber: "100001", Name: "Hammer", Category: "Tools", Amount: 5}
	require.NoError(t, repo.Create(item))

	found, err := repo.FindByInventoryNumber("100001")
	require.NoError(t, err)
	assert.Equal(t, "Hammer", found.Name)
}

func TestGormItemRepository_FindAll_FilterByCategory(t *testing.T) {
	db := setupItemTestDB(t)
	repo := repository.NewGormItemRepository(db)

	require.NoError(t, repo.Create(&models.Item{InventoryNumber: "100002", Name: "Drill", Category: "Tools", Amount: 1}))
	require.NoError(t, repo.Create(&models.Item{InventoryNumber: "100003", Name: "Laptop", Category: "Electronics", Amount: 1}))

	items, err := repo.FindAll(repository.ItemFilter{Category: "Tools"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "Drill", items[0].Name)
}

func TestGormItemRepository_FindAll_Search(t *testing.T) {
	db := setupItemTestDB(t)
	repo := repository.NewGormItemRepository(db)

	require.NoError(t, repo.Create(&models.Item{InventoryNumber: "100004", Name: "Cordless Drill", Category: "Tools", Amount: 1}))
	require.NoError(t, repo.Create(&models.Item{InventoryNumber: "100005", Name: "Laptop", Category: "Electronics", Amount: 1}))

	items, err := repo.FindAll(repository.ItemFilter{Search: "drill"})
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "100004", items[0].InventoryNumber)
}

func TestGormItemRepository_Update(t *testing.T) {
	db := setupItemTestDB(t)
	repo := repository.NewGormItemRepository(db)

	item := &models.Item{InventoryNumber: "100006", Name: "Original", Category: "Tools", Amount: 1}
	require.NoError(t, repo.Create(item))

	item.Name = "Updated"
	require.NoError(t, repo.Update(item))

	found, err := repo.FindByInventoryNumber("100006")
	require.NoError(t, err)
	assert.Equal(t, "Updated", found.Name)
}

func TestGormItemRepository_Delete(t *testing.T) {
	db := setupItemTestDB(t)
	repo := repository.NewGormItemRepository(db)

	item := &models.Item{InventoryNumber: "100007", Name: "ToDelete", Category: "Tools", Amount: 1}
	require.NoError(t, repo.Create(item))

	require.NoError(t, repo.Delete("100007"))

	_, err := repo.FindByInventoryNumber("100007")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
