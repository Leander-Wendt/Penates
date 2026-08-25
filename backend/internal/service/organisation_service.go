package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

// OrganisationService implements business logic for Organisation resources.
type OrganisationService struct {
	repo repository.OrganisationRepository
}

// NewOrganisationService builds an OrganisationService over the given OrganisationRepository.
func NewOrganisationService(repo repository.OrganisationRepository) *OrganisationService {
	return &OrganisationService{repo: repo}
}

// Create validates and persists a new Organisation.
func (s *OrganisationService) Create(name, description string) (*models.Organisation, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrValidation
	}
	org := &models.Organisation{Name: name, Description: description}
	if err := s.repo.Create(org); err != nil {
		return nil, err
	}
	return org, nil
}

// GetByID returns the Organisation with the given ID.
func (s *OrganisationService) GetByID(id uint) (*models.Organisation, error) {
	org, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return org, nil
}

// List returns every Organisation.
func (s *OrganisationService) List() ([]models.Organisation, error) {
	return s.repo.FindAll()
}

// Update validates and applies changes to an existing Organisation.
func (s *OrganisationService) Update(id uint, name, description string) (*models.Organisation, error) {
	org, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrValidation
	}
	org.Name = name
	org.Description = description
	if err := s.repo.Update(org); err != nil {
		return nil, err
	}
	return org, nil
}

// Delete removes the Organisation with the given ID.
func (s *OrganisationService) Delete(id uint) error {
	return s.repo.Delete(id)
}
