package service_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func TestUploadService_ReplaceItemImage_RejectsBadContentType(t *testing.T) {
	dir := t.TempDir()
	repo := new(mocks.ItemRepository)
	svc := service.NewUploadService(repo, dir, 10)

	_, err := svc.ReplaceItemImage("123456", "application/pdf", 100, strings.NewReader("data"))
	assert.ErrorIs(t, err, service.ErrValidation)
	repo.AssertNotCalled(t, "FindByInventoryNumber", mock.Anything)
}

func TestUploadService_ReplaceItemImage_RejectsOversizedFile(t *testing.T) {
	dir := t.TempDir()
	repo := new(mocks.ItemRepository)
	svc := service.NewUploadService(repo, dir, 1)

	oversize := int64(2 * 1024 * 1024)
	_, err := svc.ReplaceItemImage("123456", "image/png", oversize, strings.NewReader("data"))
	assert.ErrorIs(t, err, service.ErrValidation)
	repo.AssertNotCalled(t, "FindByInventoryNumber", mock.Anything)
}

func TestUploadService_ReplaceItemImage_SavesFileAndUpdatesItem(t *testing.T) {
	dir := t.TempDir()
	repo := new(mocks.ItemRepository)
	item := &models.Item{InventoryNumber: "123456", Name: "Drill"}
	repo.On("FindByInventoryNumber", "123456").Return(item, nil)
	var updated *models.Item
	repo.On("Update", mock.AnythingOfType("*models.Item")).Run(func(args mock.Arguments) {
		updated = args.Get(0).(*models.Item)
	}).Return(nil)
	svc := service.NewUploadService(repo, dir, 10)

	result, err := svc.ReplaceItemImage("123456", "image/png", 4, strings.NewReader("data"))
	require.NoError(t, err)
	assert.NotEmpty(t, result.ImagePath)
	assert.True(t, strings.HasPrefix(result.ImagePath, "123456-"))
	assert.True(t, strings.HasSuffix(result.ImagePath, ".png"))
	assert.Equal(t, result.ImagePath, updated.ImagePath)

	contents, err := os.ReadFile(filepath.Join(dir, result.ImagePath))
	require.NoError(t, err)
	assert.Equal(t, "data", string(contents))
}

func TestUploadService_ReplaceItemImage_DeletesPreviousFile(t *testing.T) {
	dir := t.TempDir()
	oldFile := filepath.Join(dir, "123456-old.png")
	require.NoError(t, os.WriteFile(oldFile, []byte("old"), 0o644))

	repo := new(mocks.ItemRepository)
	item := &models.Item{InventoryNumber: "123456", Name: "Drill", ImagePath: "123456-old.png"}
	repo.On("FindByInventoryNumber", "123456").Return(item, nil)
	repo.On("Update", mock.AnythingOfType("*models.Item")).Return(nil)
	svc := service.NewUploadService(repo, dir, 10)

	_, err := svc.ReplaceItemImage("123456", "image/png", 4, strings.NewReader("data"))
	require.NoError(t, err)

	_, statErr := os.Stat(oldFile)
	assert.True(t, os.IsNotExist(statErr))
}

func TestUploadService_DeleteItemImage_RemovesFileAndClearsPath(t *testing.T) {
	dir := t.TempDir()
	existingFile := filepath.Join(dir, "123456-existing.png")
	require.NoError(t, os.WriteFile(existingFile, []byte("data"), 0o644))

	repo := new(mocks.ItemRepository)
	item := &models.Item{InventoryNumber: "123456", ImagePath: "123456-existing.png"}
	repo.On("FindByInventoryNumber", "123456").Return(item, nil)
	var updated *models.Item
	repo.On("Update", mock.AnythingOfType("*models.Item")).Run(func(args mock.Arguments) {
		updated = args.Get(0).(*models.Item)
	}).Return(nil)
	svc := service.NewUploadService(repo, dir, 10)

	result, err := svc.DeleteItemImage("123456")
	require.NoError(t, err)
	assert.Empty(t, result.ImagePath)
	assert.Empty(t, updated.ImagePath)

	_, statErr := os.Stat(existingFile)
	assert.True(t, os.IsNotExist(statErr))
}
