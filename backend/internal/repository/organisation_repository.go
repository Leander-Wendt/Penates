package repository

import (
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// OrganisationRepository defines persistence operations for Organisation.
type OrganisationRepository interface {
	Create(org *models.Organisation) error
	FindByID(id uint) (*models.Organisation, error)
	FindAll() ([]models.Organisation, error)
	Update(org *models.Organisation) error
	Delete(id uint) error
}

// GormOrganisationRepository is a GORM-backed OrganisationRepository.
type GormOrganisationRepository struct {
	db *gorm.DB
}

// NewGormOrganisationRepository builds a GormOrganisationRepository over the given *gorm.DB.
func NewGormOrganisationRepository(db *gorm.DB) *GormOrganisationRepository {
	return &GormOrganisationRepository{db: db}
}

// Create persists a new Organisation.
func (r *GormOrganisationRepository) Create(org *models.Organisation) error {
	return r.db.Create(org).Error
}

// FindByID looks up an Organisation by its ID.
func (r *GormOrganisationRepository) FindByID(id uint) (*models.Organisation, error) {
	var org models.Organisation
	if err := r.db.First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// FindAll returns every Organisation.
func (r *GormOrganisationRepository) FindAll() ([]models.Organisation, error) {
	var orgs []models.Organisation
	if err := r.db.Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

// Update saves changes to an existing Organisation.
func (r *GormOrganisationRepository) Update(org *models.Organisation) error {
	return r.db.Save(org).Error
}

// Delete soft-deletes the Organisation with the given ID.
func (r *GormOrganisationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Organisation{}, id).Error
}
