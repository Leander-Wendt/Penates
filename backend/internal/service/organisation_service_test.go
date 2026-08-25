package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func TestOrganisationService_Create(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	repo.On("Create", mock.AnythingOfType("*models.Organisation")).Return(nil)
	svc := service.NewOrganisationService(repo)

	org, err := svc.Create("Acme", "desc")
	require.NoError(t, err)
	assert.Equal(t, "Acme", org.Name)
	repo.AssertExpectations(t)
}

func TestOrganisationService_Create_EmptyNameFails(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	svc := service.NewOrganisationService(repo)

	_, err := svc.Create("", "desc")
	assert.ErrorIs(t, err, service.ErrValidation)
	repo.AssertNotCalled(t, "Create")
}

func TestOrganisationService_GetByID_NotFound(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	repo.On("FindByID", uint(42)).Return(nil, gorm.ErrRecordNotFound)
	svc := service.NewOrganisationService(repo)

	_, err := svc.GetByID(42)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestOrganisationService_GetByID_Found(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	repo.On("FindByID", uint(1)).Return(&models.Organisation{ID: 1, Name: "Acme"}, nil)
	svc := service.NewOrganisationService(repo)

	org, err := svc.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, "Acme", org.Name)
}

func TestOrganisationService_List(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	repo.On("FindAll").Return([]models.Organisation{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}, nil)
	svc := service.NewOrganisationService(repo)

	orgs, err := svc.List()
	require.NoError(t, err)
	assert.Len(t, orgs, 2)
}

func TestOrganisationService_Update_NotFound(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	repo.On("FindByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)
	svc := service.NewOrganisationService(repo)

	_, err := svc.Update(1, "New Name", "New Desc")
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestOrganisationService_Update_Success(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	existing := &models.Organisation{ID: 1, Name: "Old"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*models.Organisation")).Return(nil)
	svc := service.NewOrganisationService(repo)

	updated, err := svc.Update(1, "New Name", "New Desc")
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
}

func TestOrganisationService_Delete(t *testing.T) {
	repo := new(mocks.OrganisationRepository)
	repo.On("Delete", uint(1)).Return(nil)
	svc := service.NewOrganisationService(repo)

	err := svc.Delete(1)
	require.NoError(t, err)
}
