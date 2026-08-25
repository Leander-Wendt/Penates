package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// LoanRequestRepository is a testify mock implementing repository.LoanRequestRepository.
type LoanRequestRepository struct {
	mock.Mock
}

// Create records a call to Create.
func (m *LoanRequestRepository) Create(loanRequest *models.LoanRequest, itemInventoryNumbers []string) error {
	args := m.Called(loanRequest, itemInventoryNumbers)
	return args.Error(0)
}

// FindByID records a call to FindByID.
func (m *LoanRequestRepository) FindByID(id uint) (*models.LoanRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanRequest), args.Error(1)
}

// FindAll records a call to FindAll.
func (m *LoanRequestRepository) FindAll() ([]models.LoanRequest, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LoanRequest), args.Error(1)
}

// FindByRequestingUserID records a call to FindByRequestingUserID.
func (m *LoanRequestRepository) FindByRequestingUserID(userID uint) ([]models.LoanRequest, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LoanRequest), args.Error(1)
}

// Update records a call to Update.
func (m *LoanRequestRepository) Update(loanRequest *models.LoanRequest) error {
	args := m.Called(loanRequest)
	return args.Error(0)
}

// Delete records a call to Delete.
func (m *LoanRequestRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
