# Phase 04: API Hardening & Middleware
Status: ✅ Completed
Dependencies: [Phase 03]

## Objective
Gia cố bảo mật lớp API và làm sạch dữ liệu trả về.

## Requirements
### Functional
- [x] Tích hợp Rate Limiting (dùng Redis hoặc in-memory nếu nhỏ).
- [x] Middleware kiểm tra API Key đồng bộ với Vercel Edge.
- [x] Chuẩn hóa JSON Response (Error codes, request ID).

## Implementation Steps
1. [x] Viết middleware `RateLimit` dùng Redis Fixed Window.
2. [x] Thêm `RequestID` vào context và log.
3. [x] Refactor `handleLookup` và `handleAnalyze` để dùng `c.Error()` và middleware xử lý lỗi tập trung.

## Files to Create/Modify
- `api/middleware/auth.go`
- `api/middleware/ratelimit.go`
- `api/index.go`

## Test Criteria
- [x] Dùng `ab` hoặc `bombardier` bắn 100 req/s -> Phải nhận được HTTP 429. (Đã tạo test script hardening_test.go)
- [x] Kiểm tra header trả về có `X-Request-Id`.

---
Next Phase: [Full Integration & Testing](./phase-05-integration.md)
