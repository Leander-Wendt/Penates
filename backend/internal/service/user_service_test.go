package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func TestUserService_Create_HashesPassword(t *testing.T) {
	repo := new(mocks.UserRepository)
	var captured *models.User
	repo.On("Create", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		captured = args.Get(0).(*models.User)
	}).Return(nil)
	svc := service.NewUserService(repo)

	user, err := svc.Create("new@example.com", "plaintext-pw", models.RoleStudent, 1)
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", user.Email)
	assert.NotEqual(t, "plaintext-pw", captured.PasswordHash)
	assert.True(t, auth.CheckPasswordHash("plaintext-pw", captured.PasswordHash))
	repo.AssertExpectations(t)
}

func TestUserService_Create_InvalidRoleFails(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := service.NewUserService(repo)

	_, err := svc.Create("new@example.com", "plaintext-pw", models.Role("bogus"), 1)
	assert.ErrorIs(t, err, service.ErrValidation)
	repo.AssertNotCalled(t, "Create")
}

func TestUserService_Create_EmptyEmailFails(t *testing.T) {
	repo := new(mocks.UserRepository)
	svc := service.NewUserService(repo)

	_, err := svc.Create("", "plaintext-pw", models.RoleStudent, 1)
	assert.ErrorIs(t, err, service.ErrValidation)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("FindByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)
	svc := service.NewUserService(repo)

	_, err := svc.GetByID(1)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestUserService_List(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("FindAll").Return([]models.User{{ID: 1}, {ID: 2}}, nil)
	svc := service.NewUserService(repo)

	users, err := svc.List()
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserService_Update_ChangesRoleAndOrg(t *testing.T) {
	repo := new(mocks.UserRepository)
	existing := &models.User{ID: 1, Email: "a@example.com", Role: models.RoleStudent, OrganisationID: 1}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
	svc := service.NewUserService(repo)

	updated, err := svc.Update(1, models.RoleLogistics, 2)
	require.NoError(t, err)
	assert.Equal(t, models.RoleLogistics, updated.Role)
	assert.Equal(t, uint(2), updated.OrganisationID)
}

func TestUserService_Delete(t *testing.T) {
	repo := new(mocks.UserRepository)
	repo.On("Delete", uint(1)).Return(nil)
	svc := service.NewUserService(repo)

	require.NoError(t, svc.Delete(1))
}
