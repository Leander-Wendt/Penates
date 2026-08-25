package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Leander-Wendt/Penates/backend/internal/handler"
	"github.com/Leander-Wendt/Penates/backend/internal/middleware"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func newUserTestRouter(svc *mocks.UserService, userID uint, role models.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	})
	h := handler.NewUserHandler(svc)
	r.GET("/users/me", h.Me)
	r.GET("/users", middleware.RequireRoles(models.RoleAdmin, models.RoleLogistics), h.List)
	r.POST("/users", middleware.RequireRoles(models.RoleAdmin, models.RoleLogistics), h.Create)
	r.GET("/users/:id", h.Get)
	r.PUT("/users/:id", middleware.RequireRoles(models.RoleAdmin, models.RoleLogistics), h.Update)
	r.DELETE("/users/:id", middleware.RequireRoles(models.RoleAdmin), h.Delete)
	return r
}

func TestUserHandler_Me(t *testing.T) {
	svc := new(mocks.UserService)
	svc.On("GetByID", uint(7)).Return(&models.User{ID: 7, Email: "me@example.com"}, nil)
	r := newUserTestRouter(svc, 7, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_Get_OwnRecordAllowedForStudent(t *testing.T) {
	svc := new(mocks.UserService)
	svc.On("GetByID", uint(7)).Return(&models.User{ID: 7, Email: "me@example.com"}, nil)
	r := newUserTestRouter(svc, 7, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/users/7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_Get_OtherRecordForbiddenForStudent(t *testing.T) {
	svc := new(mocks.UserService)
	r := newUserTestRouter(svc, 7, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/users/8", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	svc.AssertNotCalled(t, "GetByID")
}

func TestUserHandler_Get_AdminCanGetAnyRecord(t *testing.T) {
	svc := new(mocks.UserService)
	svc.On("GetByID", uint(8)).Return(&models.User{ID: 8, Email: "other@example.com"}, nil)
	r := newUserTestRouter(svc, 7, models.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/users/8", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserHandler_Delete_AdminAllowed(t *testing.T) {
	svc := new(mocks.UserService)
	svc.On("Delete", uint(8)).Return(nil)
	r := newUserTestRouter(svc, 7, models.RoleAdmin)

	req := httptest.NewRequest(http.MethodDelete, "/users/8", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestUserHandler_Delete_LogisticsForbidden(t *testing.T) {
	svc := new(mocks.UserService)
	r := newUserTestRouter(svc, 7, models.RoleLogistics)

	req := httptest.NewRequest(http.MethodDelete, "/users/8", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("WWW-Authenticate"))
	svc.AssertNotCalled(t, "Delete")
}

func TestUserHandler_Create(t *testing.T) {
	svc := new(mocks.UserService)
	svc.On("Create", "new@example.com", "password123", models.RoleStudent, uint(1)).
		Return(&models.User{ID: 9, Email: "new@example.com", Role: models.RoleStudent, OrganisationID: 1}, nil)
	r := newUserTestRouter(svc, 1, models.RoleAdmin)

	payload, _ := json.Marshal(map[string]interface{}{
		"email": "new@example.com", "password": "password123", "role": "student", "organisationId": 1,
	})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserHandler_List(t *testing.T) {
	svc := new(mocks.UserService)
	svc.On("List").Return([]models.User{{ID: 1}, {ID: 2}}, nil)
	r := newUserTestRouter(svc, 1, models.RoleLogistics)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
