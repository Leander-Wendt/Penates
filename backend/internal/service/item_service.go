package service

import (
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

var inventoryNumberPattern = regexp.MustCompile(`^[0-9]{6}$`)

// ItemService implements business logic for Item resources.
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService builds an ItemService over the given ItemRepository.
func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// validate checks the invariants shared by Create and Update.
func validateItem(item models.Item) error {
	if !inventoryNumberPattern.MatchString(item.InventoryNumber) {
		return ErrValidation
	}
	if strings.TrimSpace(item.Name) == "" {
		return ErrValidation
	}
	if item.Amount < 0 {
		return ErrValidation
	}
	return nil
}

// Create validates and persists a new Item.
func (s *ItemService) Create(item models.Item) (*models.Item, error) {
	if err := validateItem(item); err != nil {
		return nil, err
	}
	if err := s.repo.Create(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

// GetByInventoryNumber returns the Item with the given inventory number.
func (s *ItemService) GetByInventoryNumber(inventoryNumber string) (*models.Item, error) {
	item, err := s.repo.FindByInventoryNumber(inventoryNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item, nil
}

// List returns Items matching the given filter.
func (s *ItemService) List(filter repository.ItemFilter) ([]models.Item, error) {
	return s.repo.FindAll(filter)
}

// Update applies changes to an existing Item, keeping its inventory number and image path.
func (s *ItemService) Update(inventoryNumber string, updates models.Item) (*models.Item, error) {
	existing, err := s.GetByInventoryNumber(inventoryNumber)
	if err != nil {
		return nil, err
	}
	updates.InventoryNumber = existing.InventoryNumber
	if err := validateItem(updates); err != nil {
		return nil, err
	}
	updates.ImagePath = existing.ImagePath
	if err := s.repo.Update(&updates); err != nil {
		return nil, err
	}
	return &updates, nil
}

// Delete removes the Item with the given inventory number.
func (s *ItemService) Delete(inventoryNumber string) error {
	return s.repo.Delete(inventoryNumber)
}
