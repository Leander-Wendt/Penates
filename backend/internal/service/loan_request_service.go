package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

var validLoanRequestStatuses = map[string]bool{
	"pending":   true,
	"approved":  true,
	"rejected":  true,
	"returned":  true,
	"cancelled": true,
}

// LoanRequestService implements business logic for LoanRequest resources.
type LoanRequestService struct {
	repo repository.LoanRequestRepository
}

// NewLoanRequestService builds a LoanRequestService over the given LoanRequestRepository.
func NewLoanRequestService(repo repository.LoanRequestRepository) *LoanRequestService {
	return &LoanRequestService{repo: repo}
}

// Create validates and persists a new LoanRequest with status "pending" for the given items.
func (s *LoanRequestService) Create(requestingUserID uint, itemInventoryNumbers []string, locationOfItems string) (*models.LoanRequest, error) {
	if len(itemInventoryNumbers) == 0 {
		return nil, ErrValidation
	}
	loanRequest := &models.LoanRequest{
		RequestingUserID: requestingUserID,
		LocationOfItems:  locationOfItems,
		Status:           "pending",
	}
	if err := s.repo.Create(loanRequest, itemInventoryNumbers); err != nil {
		return nil, err
	}
	return loanRequest, nil
}

// GetByID returns the LoanRequest with the given ID.
func (s *LoanRequestService) GetByID(id uint) (*models.LoanRequest, error) {
	loanRequest, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return loanRequest, nil
}

// ListForUser returns every LoanRequest for Admin/Logistics roles, or only the
// caller's own LoanRequests for a Student.
func (s *LoanRequestService) ListForUser(userID uint, role models.Role) ([]models.LoanRequest, error) {
	if role == models.RoleAdmin || role == models.RoleLogistics {
		return s.repo.FindAll()
	}
	return s.repo.FindByRequestingUserID(userID)
}

// canModify reports whether callerID may PUT/DELETE the given LoanRequest: they
// must be its owner and it must still be pending.
func canModify(loanRequest *models.LoanRequest, callerID uint) bool {
	return loanRequest.RequestingUserID == callerID && loanRequest.Status == "pending"
}

// Update changes the LocationOfItems of a LoanRequest. Only the owning user may
// do so, and only while the request's status is still "pending".
func (s *LoanRequestService) Update(id uint, callerID uint, callerRole models.Role, locationOfItems string) (*models.LoanRequest, error) {
	loanRequest, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !canModify(loanRequest, callerID) {
		return nil, ErrForbidden
	}
	loanRequest.LocationOfItems = locationOfItems
	if err := s.repo.Update(loanRequest); err != nil {
		return nil, err
	}
	return loanRequest, nil
}

// Delete removes a LoanRequest. Only the owning user may do so, and only while
// the request's status is still "pending".
func (s *LoanRequestService) Delete(id uint, callerID uint, callerRole models.Role) error {
	loanRequest, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if !canModify(loanRequest, callerID) {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// UpdateStatus transitions a LoanRequest to a new status. Intended for use by
// Admin/Logistics callers via the dedicated status endpoint.
func (s *LoanRequestService) UpdateStatus(id uint, status string) (*models.LoanRequest, error) {
	if !validLoanRequestStatuses[status] {
		return nil, ErrValidation
	}
	loanRequest, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	loanRequest.Status = status
	if err := s.repo.Update(loanRequest); err != nil {
		return nil, err
	}
	return loanRequest, nil
}
