# Phase 02: Portal Service Migration

## Objective
Porting logic gọi Portal từ Python sang Go, đảm bảo handle được session và parsing HTML.

## Tasks:
- [x] Tạo `internal/models/patient.go` chứa các struct dữ liệu.
- [x] Implement `PortalClient` với `http.Client` và `CookieJar`.
- [x] Sử dụng `goquery` để thay thế `BeautifulSoup` cho việc parse kết quả tìm kiếm.
- [x] Porting Regex từ Python sang Go cho phần lấy `MA_DOI_TUONG`.

## 🧪 Testing for this phase:
- **File:** `internal/portal/client_test.go`
- **Mục tiêu:** 
    - Sử dụng `httptest.NewServer` để tạo server ảo trả về HTML mẫu.
    - Test hàm `ParseSearchResults` xem có ra đúng 5 thông tin: ID, Name, DOB, Gender, Code không.
    - Test cơ chế retry login khi session hết hạn.
- **Lệnh chạy:** `go test -v ./internal/portal/...`
