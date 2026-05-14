# Phase 04: API Handlers & Routing

## Objective
Kết nối các service lại và phơi ra ngoài API cho Frontend dùng.

## Tasks:
- [x] Viết `api/index.go` sử dụng Gin hoặc Echo.
- [x] Implement Middleware `AuthRequired` kiểm tra API Key.
- [x] Kết nối `PortalClient` và `AnalyzerService` qua Dependency Injection.

## 🧪 Testing for this phase:
- **File:** `api/index_test.go`
- **Mục tiêu:**
    - Sử dụng `net/http/httptest` để gọi thử vào endpoint `/api/lookup`.
    - Kiểm tra xem nếu không có API Key thì có trả về 403 không.
    - Kiểm tra cấu trúc JSON trả về có khớp 100% với Frontend (Android app) đang yêu cầu không.
- **Lệnh chạy:** `go test -v ./api/...`
