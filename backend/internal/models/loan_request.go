package models

import (
	"time"

	"gorm.io/gorm"
)

// LoanRequest represents a request from a User to borrow one or more Items.
type LoanRequest struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	RequestingUserID uint           `gorm:"not null" json:"requestingUserId"`
	RequestingUser   User           `json:"requestingUser"`
	Items            []Item         `gorm:"-" json:"items"`
	DateOfLending    *time.Time     `json:"dateOfLending"`
	DateOfReturn     *time.Time     `json:"dateOfReturn"`
	LocationOfItems  string         `json:"locationOfItems"`
	Status           string         `gorm:"not null;default:'pending';check:status IN ('pending','approved','rejected','returned','cancelled')" json:"status"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// LoanRequestItem is the explicit join model linking a LoanRequest to an Item.
type LoanRequestItem struct {
	LoanRequestID       uint   `gorm:"primaryKey" json:"loanRequestId"`
	ItemInventoryNumber string `gorm:"primaryKey;type:varchar(6)" json:"itemInventoryNumber"`
}
