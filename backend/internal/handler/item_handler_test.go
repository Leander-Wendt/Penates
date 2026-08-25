package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Leander-Wendt/Penates/backend/internal/handler"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/repository"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func newItemTestRouter(svc *mocks.ItemService, role models.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Set("role", role)
		c.Next()
	})
	h := handler.NewItemHandler(svc)
	r.GET("/items", h.List)
	r.POST("/items", h.Create)
	r.GET("/items/:inventoryNumber", h.Get)
	r.PUT("/items/:inventoryNumber", h.Update)
	r.DELETE("/items/:inventoryNumber", h.Delete)
	return r
}

func sampleTestItem() models.Item {
	return models.Item{
		InventoryNumber: "123456",
		Name:            "Drill",
		Manufacturer:    "Bosch",
		Price:           99.99,
	}
}

func TestItemHandler_Get_StudentRedactsFields(t *testing.T) {
	svc := new(mocks.ItemService)
	item := sampleTestItem()
	svc.On("GetByInventoryNumber", "123456").Return(&item, nil)
	r := newItemTestRouter(svc, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/items/123456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.NotContains(t, body, "manufacturer")
	assert.NotContains(t, body, "price")
}

func TestItemHandler_Get_AdminIncludesFields(t *testing.T) {
	svc := new(mocks.ItemService)
	item := sampleTestItem()
	svc.On("GetByInventoryNumber", "123456").Return(&item, nil)
	r := newItemTestRouter(svc, models.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/items/123456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Contains(t, body, "manufacturer")
	assert.Contains(t, body, "price")
}

func TestItemHandler_List_PassesQueryFilters(t *testing.T) {
	svc := new(mocks.ItemService)
	svc.On("List", repository.ItemFilter{Search: "drill", Category: "Tools", Location: "Room A"}).
		Return([]models.Item{sampleTestItem()}, nil)
	r := newItemTestRouter(svc, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/items?search=drill&category=Tools&location=Room+A", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestItemHandler_Create(t *testing.T) {
	svc := new(mocks.ItemService)
	item := sampleTestItem()
	svc.On("Create", mock.MatchedBy(func(i models.Item) bool { return i.InventoryNumber == "123456" })).Return(&item, nil)
	r := newItemTestRouter(svc, models.RoleAdmin)

	payload, _ := json.Marshal(map[string]interface{}{
		"inventoryNumber": "123456", "name": "Drill", "amount": 1,
	})
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestItemHandler_Delete(t *testing.T) {
	svc := new(mocks.ItemService)
	svc.On("Delete", "123456").Return(nil)
	r := newItemTestRouter(svc, models.RoleLogistics)

	req := httptest.NewRequest(http.MethodDelete, "/items/123456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
