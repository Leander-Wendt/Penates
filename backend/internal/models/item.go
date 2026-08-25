package models

import (
	"time"

	"gorm.io/gorm"
)

// Item represents a piece of inventory that can be loaned out.
type Item struct {
	InventoryNumber              string         `gorm:"primaryKey;type:varchar(6);check:inventory_number ~ '^[0-9]{6}$'" json:"inventoryNumber"`
	Name                         string         `gorm:"not null" json:"name"`
	Description                  string         `json:"description"`
	ImagePath                    string         `json:"imagePath"`
	Category                     string         `gorm:"index" json:"category"`
	Location                     string         `json:"location"`
	Amount                       int            `gorm:"not null;default:1;check:amount >= 0" json:"amount"`
	Manufacturer                 string         `json:"manufacturer"`
	SerialNumber                 string         `json:"serialNumber"`
	Note                         string         `json:"note"`
	LastTechnicalInspectionDate  *time.Time     `gorm:"type:date" json:"lastTechnicalInspectionDate"`
	ElectricalAppliance          bool           `gorm:"not null;default:false" json:"electricalAppliance"`
	LastElectricalInspectionDate *time.Time     `gorm:"type:date" json:"lastElectricalInspectionDate"`
	DateOfPurchase               *time.Time     `gorm:"type:date" json:"dateOfPurchase"`
	Price                        float64        `gorm:"type:numeric(10,2)" json:"price"`
	ResolutionNumber             string         `json:"resolutionNumber"`
	CreatedAt                    time.Time      `json:"createdAt"`
	UpdatedAt                    time.Time      `json:"updatedAt"`
	DeletedAt                    gorm.DeletedAt `gorm:"index" json:"-"`
}
