//go:build integration

package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/database"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil"
)

func setupLoanRequestTestDB(t *testing.T) (*gorm.DB, uint) {
	t.Helper()
	db := testutil.NewTestDB(t)
	require.NoError(t, database.AutoMigrate(db))
	org := models.Organisation{Name: "Test Org"}
	require.NoError(t, db.Create(&org).Error)
	user := models.User{Email: "student@example.com", PasswordHash: "hash", Role: models.RoleStudent, OrganisationID: org.ID}
	require.NoError(t, db.Create(&user).Error)
	require.NoError(t, db.Create(&models.Item{InventoryNumber: "200001", Name: "Item A", Amount: 1}).Error)
	require.NoError(t, db.Create(&models.Item{InventoryNumber: "200002", Name: "Item B", Amount: 1}).Error)
	return db, user.ID
}

func TestGormLoanRequestRepository_CreateWithItems(t *testing.T) {
	db, userID := setupLoanRequestTestDB(t)
	repo := repository.NewGormLoanRequestRepository(db)

	lr := &models.LoanRequest{RequestingUserID: userID, Status: "pending"}
	require.NoError(t, repo.Create(lr, []string{"200001", "200002"}))
	assert.NotZero(t, lr.ID)
	assert.Len(t, lr.Items, 2)

	var joinCount int64
	require.NoError(t, db.Model(&models.LoanRequestItem{}).Where("loan_request_id = ?", lr.ID).Count(&joinCount).Error)
	assert.Equal(t, int64(2), joinCount)
}

func TestGormLoanRequestRepository_FindByID(t *testing.T) {
	db, userID := setupLoanRequestTestDB(t)
	repo := repository.NewGormLoanRequestRepository(db)

	lr := &models.LoanRequest{RequestingUserID: userID, Status: "pending"}
	require.NoError(t, repo.Create(lr, []string{"200001"}))

	found, err := repo.FindByID(lr.ID)
	require.NoError(t, err)
	assert.Equal(t, "student@example.com", found.RequestingUser.Email)
	require.Len(t, found.Items, 1)
	assert.Equal(t, "200001", found.Items[0].InventoryNumber)
}

func TestGormLoanRequestRepository_FindByRequestingUserID(t *testing.T) {
	db, userID := setupLoanRequestTestDB(t)
	repo := repository.NewGormLoanRequestRepository(db)

	require.NoError(t, repo.Create(&models.LoanRequest{RequestingUserID: userID, Status: "pending"}, []string{"200001"}))

	found, err := repo.FindByRequestingUserID(userID)
	require.NoError(t, err)
	assert.Len(t, found, 1)
}

func TestGormLoanRequestRepository_Update(t *testing.T) {
	db, userID := setupLoanRequestTestDB(t)
	repo := repository.NewGormLoanRequestRepository(db)

	lr := &models.LoanRequest{RequestingUserID: userID, Status: "pending"}
	require.NoError(t, repo.Create(lr, []string{"200001"}))

	lr.Status = "approved"
	require.NoError(t, repo.Update(lr))

	found, err := repo.FindByID(lr.ID)
	require.NoError(t, err)
	assert.Equal(t, "approved", found.Status)
}

func TestGormLoanRequestRepository_Delete(t *testing.T) {
	db, userID := setupLoanRequestTestDB(t)
	repo := repository.NewGormLoanRequestRepository(db)

	lr := &models.LoanRequest{RequestingUserID: userID, Status: "pending"}
	require.NoError(t, repo.Create(lr, []string{"200001"}))

	require.NoError(t, repo.Delete(lr.ID))

	_, err := repo.FindByID(lr.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
