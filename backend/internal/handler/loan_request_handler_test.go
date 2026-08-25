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
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func newLoanRequestTestRouter(svc *mocks.LoanRequestService, userID uint, role models.Role) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	})
	h := handler.NewLoanRequestHandler(svc)
	r.GET("/loan-requests", h.List)
	r.POST("/loan-requests", h.Create)
	r.GET("/loan-requests/:id", h.Get)
	r.PUT("/loan-requests/:id", h.Update)
	r.DELETE("/loan-requests/:id", h.Delete)
	r.PATCH("/loan-requests/:id/status", middleware.RequireRoles(models.RoleAdmin, models.RoleLogistics), h.UpdateStatus)
	return r
}

func TestLoanRequestHandler_List(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("ListForUser", uint(5), models.RoleStudent).Return([]models.LoanRequest{{ID: 1, RequestingUserID: 5}}, nil)
	r := newLoanRequestTestRouter(svc, 5, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/loan-requests", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoanRequestHandler_Create(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("Create", uint(5), []string{"123456"}, "Room A").Return(&models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending"}, nil)
	r := newLoanRequestTestRouter(svc, 5, models.RoleStudent)

	payload, _ := json.Marshal(map[string]interface{}{
		"itemInventoryNumbers": []string{"123456"}, "locationOfItems": "Room A",
	})
	req := httptest.NewRequest(http.MethodPost, "/loan-requests", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLoanRequestHandler_Get_OwnerAllowed(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("GetByID", uint(1)).Return(&models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending"}, nil)
	r := newLoanRequestTestRouter(svc, 5, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/loan-requests/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoanRequestHandler_Get_NonOwnerForbidden(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("GetByID", uint(1)).Return(&models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending"}, nil)
	r := newLoanRequestTestRouter(svc, 6, models.RoleStudent)

	req := httptest.NewRequest(http.MethodGet, "/loan-requests/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLoanRequestHandler_Get_AdminAllowed(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("GetByID", uint(1)).Return(&models.LoanRequest{ID: 1, RequestingUserID: 5, Status: "pending"}, nil)
	r := newLoanRequestTestRouter(svc, 999, models.RoleAdmin)

	req := httptest.NewRequest(http.MethodGet, "/loan-requests/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoanRequestHandler_Update_ForbiddenBySerivce(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("Update", uint(1), uint(6), models.RoleStudent, "New Loc").Return(nil, service.ErrForbidden)
	r := newLoanRequestTestRouter(svc, 6, models.RoleStudent)

	payload, _ := json.Marshal(map[string]string{"locationOfItems": "New Loc"})
	req := httptest.NewRequest(http.MethodPut, "/loan-requests/1", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestLoanRequestHandler_Delete(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("Delete", uint(1), uint(5), models.RoleStudent).Return(nil)
	r := newLoanRequestTestRouter(svc, 5, models.RoleStudent)

	req := httptest.NewRequest(http.MethodDelete, "/loan-requests/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestLoanRequestHandler_UpdateStatus_AdminAllowed(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("UpdateStatus", uint(1), "approved").Return(&models.LoanRequest{ID: 1, Status: "approved"}, nil)
	r := newLoanRequestTestRouter(svc, 1, models.RoleAdmin)

	payload, _ := json.Marshal(map[string]string{"status": "approved"})
	req := httptest.NewRequest(http.MethodPatch, "/loan-requests/1/status", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoanRequestHandler_UpdateStatus_StudentForbidden(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	r := newLoanRequestTestRouter(svc, 1, models.RoleStudent)

	payload, _ := json.Marshal(map[string]string{"status": "approved"})
	req := httptest.NewRequest(http.MethodPatch, "/loan-requests/1/status", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	svc.AssertNotCalled(t, "UpdateStatus")
}

func TestLoanRequestHandler_UpdateStatus_InvalidStatus(t *testing.T) {
	svc := new(mocks.LoanRequestService)
	svc.On("UpdateStatus", uint(1), "bogus").Return(nil, service.ErrValidation)
	r := newLoanRequestTestRouter(svc, 1, models.RoleAdmin)

	payload, _ := json.Marshal(map[string]string{"status": "bogus"})
	req := httptest.NewRequest(http.MethodPatch, "/loan-requests/1/status", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
