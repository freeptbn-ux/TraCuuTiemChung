package portal

import (
	"fmt"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
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

	LoginURL  string
	IndexURL  string
	SearchURL string
	DetailURL string
}

// NewPortalClient creates a new portal client with cookie jar support
func NewPortalClient(username, password string) *PortalClient {
	jar, _ := cookiejar.New(nil)
	client := resty.New()
	client.SetCookieJar(jar)
	client.SetTimeout(30 * time.Second)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0")

	return &PortalClient{
		restClient: client,
		username:   username,
		password:   password,
		LoginURL:   LoginURL,
		IndexURL:   IndexURL,
		SearchURL:  SearchURL,
		DetailURL:  DetailURL,
	}
}

// Login performs login to the portal and stores cookies in the jar
func (pc *PortalClient) Login() error {
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
			"remember_me":                "false",
		}).
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

	resp, err := pc.restClient.R().
		SetQueryParams(searchParams).
		SetHeader("X-Requested-With", "XMLHttpRequest").
		SetHeader("Referer", pc.IndexURL).
		Get(pc.SearchURL)

	if err != nil {
		return nil, err
	}

	htmlContent := resp.String()
	// Detect session expiry (usually redirects to Login or shows Login form in AJAX response)
	if strings.Contains(htmlContent, "UserName") && strings.Contains(htmlContent, "__RequestVerificationToken") {
		// Session expired, retry login
		if err := pc.Login(); err != nil {
			return nil, err
		}
		resp, err = pc.restClient.R().
			SetQueryParams(searchParams).
			SetHeader("X-Requested-With", "XMLHttpRequest").
			SetHeader("Referer", pc.IndexURL).
			Get(pc.SearchURL)
		if err != nil {
			return nil, err
		}
		htmlContent = resp.String()
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
		if err := pc.Login(); err != nil {
			return nil, err
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
