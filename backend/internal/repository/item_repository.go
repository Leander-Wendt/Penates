package repository

import (
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// ItemFilter narrows an item search by optional search term, category, and location.
type ItemFilter struct {
	Search   string
	Category string
	Location string
}

// ItemRepository defines persistence operations for Item.
type ItemRepository interface {
	Create(item *models.Item) error
	FindByInventoryNumber(inventoryNumber string) (*models.Item, error)
	FindAll(filter ItemFilter) ([]models.Item, error)
	Update(item *models.Item) error
	Delete(inventoryNumber string) error
}

// GormItemRepository is a GORM-backed ItemRepository.
type GormItemRepository struct {
	db *gorm.DB
}

// NewGormItemRepository builds a GormItemRepository over the given *gorm.DB.
func NewGormItemRepository(db *gorm.DB) *GormItemRepository {
	return &GormItemRepository{db: db}
}

// Create persists a new Item.
func (r *GormItemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

// FindByInventoryNumber looks up an Item by its inventory number.
func (r *GormItemRepository) FindByInventoryNumber(inventoryNumber string) (*models.Item, error) {
	var item models.Item
	if err := r.db.First(&item, "inventory_number = ?", inventoryNumber).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// FindAll returns Items matching the given filter.
func (r *GormItemRepository) FindAll(filter ItemFilter) ([]models.Item, error) {
	query := r.db.Model(&models.Item{})
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ? OR inventory_number ILIKE ?", like, like, like)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Location != "" {
		query = query.Where("location = ?", filter.Location)
	}
	var items []models.Item
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Update saves changes to an existing Item.
func (r *GormItemRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

// Delete soft-deletes the Item with the given inventory number.
func (r *GormItemRepository) Delete(inventoryNumber string) error {
	return r.db.Delete(&models.Item{}, "inventory_number = ?", inventoryNumber).Error
}
