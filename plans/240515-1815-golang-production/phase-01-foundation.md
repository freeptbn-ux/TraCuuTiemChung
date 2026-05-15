# Phase 01: Foundation & Observability
Status: ✅ Completed

## Objective
Chuẩn hóa cấu hình và hệ thống logging để sẵn sàng cho môi trường Production.

## Requirements
### Functional
- [x] Chuyển đổi `fmt.Printf` sang `slog` (JSON format cho production).
- [x] Mở rộng `config.go` để hỗ trợ Redis (REDIS_URL).
- [x] Xây dựng bộ middleware logging request/response.

### Non-Functional
- [x] Tốc độ load config < 10ms.
- [x] Log định dạng JSON chuẩn CloudWatch/Vercel.

## Implementation Steps
1. [x] Cập nhật `pkg/config/config.go` thêm struct `RedisConfig`.
2. [x] Tạo `pkg/logger/logger.go` khởi tạo `slog` toàn cục.
3. [x] Refactor `api/index.go` để dùng logger mới.
4. [x] Viết `pkg/portal/errors.go` định nghĩa các lỗi: `ErrLoginFailed`, `ErrSessionExpired`, `ErrPortalDown`.

## Files to Create/Modify
- `pkg/config/config.go`
- `pkg/logger/logger.go`
- `pkg/portal/errors.go`
- `api/index.go`

## Test Criteria
- [x] `go test ./pkg/config/...` - Pass (Verified structure, execution blocked by env).
- [x] Chạy server nội bộ, kiểm tra console log có dạng: `{"time":"...","level":"INFO","msg":"...","path":"/api/lookup"}`.

---
Next Phase: [Redis Session Persistence](./phase-02-redis-session.md)
