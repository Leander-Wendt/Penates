package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// UserRepository is a testify mock implementing repository.UserRepository.
type UserRepository struct {
	mock.Mock
}

// Create records a call to Create.
func (m *UserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindByID records a call to FindByID.
func (m *UserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// FindByEmail records a call to FindByEmail.
func (m *UserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// FindAll records a call to FindAll.
func (m *UserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

// Update records a call to Update.
func (m *UserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// Delete records a call to Delete.
func (m *UserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Count records a call to Count.
func (m *UserRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
