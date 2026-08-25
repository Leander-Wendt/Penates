package repository

import (
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// UserRepository defines persistence operations for User.
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	Count() (int64, error)
}

// GormUserRepository is a GORM-backed UserRepository.
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository builds a GormUserRepository over the given *gorm.DB.
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Create persists a new User. The caller must have already hashed the password.
func (r *GormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID looks up a User by ID, preloading its Organisation.
func (r *GormUserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Organisation").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail looks up a User by email, preloading its Organisation.
func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Organisation").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll returns every User, preloading their Organisation.
func (r *GormUserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("Organisation").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update saves changes to an existing User.
func (r *GormUserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft-deletes the User with the given ID.
func (r *GormUserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// Count returns the total number of Users, including soft-deleted ones excluded by default scope.
func (r *GormUserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
