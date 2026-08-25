package mocks

import (
	"io"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

func returnOrgSlice(args mock.Arguments) ([]models.Organisation, error) {
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Organisation), args.Error(1)
}

// OrganisationService is a testify mock implementing the handler package's OrganisationService interface.
type OrganisationService struct {
	mock.Mock
}

// Create records a call to Create.
func (m *OrganisationService) Create(name, description string) (*models.Organisation, error) {
	args := m.Called(name, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organisation), args.Error(1)
}

// GetByID records a call to GetByID.
func (m *OrganisationService) GetByID(id uint) (*models.Organisation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organisation), args.Error(1)
}

// List records a call to List.
func (m *OrganisationService) List() ([]models.Organisation, error) {
	args := m.Called()
	return returnOrgSlice(args)
}

// Update records a call to Update.
func (m *OrganisationService) Update(id uint, name, description string) (*models.Organisation, error) {
	args := m.Called(id, name, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Organisation), args.Error(1)
}

// Delete records a call to Delete.
func (m *OrganisationService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// UserService is a testify mock implementing the handler package's UserService interface.
type UserService struct {
	mock.Mock
}

// Create records a call to Create.
func (m *UserService) Create(email, password string, role models.Role, organisationID uint) (*models.User, error) {
	args := m.Called(email, password, role, organisationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// GetByID records a call to GetByID.
func (m *UserService) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// List records a call to List.
func (m *UserService) List() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

// Update records a call to Update.
func (m *UserService) Update(id uint, role models.Role, organisationID uint) (*models.User, error) {
	args := m.Called(id, role, organisationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Delete records a call to Delete.
func (m *UserService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// AuthService is a testify mock implementing the handler package's AuthService interface.
type AuthService struct {
	mock.Mock
}

// Login records a call to Login.
func (m *AuthService) Login(email, password string) (string, time.Time, error) {
	args := m.Called(email, password)
	expiresAt, _ := args.Get(1).(time.Time)
	return args.String(0), expiresAt, args.Error(2)
}

// ItemService is a testify mock implementing the handler package's ItemService interface.
type ItemService struct {
	mock.Mock
}

// Create records a call to Create.
func (m *ItemService) Create(item models.Item) (*models.Item, error) {
	args := m.Called(item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

// GetByInventoryNumber records a call to GetByInventoryNumber.
func (m *ItemService) GetByInventoryNumber(inventoryNumber string) (*models.Item, error) {
	args := m.Called(inventoryNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

// List records a call to List.
func (m *ItemService) List(filter repository.ItemFilter) ([]models.Item, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Item), args.Error(1)
}

// Update records a call to Update.
func (m *ItemService) Update(inventoryNumber string, updates models.Item) (*models.Item, error) {
	args := m.Called(inventoryNumber, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

// Delete records a call to Delete.
func (m *ItemService) Delete(inventoryNumber string) error {
	args := m.Called(inventoryNumber)
	return args.Error(0)
}

// LoanRequestService is a testify mock implementing the handler package's LoanRequestService interface.
type LoanRequestService struct {
	mock.Mock
}

// Create records a call to Create.
func (m *LoanRequestService) Create(requestingUserID uint, itemInventoryNumbers []string, locationOfItems string) (*models.LoanRequest, error) {
	args := m.Called(requestingUserID, itemInventoryNumbers, locationOfItems)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanRequest), args.Error(1)
}

// GetByID records a call to GetByID.
func (m *LoanRequestService) GetByID(id uint) (*models.LoanRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanRequest), args.Error(1)
}

// ListForUser records a call to ListForUser.
func (m *LoanRequestService) ListForUser(userID uint, role models.Role) ([]models.LoanRequest, error) {
	args := m.Called(userID, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.LoanRequest), args.Error(1)
}

// Update records a call to Update.
func (m *LoanRequestService) Update(id uint, callerID uint, callerRole models.Role, locationOfItems string) (*models.LoanRequest, error) {
	args := m.Called(id, callerID, callerRole, locationOfItems)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanRequest), args.Error(1)
}

// Delete records a call to Delete.
func (m *LoanRequestService) Delete(id uint, callerID uint, callerRole models.Role) error {
	args := m.Called(id, callerID, callerRole)
	return args.Error(0)
}

// UpdateStatus records a call to UpdateStatus.
func (m *LoanRequestService) UpdateStatus(id uint, status string) (*models.LoanRequest, error) {
	args := m.Called(id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanRequest), args.Error(1)
}

// UploadService is a testify mock implementing the handler package's UploadService interface.
type UploadService struct {
	mock.Mock
}

// ReplaceItemImage records a call to ReplaceItemImage.
func (m *UploadService) ReplaceItemImage(inventoryNumber, contentType string, size int64, reader io.Reader) (*models.Item, error) {
	args := m.Called(inventoryNumber, contentType, size, reader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}

// DeleteItemImage records a call to DeleteItemImage.
func (m *UploadService) DeleteItemImage(inventoryNumber string) (*models.Item, error) {
	args := m.Called(inventoryNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Item), args.Error(1)
}
