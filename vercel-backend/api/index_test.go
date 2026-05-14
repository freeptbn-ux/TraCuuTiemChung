package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"vercel-backend/internal/config"
)

func TestAuthRequired(t *testing.T) {
	// Setup config for test
	os.Setenv("X_API_KEY", "test_secret")
	cfg = config.LoadConfig()

	tests := []struct {
		name           string
		apiKey         string
		expectedStatus int
	}{
		{
			name:           "Valid API Key",
			apiKey:         "test_secret",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid API Key",
			apiKey:         "wrong_secret",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Missing API Key",
			apiKey:         "",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/lookup", bytes.NewBuffer([]byte(`{"phone":"0123"}`)))
			if tt.apiKey != "" {
				req.Header.Set("X-API-KEY", tt.apiKey)
			}

			w := httptest.NewRecorder()
			Handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestHealthCheck(t *testing.T) {
	os.Setenv("X_API_KEY", "test_secret")
	cfg = config.LoadConfig()

	req, _ := http.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()
	Handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestLookupAuth(t *testing.T) {
	os.Setenv("X_API_KEY", "test_secret")
	cfg = config.LoadConfig()

	reqBody, _ := json.Marshal(map[string]string{"phone": "0123456789"})
	req, _ := http.NewRequest("POST", "/api/lookup", bytes.NewBuffer(reqBody))
	// No auth header

	w := httptest.NewRecorder()
	Handler(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
