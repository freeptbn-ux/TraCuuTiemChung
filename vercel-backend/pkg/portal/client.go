package portal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"vercel-backend/pkg/models"
)

const (
	LoginURL  = "https://tiemchung.vncdc.gov.vn/Account/Login"
	IndexURL  = "https://tiemchung.vncdc.gov.vn/TiemChung/DoiTuong/Index"
	SearchURL = "https://tiemchung.vncdc.gov.vn/TiemChung/DoiTuong/TimKiem"
	DetailURL = "https://tiemchung.vncdc.gov.vn/TiemChung/DoiTuong/Detail"
)

type PortalClient struct {
	restClient *resty.Client
	username   string
	password   string
	redis      *redis.Client
	locker     Locker

	LoginURL  string
	IndexURL  string
	SearchURL string
	DetailURL string
}

// NewPortalClient creates a new portal client with optional Redis support
func NewPortalClient(username, password string, redisClient *redis.Client) *PortalClient {
	client := resty.New()
	
	if redisClient != nil {
		jar := NewRedisCookieJar(redisClient, username)
		client.SetCookieJar(jar)
	} else {
		jar, _ := cookiejar.New(nil)
		client.SetCookieJar(jar)
	}
	client.SetTimeout(30 * time.Second)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0")

	return &PortalClient{
		restClient: client,
		username:   username,
		password:   password,
		redis:      redisClient,
		locker:     NewRedisLocker(redisClient),
		LoginURL:   LoginURL,
		IndexURL:   IndexURL,
		SearchURL:  SearchURL,
		DetailURL:  DetailURL,
	}
}

// RedisClient returns the underlying redis client
func (pc *PortalClient) RedisClient() *redis.Client {
	return pc.redis
}

// Login performs login with distributed lock to avoid race conditions
func (pc *PortalClient) Login() error {
	if pc.redis == nil {
		return pc.performLogin()
	}

	lockKey := fmt.Sprintf("portal:lock:login:%s", pc.username)
	ctx := context.Background()

	// Wait and retry logic
	maxRetries := 20 // 20 * 500ms = 10s
	for i := 0; i < maxRetries; i++ {
		ok, value, err := pc.locker.Acquire(ctx, lockKey, 15*time.Second)
		if err != nil {
			return err
		}

		if ok {
			slog.Info("Acquiring login lock...", "key", lockKey)
			defer func() {
				if err := pc.locker.Release(ctx, lockKey, value); err != nil {
					slog.Error("Failed to release lock", "key", lockKey, "error", err)
				}
			}()

			// Double check if already logged in by another worker
			if pc.IsLoggedIn() {
				slog.Info("Already logged in by another worker, skipping login", "key", lockKey)
				return nil
			}

			return pc.performLogin()
		}

		// Lock busy, wait and check if someone else finished
		slog.Debug("Login lock busy, waiting...", "key", lockKey, "attempt", i+1)
		time.Sleep(500 * time.Millisecond)

		if pc.IsLoggedIn() {
			slog.Info("Detected successful login from another worker while waiting", "key", lockKey)
			return nil
		}
	}

	return fmt.Errorf("timeout waiting for login lock after 10s")
}

// IsLoggedIn checks if the current session is valid by checking the cookie jar
func (pc *PortalClient) IsLoggedIn() bool {
	if pc.restClient.GetClient().Jar == nil {
		return false
	}
	u, _ := url.Parse(pc.LoginURL)
	cookies := pc.restClient.GetClient().Jar.Cookies(u)
	for _, c := range cookies {
		if c.Name == ".ASPXAUTH" {
			return true
		}
	}
	return false
}

