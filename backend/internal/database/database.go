package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// Open connects to the Postgres database at the given DSN and returns a *gorm.DB.
func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

// AutoMigrate creates or updates all Penates tables in FK dependency order:
// Organisation, User, Item, LoanRequest, LoanRequestItem.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Organisation{},
		&models.User{},
		&models.Item{},
		&models.LoanRequest{},
		&models.LoanRequestItem{},
	)
}

// BootstrapAdmin creates an initial Admin user from the given credentials if
// the users table is empty. It attaches the admin to an Organisation named
// orgName, creating it if it does not already exist. If any users already
// exist, BootstrapAdmin is a no-op.
func BootstrapAdmin(db *gorm.DB, email, password, orgName string) error {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return nil
	}

	var org models.Organisation
	err := db.Where("name = ?", orgName).First(&org).Error
	if err == gorm.ErrRecordNotFound {
		org = models.Organisation{Name: orgName}
		if err := db.Create(&org).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	admin := models.User{
		Email:          email,
		PasswordHash:   hash,
		Role:           models.RoleAdmin,
		OrganisationID: org.ID,
	}
	return db.Create(&admin).Error
}
