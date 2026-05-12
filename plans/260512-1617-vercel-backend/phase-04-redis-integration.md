# Phase 04: Redis Session Integration (Upstash)
Status: ✅ Completed
Dependencies: phase-03

## Objective
Giải quyết vấn đề Timeout và độ trễ trên Vercel Serverless bằng cách cache lại thông tin đăng nhập Portal (Cookies + Verification Token) vào Upstash Redis.

## Requirements
### Functional
- [x] Sử dụng thư viện `upstash-redis` HTTP REST client (không dùng thư viện redis chuẩn để tránh lỗi TCP tcp_keepalive trên serverless).
- [x] Chuyển logic `LoginManager` cũ sang `AuthService`.
- [x] Flow: Khi API cần request tới CDC, gọi Redis kiểm tra session:
  - Nếu CÓ: Lấy ra dùng ngay.
  - Nếu KHÔNG (hoặc hết hạn, hoặc bị portal từ chối do session chết): Gọi `AuthService.login()`, lưu cookie/token mới vào Redis kèm TTL (vd: 30 phút), sau đó mới thực hiện request.

### Non-Functional
- [x] Code không block I/O lâu.
- [x] Tránh login đồng thời khi có nhiều request (có thể dùng cờ lock đơn giản hoặc kệ, login nhiều lần portal thường không cấm nhưng tốn tài nguyên).

## Implementation Steps
1. [x] Tạo db trên **Upstash Redis** (đăng nhập bằng tk Github/Google), lấy REST URL & Token điền vào `.env`.
2. [x] Tạo `vercel-backend/services/redis_cache.py`.
3. [x] Tạo `vercel-backend/services/auth_service.py` để xử lý login CDC portal.
4. [x] Tích hợp `AuthService` và `RedisCache` vào `PortalClient` ở Phase 03.


## Files to Create/Modify
- `vercel-backend/services/redis_cache.py`
- `vercel-backend/services/auth_service.py`
- `vercel-backend/services/portal_client.py`

## Test Criteria
- [ ] Test request đầu tiên: Log hiện "Đang login và lưu cache". Phản hồi mất 3-5s.
- [ ] Test request thứ hai: Log hiện "Dùng session từ Redis". Phản hồi mất < 1s.

---
Next Phase: `phase-05-android-integration.md`
