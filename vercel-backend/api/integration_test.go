package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"vercel-backend/pkg/config"
	"vercel-backend/pkg/portal"
)

func TestFullIntegration(t *testing.T) {
	// 1. Setup Mock Portal Server
	mockPortal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/Account/Login":
			if r.Method == "GET" {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(`<html><input name="__RequestVerificationToken" value="mock-token"/></html>`))
			} else {
				// Login POST
				http.SetCookie(w, &http.Cookie{Name: ".ASPXAUTH", Value: "mock-auth-cookie", Path: "/"})
				w.WriteHeader(http.StatusOK)
			}
		case "/TiemChung/DoiTuong/TimKiem":
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
				<table id="doiTuongSearchResult">
					<tbody>
						<tr data-id="12345">
							<td>1</td>
							<td>Nguyễn Văn A</td>
							<td>Nam</td>
							<td>01/01/2020</td>
							<td>01/01/2020</td>
							<td>Nam</td>
						</tr>
					</tbody>
				</table>
				<script>var data = {"MA_DOI_TUONG":"98765"};</script>
			`))
		case "/TiemChung/DoiTuong/Detail":
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`
				<html>
					<input id="txtHoTen" value="Nguyễn Văn A"/>
					<input id="txtNgaySinh" value="01/01/2020"/>
					<input id="CurrentSystemDate" value="15/05/2026"/>
					<table id="tblVacxin">
						<tr>
							<td>1</td>
							<td>BCG</td>
							<td>Mũi 1</td>
							<td>VNVC</td>
							<td>05/01/2020</td>
						</tr>
					</table>
				</html>
			`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockPortal.Close()

	// 2. Override Config and PortalClient
	os.Setenv("X_API_KEY", "test-api-key")
	os.Setenv("PORTAL_USERNAME", "test-user")
	os.Setenv("PORTAL_PASSWORD", "test-pass")
	
	cfg = config.LoadConfig()
	pc = portal.NewPortalClient(cfg.PORTAL_USERNAME, cfg.PORTAL_PASSWORD, nil)
	pc.LoginURL = mockPortal.URL + "/Account/Login"
	pc.SearchURL = mockPortal.URL + "/TiemChung/DoiTuong/TimKiem"
	pc.DetailURL = mockPortal.URL + "/TiemChung/DoiTuong/Detail"
	pc.IndexURL = mockPortal.URL + "/TiemChung/DoiTuong/Index"

	// 3. Test /api/lookup
	t.Run("Lookup Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{"phone": "0388634123"})
		req, _ := http.NewRequest("POST", "/api/lookup", bytes.NewBuffer(reqBody))
		req.Header.Set("X-API-KEY", "test-api-key")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		Handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp StandardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		
		data, ok := resp.Data.([]interface{})
		assert.True(t, ok)
		assert.NotEmpty(t, data)
		
		patient := data[0].(map[string]interface{})
		assert.Equal(t, "Nguyễn Văn A", patient["name"])
		assert.Equal(t, "12345", patient["id"])
	})

	// 4. Test /api/analyze
	t.Run("Analyze Success", func(t *testing.T) {
		reqBody, _ := json.Marshal(map[string]string{"patient_id": "12345"})
		req, _ := http.NewRequest("POST", "/api/analyze", bytes.NewBuffer(reqBody))
		req.Header.Set("X-API-KEY", "test-api-key")
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		Handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp StandardResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		
		analysis := resp.Data.(map[string]interface{})
		assert.Equal(t, "Nguyễn Văn A", analysis["patient_name"])
		assert.NotEmpty(t, analysis["administered_vaccines"])
		assert.NotEmpty(t, analysis["missing_vaccines"])
	})

	// 5. Test Password Change (Criteria)
	t.Run("Handle Portal Auth Failure", func(t *testing.T) {
		// Modify mock to return 401 or redirect to login without cookie
		// Simulating password change by making login fail
		originalURL := pc.LoginURL
		pc.LoginURL = mockPortal.URL + "/WrongPath" 
		
		reqBody, _ := json.Marshal(map[string]string{"phone": "0388634123"})
		req, _ := http.NewRequest("POST", "/api/lookup", bytes.NewBuffer(reqBody))
		req.Header.Set("X-API-KEY", "test-api-key")

		w := httptest.NewRecorder()
		Handler(w, req)

		// It should return 500 or 401 depending on ErrorHandler, but not crash
		assert.True(t, w.Code >= 400)
		
		// Restore
		pc.LoginURL = originalURL
	})
}
