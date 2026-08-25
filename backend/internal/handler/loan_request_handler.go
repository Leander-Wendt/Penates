package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// LoanRequestService defines the operations LoanRequestHandler needs from the service layer.
type LoanRequestService interface {
	Create(requestingUserID uint, itemInventoryNumbers []string, locationOfItems string) (*models.LoanRequest, error)
	GetByID(id uint) (*models.LoanRequest, error)
	ListForUser(userID uint, role models.Role) ([]models.LoanRequest, error)
	Update(id uint, callerID uint, callerRole models.Role, locationOfItems string) (*models.LoanRequest, error)
	Delete(id uint, callerID uint, callerRole models.Role) error
	UpdateStatus(id uint, status string) (*models.LoanRequest, error)
}

// LoanRequestHandler exposes HTTP handlers for the LoanRequest resource.
type LoanRequestHandler struct {
	svc LoanRequestService
}

// NewLoanRequestHandler builds a LoanRequestHandler over the given LoanRequestService.
func NewLoanRequestHandler(svc LoanRequestService) *LoanRequestHandler {
	return &LoanRequestHandler{svc: svc}
}

// createLoanRequestRequest is the JSON body accepted by Create.
type createLoanRequestRequest struct {
	ItemInventoryNumbers []string `json:"itemInventoryNumbers" binding:"required"`
	LocationOfItems      string   `json:"locationOfItems"`
}

// updateLoanRequestRequest is the JSON body accepted by Update.
type updateLoanRequestRequest struct {
	LocationOfItems string `json:"locationOfItems"`
}

// updateLoanRequestStatusRequest is the JSON body accepted by UpdateStatus.
type updateLoanRequestStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// List handles GET /api/v1/loan-requests. Admin/Logistics see every
// LoanRequest; a Student sees only their own.
func (h *LoanRequestHandler) List(c *gin.Context) {
	loanRequests, err := h.svc.ListForUser(currentUserID(c), currentRole(c))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, loanRequests)
}

// Create handles POST /api/v1/loan-requests.
func (h *LoanRequestHandler) Create(c *gin.Context) {
	var req createLoanRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	loanRequest, err := h.svc.Create(currentUserID(c), req.ItemInventoryNumbers, req.LocationOfItems)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, loanRequest)
}

// Get handles GET /api/v1/loan-requests/:id. Admin/Logistics may fetch any
// LoanRequest; a Student may only fetch their own.
func (h *LoanRequestHandler) Get(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	loanRequest, err := h.svc.GetByID(id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	role := currentRole(c)
	if role != models.RoleAdmin && role != models.RoleLogistics && loanRequest.RequestingUserID != currentUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	c.JSON(http.StatusOK, loanRequest)
}

// Update handles PUT /api/v1/loan-requests/:id. The service enforces that
// only the owning user may update it, and only while it is still pending.
func (h *LoanRequestHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateLoanRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	loanRequest, err := h.svc.Update(id, currentUserID(c), currentRole(c), req.LocationOfItems)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, loanRequest)
}

// Delete handles DELETE /api/v1/loan-requests/:id. The service enforces that
// only the owning user may delete it, and only while it is still pending.
func (h *LoanRequestHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(id, currentUserID(c), currentRole(c)); err != nil {
		writeServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateStatus handles PATCH /api/v1/loan-requests/:id/status.
func (h *LoanRequestHandler) UpdateStatus(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateLoanRequestStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	loanRequest, err := h.svc.UpdateStatus(id, req.Status)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, loanRequest)
}
