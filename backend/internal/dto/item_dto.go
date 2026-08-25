package dto

import (
	"time"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// ItemDTO is the API representation of an Item, with fields restricted to
// Admin/Logistics roles omitted for other roles.
type ItemDTO struct {
	InventoryNumber              string     `json:"inventoryNumber"`
	Name                         string     `json:"name"`
	Description                  string     `json:"description"`
	ImagePath                    string     `json:"imagePath"`
	Category                     string     `json:"category"`
	Location                     string     `json:"location"`
	Amount                       int        `json:"amount"`
	Note                         string     `json:"note"`
	ElectricalAppliance          bool       `json:"electricalAppliance"`
	CreatedAt                    time.Time  `json:"createdAt"`
	UpdatedAt                    time.Time  `json:"updatedAt"`
	Manufacturer                 *string    `json:"manufacturer,omitempty"`
	SerialNumber                 *string    `json:"serialNumber,omitempty"`
	LastTechnicalInspectionDate  *time.Time `json:"lastTechnicalInspectionDate,omitempty"`
	LastElectricalInspectionDate *time.Time `json:"lastElectricalInspectionDate,omitempty"`
	DateOfPurchase               *time.Time `json:"dateOfPurchase,omitempty"`
	Price                        *float64   `json:"price,omitempty"`
	ResolutionNumber             *string    `json:"resolutionNumber,omitempty"`
}

// canSeeRestrictedFields reports whether the given role is allowed to see the
// Admin/Logistics-only Item fields.
func canSeeRestrictedFields(role models.Role) bool {
	return role == models.RoleAdmin || role == models.RoleLogistics
}

// NewItemDTO builds an ItemDTO from an Item, populating the restricted fields
// only when role is Admin or Logistics.
func NewItemDTO(item models.Item, role models.Role) ItemDTO {
	d := ItemDTO{
		InventoryNumber:     item.InventoryNumber,
		Name:                item.Name,
		Description:         item.Description,
		ImagePath:           item.ImagePath,
		Category:            item.Category,
		Location:            item.Location,
		Amount:              item.Amount,
		Note:                item.Note,
		ElectricalAppliance: item.ElectricalAppliance,
		CreatedAt:           item.CreatedAt,
		UpdatedAt:           item.UpdatedAt,
	}
	if canSeeRestrictedFields(role) {
		manufacturer := item.Manufacturer
		serialNumber := item.SerialNumber
		price := item.Price
		resolutionNumber := item.ResolutionNumber
		d.Manufacturer = &manufacturer
		d.SerialNumber = &serialNumber
		d.LastTechnicalInspectionDate = item.LastTechnicalInspectionDate
		d.LastElectricalInspectionDate = item.LastElectricalInspectionDate
		d.DateOfPurchase = item.DateOfPurchase
		d.Price = &price
		d.ResolutionNumber = &resolutionNumber
	}
	return d
}

// NewItemDTOs builds a slice of ItemDTO from a slice of Item for the given role.
func NewItemDTOs(items []models.Item, role models.Role) []ItemDTO {
	dtos := make([]ItemDTO, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, NewItemDTO(item, role))
	}
	return dtos
}
