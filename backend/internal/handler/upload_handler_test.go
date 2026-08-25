package handler_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Leander-Wendt/Penates/backend/internal/handler"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func newUploadTestRouter(svc handler.UploadService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewUploadHandler(svc)
	r.POST("/items/:inventoryNumber/image", h.Upload)
	r.DELETE("/items/:inventoryNumber/image", h.Delete)
	return r
}

func buildMultipartRequest(t *testing.T, fieldName, filename, contentType string, content []byte) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreatePart(map[string][]string{
		"Content-Disposition": {`form-data; name="` + fieldName + `"; filename="` + filename + `"`},
		"Content-Type":        {contentType},
	})
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return body, writer.FormDataContentType()
}

func TestUploadHandler_Upload_SavesFileAndUpdatesItem(t *testing.T) {
	dir := t.TempDir()
	itemRepo := new(mocks.ItemRepository)
	item := &models.Item{InventoryNumber: "123456", Name: "Drill"}
	itemRepo.On("FindByInventoryNumber", "123456").Return(item, nil)
	itemRepo.On("Update", mock.AnythingOfType("*models.Item")).Return(nil)
	uploadSvc := service.NewUploadService(itemRepo, dir, 10)
	r := newUploadTestRouter(uploadSvc)

	body, contentType := buildMultipartRequest(t, "image", "photo.png", "image/png", []byte("fake-image-bytes"))
	req := httptest.NewRequest(http.MethodPost, "/items/123456/image", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.True(t, strings.HasPrefix(entries[0].Name(), "123456-"))

	contents, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	require.NoError(t, err)
	assert.Equal(t, "fake-image-bytes", string(contents))

	itemRepo.AssertExpectations(t)
}

func TestUploadHandler_Upload_RejectsBadContentType(t *testing.T) {
	dir := t.TempDir()
	itemRepo := new(mocks.ItemRepository)
	uploadSvc := service.NewUploadService(itemRepo, dir, 10)
	r := newUploadTestRouter(uploadSvc)

	body, contentType := buildMultipartRequest(t, "image", "doc.pdf", "application/pdf", []byte("not-an-image"))
	req := httptest.NewRequest(http.MethodPost, "/items/123456/image", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	itemRepo.AssertNotCalled(t, "FindByInventoryNumber", mock.Anything)
}

func TestUploadHandler_Delete(t *testing.T) {
	svc := new(mocks.UploadService)
	svc.On("DeleteItemImage", "123456").Return(&models.Item{InventoryNumber: "123456"}, nil)
	r := newUploadTestRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/items/123456/image", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
