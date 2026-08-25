package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/auth"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

// UserService implements business logic for User resources.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService builds a UserService over the given UserRepository.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// isValidRole reports whether role is one of the known Role values.
func isValidRole(role models.Role) bool {
	switch role {
	case models.RoleAdmin, models.RoleLogistics, models.RoleStudent:
		return true
	default:
		return false
	}
}

// Create validates the input, hashes the password, and persists a new User.
// The repository never receives a plaintext password.
func (s *UserService) Create(email, password string, role models.Role, organisationID uint) (*models.User, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return nil, ErrValidation
	}
	if !isValidRole(role) {
		return nil, ErrValidation
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:          email,
		PasswordHash:   hash,
		Role:           role,
		OrganisationID: organisationID,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID returns the User with the given ID.
func (s *UserService) GetByID(id uint) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// List returns every User.
func (s *UserService) List() ([]models.User, error) {
	return s.repo.FindAll()
}

// Update changes the Role and OrganisationID of an existing User.
func (s *UserService) Update(id uint, role models.Role, organisationID uint) (*models.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if !isValidRole(role) {
		return nil, ErrValidation
	}
	user.Role = role
	user.OrganisationID = organisationID
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Delete removes the User with the given ID.
func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}
