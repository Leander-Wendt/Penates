package repository

import (
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// LoanRequestRepository defines persistence operations for LoanRequest.
type LoanRequestRepository interface {
	Create(loanRequest *models.LoanRequest, itemInventoryNumbers []string) error
	FindByID(id uint) (*models.LoanRequest, error)
	FindAll() ([]models.LoanRequest, error)
	FindByRequestingUserID(userID uint) ([]models.LoanRequest, error)
	Update(loanRequest *models.LoanRequest) error
	Delete(id uint) error
}

// GormLoanRequestRepository is a GORM-backed LoanRequestRepository.
type GormLoanRequestRepository struct {
	db *gorm.DB
}

// NewGormLoanRequestRepository builds a GormLoanRequestRepository over the given *gorm.DB.
func NewGormLoanRequestRepository(db *gorm.DB) *GormLoanRequestRepository {
	return &GormLoanRequestRepository{db: db}
}

// Create persists a new LoanRequest along with its explicit LoanRequestItem join rows.
func (r *GormLoanRequestRepository) Create(loanRequest *models.LoanRequest, itemInventoryNumbers []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Items").Create(loanRequest).Error; err != nil {
			return err
		}
		for _, inv := range itemInventoryNumbers {
			join := models.LoanRequestItem{LoanRequestID: loanRequest.ID, ItemInventoryNumber: inv}
			if err := tx.Create(&join).Error; err != nil {
				return err
			}
		}
		return r.loadItems(tx, loanRequest)
	})
}

func (r *GormLoanRequestRepository) loadItems(tx *gorm.DB, loanRequest *models.LoanRequest) error {
	var joins []models.LoanRequestItem
	if err := tx.Where("loan_request_id = ?", loanRequest.ID).Find(&joins).Error; err != nil {
		return err
	}
	if len(joins) == 0 {
		loanRequest.Items = []models.Item{}
		return nil
	}
	invNumbers := make([]string, 0, len(joins))
	for _, j := range joins {
		invNumbers = append(invNumbers, j.ItemInventoryNumber)
	}
	var items []models.Item
	if err := tx.Where("inventory_number IN ?", invNumbers).Find(&items).Error; err != nil {
		return err
	}
	loanRequest.Items = items
	return nil
}

// FindByID looks up a LoanRequest by ID, preloading its RequestingUser and Items.
func (r *GormLoanRequestRepository) FindByID(id uint) (*models.LoanRequest, error) {
	var loanRequest models.LoanRequest
	if err := r.db.Preload("RequestingUser").First(&loanRequest, id).Error; err != nil {
		return nil, err
	}
	if err := r.loadItems(r.db, &loanRequest); err != nil {
		return nil, err
	}
	return &loanRequest, nil
}

// FindAll returns every LoanRequest, preloading RequestingUser and Items.
func (r *GormLoanRequestRepository) FindAll() ([]models.LoanRequest, error) {
	var loanRequests []models.LoanRequest
	if err := r.db.Preload("RequestingUser").Find(&loanRequests).Error; err != nil {
		return nil, err
	}
	for i := range loanRequests {
		if err := r.loadItems(r.db, &loanRequests[i]); err != nil {
			return nil, err
		}
	}
	return loanRequests, nil
}

// FindByRequestingUserID returns every LoanRequest made by the given user.
func (r *GormLoanRequestRepository) FindByRequestingUserID(userID uint) ([]models.LoanRequest, error) {
	var loanRequests []models.LoanRequest
	if err := r.db.Preload("RequestingUser").Where("requesting_user_id = ?", userID).Find(&loanRequests).Error; err != nil {
		return nil, err
	}
	for i := range loanRequests {
		if err := r.loadItems(r.db, &loanRequests[i]); err != nil {
			return nil, err
		}
	}
	return loanRequests, nil
}

// Update saves changes to an existing LoanRequest, excluding its Items association.
func (r *GormLoanRequestRepository) Update(loanRequest *models.LoanRequest) error {
	return r.db.Omit("Items").Save(loanRequest).Error
}

// Delete soft-deletes the LoanRequest with the given ID.
func (r *GormLoanRequestRepository) Delete(id uint) error {
	return r.db.Delete(&models.LoanRequest{}, id).Error
}
