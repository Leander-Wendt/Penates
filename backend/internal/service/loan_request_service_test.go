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

func TestLoanRequestService_Create_DefaultsToPending(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	repo.On("Create", mock.AnythingOfType("*models.LoanRequest"), []string{"123456"}).
		Run(func(args mock.Arguments) {
			lr := args.Get(0).(*models.LoanRequest)
			lr.ID = 1
		}).Return(nil)
	svc := service.NewLoanRequestService(repo)

	lr, err := svc.Create(1, []string{"123456"}, "Room A")
	require.NoError(t, err)
	assert.Equal(t, "pending", lr.Status)
}

func TestLoanRequestService_Create_NoItemsFails(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	svc := service.NewLoanRequestService(repo)

	_, err := svc.Create(1, []string{}, "Room A")
	assert.ErrorIs(t, err, service.ErrValidation)
	repo.AssertNotCalled(t, "Create")
}

func TestLoanRequestService_GetByID_NotFound(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	repo.On("FindByID", uint(1)).Return(nil, gorm.ErrRecordNotFound)
	svc := service.NewLoanRequestService(repo)

	_, err := svc.GetByID(1)
	assert.ErrorIs(t, err, service.ErrNotFound)
}

func TestLoanRequestService_ListForUser_AdminSeesAll(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	repo.On("FindAll").Return([]models.LoanRequest{{ID: 1}, {ID: 2}}, nil)
	svc := service.NewLoanRequestService(repo)

	list, err := svc.ListForUser(1, models.RoleAdmin)
	require.NoError(t, err)
	assert.Len(t, list, 2)
	repo.AssertNotCalled(t, "FindByRequestingUserID", mock.Anything)
}

func TestLoanRequestService_ListForUser_StudentSeesOwnOnly(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	repo.On("FindByRequestingUserID", uint(5)).Return([]models.LoanRequest{{ID: 1, RequestingUserID: 5}}, nil)
	svc := service.NewLoanRequestService(repo)

	list, err := svc.ListForUser(5, models.RoleStudent)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	repo.AssertNotCalled(t, "FindAll")
}

func TestLoanRequestService_Update_OwnerWhilePendingSucceeds(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending", LocationOfItems: "Old"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*models.LoanRequest")).Return(nil)
	svc := service.NewLoanRequestService(repo)

	updated, err := svc.Update(1, 5, models.RoleStudent, "New Location")
	require.NoError(t, err)
	assert.Equal(t, "New Location", updated.LocationOfItems)
}

func TestLoanRequestService_Update_NonOwnerForbidden(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	svc := service.NewLoanRequestService(repo)

	_, err := svc.Update(1, 6, models.RoleStudent, "New Location")
	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestLoanRequestService_Update_AdminForbiddenViaThisMethod(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	svc := service.NewLoanRequestService(repo)

	_, err := svc.Update(1, 999, models.RoleAdmin, "New Location")
	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestLoanRequestService_Update_NotPendingForbidden(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "approved"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	svc := service.NewLoanRequestService(repo)

	_, err := svc.Update(1, 5, models.RoleStudent, "New Location")
	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestLoanRequestService_Delete_OwnerWhilePendingSucceeds(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)
	svc := service.NewLoanRequestService(repo)

	require.NoError(t, svc.Delete(1, 5, models.RoleStudent))
}

func TestLoanRequestService_Delete_NotPendingForbidden(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "returned"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	svc := service.NewLoanRequestService(repo)

	err := svc.Delete(1, 5, models.RoleStudent)
	assert.ErrorIs(t, err, service.ErrForbidden)
}

func TestLoanRequestService_UpdateStatus_ValidTransition(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, Status: "pending"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*models.LoanRequest")).Return(nil)
	svc := service.NewLoanRequestService(repo)

	updated, err := svc.UpdateStatus(1, "approved")
	require.NoError(t, err)
	assert.Equal(t, "approved", updated.Status)
}

func TestLoanRequestService_UpdateStatus_InvalidStatus(t *testing.T) {
	repo := new(mocks.LoanRequestRepository)
	existing := &models.LoanRequest{ID: 1, Status: "pending"}
	repo.On("FindByID", uint(1)).Return(existing, nil)
	svc := service.NewLoanRequestService(repo)

	_, err := svc.UpdateStatus(1, "bogus")
	assert.ErrorIs(t, err, service.ErrValidation)
}