// ConnectivityStatus đại diện cho kết quả kiểm tra kết nối
type ConnectivityStatus struct {
	URL          string `json:"url"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int64  `json:"response_time_ms"`
	IsBlocked    bool   `json:"is_blocked"`
	HasToken     bool   `json:"has_token"`
	Error        string `json:"error,omitempty"`
	BodySnippet  string `json:"body_snippet,omitempty"`
	LookupResult string `json:"lookup_result,omitempty"` // Kết quả tra cứu thử nghiệm
	SearchBody   string `json:"search_body,omitempty"`
}

// CheckPortalConnectivity kiểm tra xem portal có thể truy cập được không và có dấu hiệu bị chặn IP không
func (pc *PortalClient) CheckPortalConnectivity() ConnectivityStatus {
	start := time.Now()
	resp, err := pc.restClient.R().Get(pc.LoginURL)
	duration := time.Since(start).Milliseconds()

	status := ConnectivityStatus{
		URL:          pc.LoginURL,
		ResponseTime: duration,
	}

	if err != nil {
		status.Error = err.Error()
		status.IsBlocked = true
		return status
	}

	status.StatusCode = resp.StatusCode()
	body := resp.String()
	
	if len(body) > 200 {
		status.BodySnippet = body[:200]
	} else {
		status.BodySnippet = body
	}

	if strings.Contains(body, "__RequestVerificationToken") {
		status.HasToken = true
	} else if resp.StatusCode() == 200 {
		status.IsBlocked = true
		status.Error = "Response does not look like login page"
	}

	// Thử tra cứu một số điện thoại mẫu để xem phản hồi thực tế
	if status.HasToken {
		patients, err := pc.LookupPatients("0388634123")
		if err != nil {
			status.LookupResult = fmt.Sprintf("Lookup Error: %v", err)
		} else {
			status.LookupResult = fmt.Sprintf("Found %d patients", len(patients))
			// If 0 patients, capture snippet from internal state or by re-running
			// Actually, let's just re-run and capture snippet if 0
			if len(patients) == 0 {
				searchResp, _ := pc.restClient.R().
					SetQueryParam("SoDienThoai", "0388634123").
					SetQueryParam("Length", "5").
					SetHeader("X-Requested-With", "XMLHttpRequest").
					SetHeader("Referer", pc.IndexURL).
					Get(pc.SearchURL)
				status.SearchBody = searchResp.String()
			}
		}
	}

	return status
}



// performLogin performs the actual login to the portal
func (pc *PortalClient) performLogin() error {
	if pc.username == "" || pc.password == "" {
		return fmt.Errorf("PORTAL_USERNAME or PORTAL_PASSWORD not set")
	}

	// 1. Get login page to extract __RequestVerificationToken
	resp, err := pc.restClient.R().Get(pc.LoginURL)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return err
	}

	token, exists := doc.Find("input[name='__RequestVerificationToken']").Attr("value")
	if !exists {
		return fmt.Errorf("could not find __RequestVerificationToken")
	}

	// 2. Perform login
	resp, err = pc.restClient.R().
		SetMultipartFormData(map[string]string{
			"__RequestVerificationToken": token,
			"UserName":                   pc.username,
			"password":                   pc.password,
			"remember_me":                "true",
		}).
		SetHeader("Referer", pc.LoginURL).
		Post(pc.LoginURL)

	if err != nil {
		return err
	}

	// Check if login was successful (check for .ASPXAUTH cookie)
	found := false
	if pc.restClient.GetClient().Jar != nil {
		for _, cookie := range pc.restClient.GetClient().Jar.Cookies(resp.RawResponse.Request.URL) {
			if cookie.Name == ".ASPXAUTH" {
				found = true
				break
			}
		}
	}

	if !found {
		// Sometimes the redirect happens, check the response URL or body
		if !strings.Contains(resp.String(), "Tài khoản hoặc mật khẩu không đúng") && resp.StatusCode() == 200 {
			// If we are at Index page, we are logged in
			if strings.Contains(resp.Request.URL, "/TiemChung/DoiTuong/Index") {
				found = true
			}
		}
	}

	if !found {
		return fmt.Errorf("login failed: .ASPXAUTH cookie not found or redirected to login")
	}

	return nil
}

// LookupPatients searches for patients by phone number
func (pc *PortalClient) LookupPatients(phone string) ([]models.Patient, error) {
	// Always login first — on Vercel stateless serverless, there is no persistent
	// session between invocations. Relying on session detection is unreliable.
	if !pc.IsLoggedIn() {
		slog.Info("No active session, logging in before lookup", "phone", phone)
		if err := pc.Login(); err != nil {
			return nil, fmt.Errorf("login failed before lookup: %w", err)
		}
	}

	searchParams := map[string]string{
		"Length":            "5",
		"LoaiDiaChi":        "0",
		"VungMienId":        "-Khu vực-",
		"ThonApId":          "-Thôn/Ấp-",
		"NgaySinhTu":        "",
		"NgaySinhToi":       "",
		"GioiTinh":          "-1",
		"LuaTuoi":           "-1",
		"MaDoiTuong":        "",
		"TenDoiTuong":       "",
		"TenMe":             "",
		"TenBo":             "",
		"MaDinhDanh":        "",
		"SoDienThoai":       phone,
		"TenNguoiGiamHo":    "",
		"TinhTrangTheoDoi":  "-1",
		"TinhTrangMangThai": "-1",
		"PageNumber":        "1",
		"PageSize":          "20",
		"CurrentSystemDate": "",
	}

	req := pc.restClient.R().
		SetQueryParams(searchParams).
		SetQueryParam("X-Requested-With", "XMLHttpRequest").
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetHeader("Referer", pc.IndexURL).
		SetHeader("Accept", "text/html, */*").
		SetHeader("Accept-Language", "vi,en;q=0.9")

	resp, err := req.Get(pc.SearchURL)
	if err != nil {
		return nil, err
	}

	htmlContent := resp.String()
	
	// Detect session expiry or server errors
	isErrorPage := strings.Contains(htmlContent, "an error occurred") || 
	               (strings.Contains(htmlContent, "UserName") && strings.Contains(htmlContent, "__RequestVerificationToken"))
	
	if isErrorPage {
		slog.Warn("Session expired or error page detected, re-logging in", "phone", phone)
		if err := pc.Login(); err != nil {
			return nil, fmt.Errorf("re-login failed: %w", err)
		}
		// Retry with same headers and params
		resp, err = pc.restClient.R().
			SetQueryParams(searchParams).
			SetQueryParam("X-Requested-With", "XMLHttpRequest").
			SetHeader("X-Requested-With", "XMLHttpRequest").
			SetHeader("Referer", pc.IndexURL).
			SetHeader("Accept", "text/html, */*").
			SetHeader("Accept-Language", "vi,en;q=0.9").
			Get(pc.SearchURL)
		if err != nil {
			return nil, err
		}
		htmlContent = resp.String()
	}

	if !strings.Contains(htmlContent, "doiTuongSearchResult") && !strings.Contains(htmlContent, "không có đối tượng nào") {
		snippet := htmlContent
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		slog.Warn("LookupPatients: Results table not found", "phone", phone, "snippet", snippet)
	}

	return pc.ParseSearchResults(htmlContent)
}

// ParseSearchResults parses the HTML content of search results
func (pc *PortalClient) ParseSearchResults(htmlContent string) ([]models.Patient, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	// Regex for MA_DOI_TUONG
	re := regexp.MustCompile(`"MA_DOI_TUONG":"(\d+)"`)
	matches := re.FindAllStringSubmatch(htmlContent, -1)
	patientCodes := make([]string, len(matches))
	for i, match := range matches {
		patientCodes[i] = match[1]
	}

	results := []models.Patient{}
	doc.Find("table#doiTuongSearchResult tbody tr").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(strings.ToLower(s.Text()), "không có đối tượng nào") {
			return
		}

		cells := s.Find("td")
		if cells.Length() < 5 {
			return
		}

		idValue, _ := s.Attr("data-id")
		if idValue == "" {
			onclick, _ := s.Attr("onclick")
			reID := regexp.MustCompile(`OnShowDetail\((\d+)\)`)
			match := reID.FindStringSubmatch(onclick)
			if len(match) > 1 {
				idValue = match[1]
			}
		}

		if strings.Contains(idValue, ",") {
			idValue = strings.Split(idValue, ",")[0]
			idValue = strings.TrimSpace(idValue)
		}

		name := strings.TrimSpace(cells.Eq(1).Text())
		dob := strings.TrimSpace(cells.Eq(4).Text())
		gender := "N/A"
		if cells.Length() > 5 {
			gender = strings.TrimSpace(cells.Eq(5).Text())
		}

		code := ""
		if i < len(patientCodes) {
			code = patientCodes[i]
		}

		results = append(results, models.Patient{
			ID:     idValue,
			Name:   name,
			DOB:    dob,
			Gender: gender,
			Code:   code,
		})
	})

	return results, nil
}

// GetVaccinationHistory retrieves history for a specific patient
func (pc *PortalClient) GetVaccinationHistory(patientID string) (*models.PatientDetail, error) {
	// Always login first — on Vercel stateless serverless, there is no persistent session.
	if !pc.IsLoggedIn() {
		slog.Info("No active session, logging in before GetVaccinationHistory", "patientID", patientID)
		if err := pc.Login(); err != nil {
			return nil, fmt.Errorf("login failed before history fetch: %w", err)
		}
	}

	resp, err := pc.restClient.R().
		SetQueryParam("doiTuongId", patientID).
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetHeader("Referer", pc.IndexURL).
		Get(pc.DetailURL)

	if err != nil {
		return nil, err
	}

	htmlContent := resp.String()
	if strings.Contains(htmlContent, "UserName") && strings.Contains(htmlContent, "__RequestVerificationToken") {
		slog.Warn("Session expired mid-request, re-logging in", "patientID", patientID)
		if err := pc.Login(); err != nil {
			return nil, fmt.Errorf("re-login failed: %w", err)
		}
		resp, err = pc.restClient.R().
			SetQueryParam("doiTuongId", patientID).
			SetHeader("X-Requested-With", "XMLHttpRequest").
			SetHeader("Referer", pc.IndexURL).
			Get(pc.DetailURL)
		if err != nil {
			return nil, err
		}
		htmlContent = resp.String()
	}

	return pc.ParsePatientDetail(htmlContent)
}

// ParsePatientDetail parses the HTML content of the patient detail page
func (pc *PortalClient) ParsePatientDetail(htmlContent string) (*models.PatientDetail, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	info := models.PatientInfo{}
	
	// Name
	nameInput := doc.Find("#txtHoTen")
	if val, exists := nameInput.Attr("value"); exists {
		info.Name = strings.TrimSpace(val)
	} else {
		info.Name = strings.TrimSpace(nameInput.Text())
	}

	// DOB
	dobInput := doc.Find("#txtNgaySinh")
	if dobInput.Length() == 0 {
		dobInput = doc.Find("#hfNgaySinhDoiTuong")
	}
	if val, exists := dobInput.Attr("value"); exists {
		info.Birth = strings.TrimSpace(val)
	} else {
		info.Birth = strings.TrimSpace(dobInput.Text())
	}

	// System Date
	sysDateInput := doc.Find("#CurrentSystemDate")
	if sysDateInput.Length() == 0 {
		sysDateInput = doc.Find("#hfNgayHienTai")
	}
	if val, exists := sysDateInput.Attr("value"); exists {
		info.SystemDate = strings.TrimSpace(val)
	} else {
		info.SystemDate = strings.TrimSpace(sysDateInput.Text())
	}
	
	if info.SystemDate == "" {
		loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			loc = time.FixedZone("GMT+7", 7*3600)
		}
		info.SystemDate = time.Now().In(loc).Format("02/01/2006")
	}

	var history []models.VaccineRecord
	doc.Find("#tblVacxin tr").Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() >= 5 {
			cell1 := s.Find("td").Eq(1).Clone()
			cell1.Find("*").Remove()
			vaccineNameRaw := strings.TrimSpace(cell1.Text())
			if vaccineNameRaw == "" {
				vaccineNameRaw = strings.TrimSpace(cells.Eq(1).Text())
			}

			doseTextRaw := strings.TrimSpace(cells.Eq(2).Text())
			dateTextRaw := strings.TrimSpace(cells.Eq(4).Text())

			if vaccineNameRaw != "" && dateTextRaw != "" {
				dateTextRaw = strings.ReplaceAll(dateTextRaw, " ", "")
				dateObj, err := time.Parse("02/01/2006", dateTextRaw)
				if err == nil {
					history = append(history, models.VaccineRecord{
						VaccineName: vaccineNameRaw,
						Dose:        doseTextRaw,
						Date:        dateObj,
					})
				}
			}
		}
	})

	return &models.PatientDetail{
		PatientInfo: info,
		History:     history,
	}, nil
}
