package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// OrganisationRepository is a testify mock implementing repository.OrganisationRepository.
type OrganisationRepository struct {
	mock.Mock
}

// Create records a call to Create.
func (m *OrganisationRepository) Create(org *models.Organisation) error {
	args := m.Called(org)
	return args.Error(0)
}

// FindByID records a call to FindByID.
func (m *OrganisationRepository) FindByID(id uint) (*models.Organisation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organisation), args.Error(1)
}

// FindAll records a call to FindAll.
func (m *OrganisationRepository) FindAll() ([]models.Organisation, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Organisation), args.Error(1)
}

// Update records a call to Update.
func (m *OrganisationRepository) Update(org *models.Organisation) error {
	args := m.Called(org)
	return args.Error(0)
}

// Delete records a call to Delete.
func (m *OrganisationRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
