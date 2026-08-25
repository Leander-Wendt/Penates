package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Leander-Wendt/Penates/backend/internal/handler"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func newOrganisationTestRouter(svc *mocks.OrganisationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewOrganisationHandler(svc)
	r.GET("/organisations", h.List)
	r.POST("/organisations", h.Create)
	r.GET("/organisations/:id", h.Get)
	r.PUT("/organisations/:id", h.Update)
	r.DELETE("/organisations/:id", h.Delete)
	return r
}

func TestOrganisationHandler_List(t *testing.T) {
	svc := new(mocks.OrganisationService)
	svc.On("List").Return([]models.Organisation{{ID: 1, Name: "Acme"}}, nil)
	r := newOrganisationTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/organisations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body []models.Organisation
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Len(t, body, 1)
}

func TestOrganisationHandler_Create(t *testing.T) {
	svc := new(mocks.OrganisationService)
	svc.On("Create", "Acme", "desc").Return(&models.Organisation{ID: 1, Name: "Acme", Description: "desc"}, nil)
	r := newOrganisationTestRouter(svc)

	payload, _ := json.Marshal(map[string]string{"name": "Acme", "description": "desc"})
	req := httptest.NewRequest(http.MethodPost, "/organisations", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestOrganisationHandler_Create_InvalidBody(t *testing.T) {
	svc := new(mocks.OrganisationService)
	r := newOrganisationTestRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/organisations", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Create")
}

func TestOrganisationHandler_Get_NotFound(t *testing.T) {
	svc := new(mocks.OrganisationService)
	svc.On("GetByID", uint(99)).Return(nil, service.ErrNotFound)
	r := newOrganisationTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/organisations/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOrganisationHandler_Update(t *testing.T) {
	svc := new(mocks.OrganisationService)
	svc.On("Update", uint(1), "New", "New Desc").Return(&models.Organisation{ID: 1, Name: "New", Description: "New Desc"}, nil)
	r := newOrganisationTestRouter(svc)

	payload, _ := json.Marshal(map[string]string{"name": "New", "description": "New Desc"})
	req := httptest.NewRequest(http.MethodPut, "/organisations/1", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOrganisationHandler_Delete(t *testing.T) {
	svc := new(mocks.OrganisationService)
	svc.On("Delete", uint(1)).Return(nil)
	r := newOrganisationTestRouter(svc)

	req := httptest.NewRequest(http.MethodDelete, "/organisations/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestOrganisationHandler_Get_InvalidID(t *testing.T) {
	svc := new(mocks.OrganisationService)
	r := newOrganisationTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/organisations/not-a-number", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
