package portal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestParseSearchResults(t *testing.T) {
	html := `
	<table id="doiTuongSearchResult">
		<tbody>
			<tr data-id="12345">
				<td>1</td>
				<td>Nguyen Van A</td>
				<td></td>
				<td></td>
				<td>01/01/2020</td>
				<td>Nam</td>
			</tr>
		</tbody>
	</table>
	{"MA_DOI_TUONG":"98765"}
	`

	pc := NewPortalClient("user", "pass", nil)
	results, err := pc.ParseSearchResults(html)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	r := results[0]
	if r.ID != "12345" {
		t.Errorf("Expected ID 12345, got %s", r.ID)
	}
	if r.Name != "Nguyen Van A" {
		t.Errorf("Expected Name Nguyen Van A, got %s", r.Name)
	}
	if r.DOB != "01/01/2020" {
		t.Errorf("Expected DOB 01/01/2020, got %s", r.DOB)
	}
	if r.Gender != "Nam" {
		t.Errorf("Expected Gender Nam, got %s", r.Gender)
	}
	if r.Code != "98765" {
		t.Errorf("Expected Code 98765, got %s", r.Code)
	}
}

func TestParsePatientDetail(t *testing.T) {
	html := `
	<input id="txtHoTen" value="Nguyen Van B">
	<input id="txtNgaySinh" value="02/02/2022">
	<input id="CurrentSystemDate" value="14/05/2026">
	<table id="tblVacxin">
		<tr><td>STT</td><td>Tên vắc xin</td><td>Mũi</td><td>Ngày nhắc</td><td>Ngày tiêm</td></tr>
		<tr>
			<td>1</td>
			<td>Vắc xin A</td>
			<td>1</td>
			<td></td>
			<td>10/05/2026</td>
		</tr>
		<tr>
			<td>2</td>
			<td>Vắc xin B <span class="sublabel">(Ghi chú)</span></td>
			<td>2</td>
			<td></td>
			<td>12/05/2026</td>
		</tr>
	</table>
	`

	pc := NewPortalClient("user", "pass", nil)
	detail, err := pc.ParsePatientDetail(html)
	if err != nil {
		t.Fatalf("Failed to parse detail: %v", err)
	}

	if detail.PatientInfo.Name != "Nguyen Van B" {
		t.Errorf("Expected Name Nguyen Van B, got %s", detail.PatientInfo.Name)
	}
	if detail.PatientInfo.Birth != "02/02/2022" {
		t.Errorf("Expected Birth 02/02/2022, got %s", detail.PatientInfo.Birth)
	}
	if detail.PatientInfo.SystemDate != "14/05/2026" {
		t.Errorf("Expected SystemDate 14/05/2026, got %s", detail.PatientInfo.SystemDate)
	}

	if len(detail.History) != 2 {
		t.Fatalf("Expected 2 records, got %d", len(detail.History))
	}

	rec1 := detail.History[0]
	if rec1.VaccineName != "Vắc xin A" {
		t.Errorf("Expected VaccineName Vắc xin A, got %s", rec1.VaccineName)
	}
	if rec1.Dose != "1" {
		t.Errorf("Expected Dose 1, got %s", rec1.Dose)
	}
	expectedDate1 := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	if !rec1.Date.Equal(expectedDate1) {
		t.Errorf("Expected Date %v, got %v", expectedDate1, rec1.Date)
	}

	rec2 := detail.History[1]
	if rec2.VaccineName != "Vắc xin B" {
		t.Errorf("Expected VaccineName Vắc xin B, got %s", rec2.VaccineName)
	}
}

func TestRetryLoginMechanism(t *testing.T) {
	loginCallCount := 0
	searchCallCount := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login":
			if r.Method == http.MethodGet {
				fmt.Fprint(w, `<input name="__RequestVerificationToken" value="fake-token">`)
			} else if r.Method == http.MethodPost {
				loginCallCount++
				http.SetCookie(w, &http.Cookie{Name: ".ASPXAUTH", Value: "fake-auth-cookie"})
				w.WriteHeader(http.StatusOK)
			}
		case "/search":
			searchCallCount++
			if searchCallCount == 1 {
				// Simulate session expired: returns login page elements
				fmt.Fprint(w, `<html><input name="UserName"><input name="__RequestVerificationToken"></html>`)
			} else {
				// Simulate success
				fmt.Fprint(w, `<table id="doiTuongSearchResult"><tbody><tr><td>1</td><td>Success</td><td></td><td></td><td>01/01/2000</td></tr></tbody></table>`)
			}
		}
	}))
	defer ts.Close()

	pc := NewPortalClient("user", "pass", nil)
	pc.LoginURL = ts.URL + "/login"
	pc.SearchURL = ts.URL + "/search"
	pc.IndexURL = ts.URL + "/index"

	results, err := pc.LookupPatients("0123456789")
	if err != nil {
		t.Fatalf("LookupPatients failed: %v", err)
	}

	if loginCallCount != 1 {
		t.Errorf("Expected 1 login call, got %d", loginCallCount)
	}
	if searchCallCount != 2 {
		t.Errorf("Expected 2 search calls, got %d", searchCallCount)
	}
	if len(results) != 1 || results[0].Name != "Success" {
		t.Errorf("Unexpected results: %+v", results)
	}
}
