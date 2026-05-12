# Phase 02: Vercel Deployment & Secret Setup
Status: ✅ Completed
Dependencies: Phase 01

## Objective
Đưa code lên Cloud và cấu hình các biến môi trường nhạy cảm một cách an toàn.

## Implementation Steps
1. [ ] **Git Push:** Đảm bảo toàn bộ code backend đã được push lên branch `main`.
2. [ ] **Create Vercel Project:** Kết nối GitHub repo với Vercel.
3. [ ] **Environment Variables Setup:** Thêm các key sau vào Vercel Settings:
    - `UPSTASH_REDIS_REST_URL`
    - `UPSTASH_REDIS_REST_TOKEN`
    - `X_API_KEY` (Key tự định nghĩa để app Android gọi API)
4. [ ] **Trigger Deploy:** Chạy deployment và lấy URL chính thức (ví dụ: `tracuutiemchung.vercel.app`).
5. [ ] **Logs Monitoring:** Kiểm tra logs trên Vercel để đảm bảo Redis kết nối thành công.

## Test Criteria
- [ ] URL chính thức truy cập được.
- [ ] API trả về kết quả đúng khi test bằng Postman/Curl (kèm X-API-KEY).
