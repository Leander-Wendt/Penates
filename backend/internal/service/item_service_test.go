package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func TestItemService_Create_Success(t *testing.T) {
	repo := new(mocks.ItemRepository)
	repo.On("Create", mock.AnythingOfType("*models.Item")).Return(nil)
	svc := service.NewItemService(repo)

	item, err := svc.Create(models.Item{InventoryNumber: "123456", Name: "Drill", Amount: 1})
	require.NoError(t, err)
	assert.Equal(t, "123456", item.InventoryNumber)
}

func TestItemService_Create_InvalidInventoryNumber(t *testing.T) {
	repo := new(mocks.ItemRepository)
	svc := service.NewItemService(repo)

	_, err := svc.Create(models.Item{InventoryNumber: "12A456", Name: "Drill", Amount: 1})
	assert.ErrorIs(t, err, service.ErrValidation)
	repo.AssertNotCalled(t, "Create")
}

func TestItemService_Create_NegativeAmount(t *testing.T) {
	repo := new(mocks.ItemRepository)
	svc := service.NewItemService(repo)

	_, err := svc.Create(models.Item{InventoryNumber: "123456", Name: "Drill", Amount: -1})
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestItemService_Create_EmptyName(t *testing.T) {
	repo := new(mocks.ItemRepository)
	svc := service.NewItemService(repo)

	_, err := svc.Create(models.Item{InventoryNumber: "123456", Name: "", Amount: 1})
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestItemService_GetByInventoryNumber_NotFound(t *testing.T) {
	repo := new(mocks.ItemRepository)
	repo.On("FindByInventoryNumber", "999999").Return(nil, gorm.ErrRecordNotFound)
	svc := service.NewItemService(repo)

	_, err := svc.GetByInventoryNumber("999999")
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestItemService_List_PassesFilter(t *testing.T) {
	repo := new(mocks.ItemRepository)
	repo.On("FindAll", repository.ItemFilter{Category: "Tools"}).Return([]models.Item{{InventoryNumber: "123456"}}, nil)
	svc := service.NewItemService(repo)

	items, err := svc.List(repository.ItemFilter{Category: "Tools"})
	require.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestItemService_Update_NotFound(t *testing.T) {
	repo := new(mocks.ItemRepository)
	repo.On("FindByInventoryNumber", "123456").Return(nil, gorm.ErrRecordNotFound)
	svc := service.NewItemService(repo)

	_, err := svc.Update("123456", models.Item{Name: "New Name", Amount: 1})
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestItemService_Update_Success(t *testing.T) {
	repo := new(mocks.ItemRepository)
	existing := &models.Item{InventoryNumber: "123456", Name: "Old", Amount: 1}
	repo.On("FindByInventoryNumber", "123456").Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*models.Item")).Return(nil)
	svc := service.NewItemService(repo)

	updated, err := svc.Update("123456", models.Item{Name: "New", Amount: 2})
	require.NoError(t, err)
	assert.Equal(t, "New", updated.Name)
	assert.Equal(t, 2, updated.Amount)
}

func TestItemService_Delete(t *testing.T) {
	repo := new(mocks.ItemRepository)
	repo.On("Delete", "123456").Return(nil)
	svc := service.NewItemService(repo)

	require.NoError(t, svc.Delete("123456"))
}
