package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// UserService defines the operations UserHandler needs from the service layer.
type UserService interface {
	Create(email, password string, role models.Role, organisationID uint) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	List() ([]models.User, error)
	Update(id uint, role models.Role, organisationID uint) (*models.User, error)
	Delete(id uint) error
}

// UserHandler exposes HTTP handlers for the User resource.
type UserHandler struct {
	svc UserService
}

// NewUserHandler builds a UserHandler over the given UserService.
func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// createUserRequest is the JSON body accepted by Create.
type createUserRequest struct {
	Email          string      `json:"email" binding:"required"`
	Password       string      `json:"password" binding:"required"`
	Role           models.Role `json:"role" binding:"required"`
	OrganisationID uint        `json:"organisationId" binding:"required"`
}

// updateUserRequest is the JSON body accepted by Update.
type updateUserRequest struct {
	Role           models.Role `json:"role" binding:"required"`
	OrganisationID uint        `json:"organisationId" binding:"required"`
}

// Me handles GET /api/v1/users/me, returning the authenticated caller's own User record.
func (h *UserHandler) Me(c *gin.Context) {
	user, err := h.svc.GetByID(currentUserID(c))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// List handles GET /api/v1/users.
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.List()
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// Create handles POST /api/v1/users.
func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.svc.Create(req.Email, req.Password, req.Role, req.OrganisationID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Get handles GET /api/v1/users/:id. Admin/Logistics may fetch any user;
// a Student may only fetch their own record.
func (h *UserHandler) Get(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	role := currentRole(c)
	if role != models.RoleAdmin && role != models.RoleLogistics && id != currentUserID(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	user, err := h.svc.GetByID(id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Update handles PUT /api/v1/users/:id.
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.svc.Update(id, req.Role, req.OrganisationID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// Delete handles DELETE /api/v1/users/:id.
func (h *UserHandler) Delete(c *gin.Context) {
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
