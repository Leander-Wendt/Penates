package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/dto"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

// ItemService defines the operations ItemHandler needs from the service layer.
type ItemService interface {
	Create(item models.Item) (*models.Item, error)
	GetByInventoryNumber(inventoryNumber string) (*models.Item, error)
	List(filter repository.ItemFilter) ([]models.Item, error)
	Update(inventoryNumber string, updates models.Item) (*models.Item, error)
	Delete(inventoryNumber string) error
}

// ItemHandler exposes HTTP handlers for the Item resource. Every response is
// built through dto.NewItemDTO, keyed off the caller's role, so that
// Admin/Logistics-only fields are never serialized for a Student.
type ItemHandler struct {
	svc ItemService
}

// NewItemHandler builds an ItemHandler over the given ItemService.
func NewItemHandler(svc ItemService) *ItemHandler {
	return &ItemHandler{svc: svc}
}

// itemRequest is the JSON body accepted by Create and Update.
type itemRequest struct {
	InventoryNumber              string     `json:"inventoryNumber"`
	Name                         string     `json:"name" binding:"required"`
	Description                  string     `json:"description"`
	Category                     string     `json:"category"`
	Location                     string     `json:"location"`
	Amount                       int        `json:"amount"`
	Manufacturer                 string     `json:"manufacturer"`
	SerialNumber                 string     `json:"serialNumber"`
	Note                         string     `json:"note"`
	LastTechnicalInspectionDate  *time.Time `json:"lastTechnicalInspectionDate"`
	ElectricalAppliance          bool       `json:"electricalAppliance"`
	LastElectricalInspectionDate *time.Time `json:"lastElectricalInspectionDate"`
	DateOfPurchase               *time.Time `json:"dateOfPurchase"`
	Price                        float64    `json:"price"`
	ResolutionNumber             string     `json:"resolutionNumber"`
}

// toModel converts an itemRequest into a models.Item.
func (req itemRequest) toModel() models.Item {
	return models.Item{
		InventoryNumber:              req.InventoryNumber,
		Name:                         req.Name,
		Description:                  req.Description,
		Category:                     req.Category,
		Location:                     req.Location,
		Amount:                       req.Amount,
		Manufacturer:                 req.Manufacturer,
		SerialNumber:                 req.SerialNumber,
		Note:                         req.Note,
		LastTechnicalInspectionDate:  req.LastTechnicalInspectionDate,
		ElectricalAppliance:          req.ElectricalAppliance,
		LastElectricalInspectionDate: req.LastElectricalInspectionDate,
		DateOfPurchase:               req.DateOfPurchase,
		Price:                        req.Price,
		ResolutionNumber:             req.ResolutionNumber,
	}
}

// List handles GET /api/v1/items.
func (h *ItemHandler) List(c *gin.Context) {
	filter := repository.ItemFilter{
		Search:   c.Query("search"),
		Category: c.Query("category"),
		Location: c.Query("location"),
	}
	items, err := h.svc.List(filter)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewItemDTOs(items, currentRole(c)))
}

// Create handles POST /api/v1/items.
func (h *ItemHandler) Create(c *gin.Context) {
	var req itemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	item, err := h.svc.Create(req.toModel())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewItemDTO(*item, currentRole(c)))
}

// Get handles GET /api/v1/items/:inventoryNumber.
func (h *ItemHandler) Get(c *gin.Context) {
	inventoryNumber := c.Param("inventoryNumber")
	item, err := h.svc.GetByInventoryNumber(inventoryNumber)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewItemDTO(*item, currentRole(c)))
}

// Update handles PUT /api/v1/items/:inventoryNumber.
func (h *ItemHandler) Update(c *gin.Context) {
	inventoryNumber := c.Param("inventoryNumber")
	var req itemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	item, err := h.svc.Update(inventoryNumber, req.toModel())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewItemDTO(*item, currentRole(c)))
}

// Delete handles DELETE /api/v1/items/:inventoryNumber.
func (h *ItemHandler) Delete(c *gin.Context) {
	inventoryNumber := c.Param("inventoryNumber")
	if err := h.svc.Delete(inventoryNumber); err != nil {
		writeServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
