package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

// ItemRepository is a testify mock implementing repository.ItemRepository.
type ItemRepository struct {
	mock.Mock
}

// Create records a call to Create.
func (m *ItemRepository) Create(item *models.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

// FindByInventoryNumber records a call to FindByInventoryNumber.
func (m *ItemRepository) FindByInventoryNumber(inventoryNumber string) (*models.Item, error) {
	args := m.Called(inventoryNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

// FindAll records a call to FindAll.
func (m *ItemRepository) FindAll(filter repository.ItemFilter) ([]models.Item, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Item), args.Error(1)
}

// Update records a call to Update.
func (m *ItemRepository) Update(item *models.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

// Delete records a call to Delete.
func (m *ItemRepository) Delete(inventoryNumber string) error {
	args := m.Called(inventoryNumber)
	return args.Error(0)
}
