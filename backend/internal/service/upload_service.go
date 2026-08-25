package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
)

var allowedImageContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// UploadService implements the item image upload flow: content-type/size
// validation, saving files to disk, and updating Item.ImagePath.
type UploadService struct {
	itemRepo     repository.ItemRepository
	uploadDir    string
	maxSizeBytes int64
}

// NewUploadService builds an UploadService writing to uploadDir and rejecting
// files larger than maxSizeMB megabytes.
func NewUploadService(itemRepo repository.ItemRepository, uploadDir string, maxSizeMB int64) *UploadService {
	return &UploadService{
		itemRepo:     itemRepo,
		uploadDir:    uploadDir,
		maxSizeBytes: maxSizeMB * 1024 * 1024,
	}
}

// ReplaceItemImage validates the given content type and size, saves the file
// under a generated name, updates the Item's ImagePath, and deletes any
// previously stored image file.
func (s *UploadService) ReplaceItemImage(inventoryNumber, contentType string, size int64, reader io.Reader) (*models.Item, error) {
	ext, ok := allowedImageContentTypes[contentType]
	if !ok {
		return nil, ErrValidation
	}
	if size > s.maxSizeBytes {
		return nil, ErrValidation
	}

	item, err := s.itemRepo.FindByInventoryNumber(inventoryNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return nil, err
	}
	filename := inventoryNumber + "-" + uuid.NewString() + ext
	fullPath := filepath.Join(s.uploadDir, filename)
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(out, reader); err != nil {
		_ = out.Close()
		return nil, err
	}
	if err := out.Close(); err != nil {
		return nil, err
	}

	previousPath := item.ImagePath
	item.ImagePath = filename
	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}
	if previousPath != "" {
		_ = os.Remove(filepath.Join(s.uploadDir, previousPath))
	}
	return item, nil
}

// DeleteItemImage removes the Item's stored image file, if any, and clears its ImagePath.
func (s *UploadService) DeleteItemImage(inventoryNumber string) (*models.Item, error) {
	item, err := s.itemRepo.FindByInventoryNumber(inventoryNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if item.ImagePath != "" {
		_ = os.Remove(filepath.Join(s.uploadDir, item.ImagePath))
		item.ImagePath = ""
		if err := s.itemRepo.Update(item); err != nil {
			return nil, err
		}
	}
	return item, nil
}
