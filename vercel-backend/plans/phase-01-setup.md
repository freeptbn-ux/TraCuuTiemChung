# Phase 01: Setup Environment

## Objective
Khởi tạo nền tảng Go và cấu trúc thư mục chuẩn.

## Tasks:
- [x] Khởi tạo module: `go mod init vercel-backend`
- [x] Cài đặt core dependencies: `gin`, `resty`, `goquery`, `godotenv`.
- [x] Tạo file `internal/config/config.go` để quản lý Env.
- [x] Setup folder structure như `plan.md`.

## 🧪 Testing for this phase:
- **File:** `internal/config/config_test.go`
- **Mục tiêu:** Kiểm tra xem app có đọc đúng biến môi trường `X_API_KEY` và `PORT` hay không.
- **Lệnh chạy:** `go test -v ./internal/config/...`
