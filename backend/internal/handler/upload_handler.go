package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
)

// UploadService defines the operations UploadHandler needs from the service layer.
type UploadService interface {
	ReplaceItemImage(inventoryNumber, contentType string, size int64, reader io.Reader) (*models.Item, error)
	DeleteItemImage(inventoryNumber string) (*models.Item, error)
}

// UploadHandler exposes HTTP handlers for the Item image upload flow.
type UploadHandler struct {
	svc UploadService
}

// NewUploadHandler builds an UploadHandler over the given UploadService.
func NewUploadHandler(svc UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// Upload handles POST /api/v1/items/:inventoryNumber/image.
func (h *UploadHandler) Upload(c *gin.Context) {
	inventoryNumber := c.Param("inventoryNumber")
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing image file"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image file"})
		return
	}
	defer func() { _ = file.Close() }()

	contentType := fileHeader.Header.Get("Content-Type")
	item, err := h.svc.ReplaceItemImage(inventoryNumber, contentType, fileHeader.Size, file)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

// Delete handles DELETE /api/v1/items/:inventoryNumber/image.
func (h *UploadHandler) Delete(c *gin.Context) {
	inventoryNumber := c.Param("inventoryNumber")
	item, err := h.svc.DeleteItemImage(inventoryNumber)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}
