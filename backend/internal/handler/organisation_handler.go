package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// OrganisationService defines the operations OrganisationHandler needs from
// the service layer.
type OrganisationService interface {
	Create(name, description string) (*models.Organisation, error)
	GetByID(id uint) (*models.Organisation, error)
	List() ([]models.Organisation, error)
	Update(id uint, name, description string) (*models.Organisation, error)
	Delete(id uint) error
}

// OrganisationHandler exposes HTTP handlers for the Organisation resource.
type OrganisationHandler struct {
	svc OrganisationService
}

// NewOrganisationHandler builds an OrganisationHandler over the given OrganisationService.
func NewOrganisationHandler(svc OrganisationService) *OrganisationHandler {
	return &OrganisationHandler{svc: svc}
}

// organisationRequest is the JSON body accepted by Create and Update.
type organisationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// List handles GET /api/v1/organisations.
func (h *OrganisationHandler) List(c *gin.Context) {
	orgs, err := h.svc.List()
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, orgs)
}

// Create handles POST /api/v1/organisations.
func (h *OrganisationHandler) Create(c *gin.Context) {
	var req organisationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	org, err := h.svc.Create(req.Name, req.Description)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, org)
}

// Get handles GET /api/v1/organisations/:id.
func (h *OrganisationHandler) Get(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	org, err := h.svc.GetByID(id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, org)
}

// Update handles PUT /api/v1/organisations/:id.
func (h *OrganisationHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req organisationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	org, err := h.svc.Update(id, req.Name, req.Description)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, org)
}

// Delete handles DELETE /api/v1/organisations/:id.
func (h *OrganisationHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.svc.Delete(id); err != nil {
		writeServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// parseUintParam extracts and parses a uint URL parameter named name.
func parseUintParam(c *gin.Context, name string) (uint, error) {
	raw := c.Param(name)
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}
