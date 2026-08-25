package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Leander-Wendt/Penates/backend/internal/handler"
	"github.com/Leander-Wendt/Penates/backend/internal/service"
	"github.com/Leander-Wendt/Penates/backend/internal/testutil/mocks"
)

func newAuthTestRouter(svc *mocks.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewAuthHandler(svc)
	r.POST("/auth/login", h.Login)
	return r
}

func TestAuthHandler_Login_Success(t *testing.T) {
	svc := new(mocks.AuthService)
	expiresAt := time.Now().Add(time.Hour)
	svc.On("Login", "user@example.com", "password123").Return("signed-token", expiresAt, nil)
	r := newAuthTestRouter(svc)

	payload, _ := json.Marshal(map[string]string{"email": "user@example.com", "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require := assert.New(t)
	require.NoError(json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal("signed-token", body["token"])
	require.Contains(body, "expiresAt")
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	svc := new(mocks.AuthService)
	svc.On("Login", "user@example.com", "wrong").Return("", time.Time{}, service.ErrInvalidCredentials)
	r := newAuthTestRouter(svc)

	payload, _ := json.Marshal(map[string]string{"email": "user@example.com", "password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Header().Get("WWW-Authenticate"))
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	svc := new(mocks.AuthService)
	r := newAuthTestRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Login")
}
