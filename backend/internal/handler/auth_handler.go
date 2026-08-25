package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/service"
)

// AuthService defines the operations AuthHandler needs from the service layer.
type AuthService interface {
	Login(email, password string) (string, time.Time, error)
}

// AuthHandler exposes HTTP handlers for authentication.
type AuthHandler struct {
	svc AuthService
}

// NewAuthHandler builds an AuthHandler over the given AuthService.
func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// loginRequest is the JSON body accepted by Login.
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// loginResponse is the JSON body returned by a successful Login.
type loginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Login handles POST /api/v1/auth/login. It is the login action itself, not a
// protected resource, so a failed attempt returns a plain 401 without a
// WWW-Authenticate challenge.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	token, expiresAt, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, loginResponse{Token: token, ExpiresAt: expiresAt})
}
