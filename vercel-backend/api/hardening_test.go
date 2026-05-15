package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"vercel-backend/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestHardening(t *testing.T) {
	os.Setenv("X_API_KEY", "test_hardening_key")
	cfg = config.LoadConfig()

	t.Run("RequestID Middleware", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/health", nil)
		w := httptest.NewRecorder()
		Handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("X-Request-Id"))

		var resp StandardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, w.Header().Get("X-Request-Id"), resp.RequestID)
	})

	t.Run("Auth Middleware with RequestID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/lookup", bytes.NewBuffer([]byte(`{}`)))
		// No API Key
		w := httptest.NewRecorder()
		Handler(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		
		var resp StandardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.NotEmpty(t, resp.RequestID)
	})

	t.Run("Centralized Error Handling", func(t *testing.T) {
		// To trigger an error, we can send invalid JSON to a protected route
		req, _ := http.NewRequest("POST", "/api/lookup", bytes.NewBuffer([]byte(`invalid json`)))
		req.Header.Set("X-API-KEY", "test_hardening_key")
		w := httptest.NewRecorder()
		Handler(w, req)

		// Gin's ShouldBindJSON returns 400 if it fails, which we handle in handleLookup
		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var resp StandardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "error", resp.Status)
		assert.Contains(t, resp.Message, "Phone number is required")
		assert.NotEmpty(t, resp.RequestID)
	})

	t.Run("Rate Limiting", func(t *testing.T) {
		// Since we use nil Redis in test Handler (unless configured), 
		// it should pass through. To test 429, we'd need a real/mock redis.
		// However, we can verify that with nil Redis it doesn't crash.
		req, _ := http.NewRequest("GET", "/api/health", nil)
		w := httptest.NewRecorder()
		Handler(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
