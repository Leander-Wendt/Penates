package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents an account that can authenticate with the API.
type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string         `gorm:"not null" json:"-"`
	Role           Role           `gorm:"type:varchar(20);not null;check:role IN ('admin','logistics','student')" json:"role"`
	OrganisationID uint           `gorm:"not null" json:"organisationId"`
	Organisation   Organisation   `json:"organisation"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
