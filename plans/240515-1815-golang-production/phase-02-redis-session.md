# Phase 02: Redis Session Persistence
Status: ✅ Completed
Dependencies: [Phase 01]

## Objective
Hiện thực hóa việc lưu trữ cookie vào Redis để dùng chung giữa các Lambda instances.

## Requirements
### Functional
- [x] Implement `RedisCookieJar` tuân thủ interface `http.CookieJar`.
- [x] Tích hợp `RedisCookieJar` vào `PortalClient`.
- [x] Cơ chế `SyncToRedis` và `LoadFromRedis`.

### Non-Functional
- [x] Không làm tăng latency của API quá 50ms (Redis latency).
- [x] Tự động refresh TTL cho session trong Redis.

## Implementation Steps
1. [x] Install `github.com/redis/go-redis/v9`.
2. [x] Tạo `pkg/portal/cookie_jar.go` chứa logic Redis.
3. [x] Cập nhật `NewPortalClient` để nhận Redis client.
4. [x] Sửa logic `Login()` để sau khi login thành công thì lưu ngay vào Redis.
5. [x] Sửa logic trước khi request: Nếu jar trống, thử load từ Redis.

## Files to Create/Modify
- `pkg/portal/cookie_jar.go`
- `pkg/portal/client.go`
- `pkg/portal/client_test.go`

## Test Criteria
- [x] Unit test: Mock Redis, giả lập cookie và kiểm tra dữ liệu được ghi đúng key.
- [x] Integration test: Gọi API lookup 2 lần, lần 2 không được phép thấy log "Performing login...".

---
Next Phase: [Concurrency & Locking](./phase-03-concurrency.md)
