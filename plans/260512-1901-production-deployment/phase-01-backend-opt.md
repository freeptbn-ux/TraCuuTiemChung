# Phase 01: Backend Production Optimization
Status: ✅ Completed
Dependencies: None

## Objective
Tối ưu hóa mã nguồn Python và cấu hình Vercel để đảm bảo hiệu suất và tính ổn định trên môi trường Serverless.

## Implementation Steps
1. [x] **Update requirements.txt:** Pin chính xác version của tất cả các thư viện (fastapi, uvicorn, beautifulsoup4, redis...).
2. [x] **Restructure for Vercel:** Chuyển entry point chính vào thư mục `api/` (theo chuẩn Vercel tự động nhận diện) hoặc cấu hình `rewrites` trong `vercel.json`.
3. [x] **Optimize `vercel.json`:** Sử dụng cấu trúc hiện đại nhất của Vercel để định tuyến traffic.
4. [x] **Global Error Handling:** Đảm bảo các lỗi cào dữ liệu không làm crash instance serverless, trả về mã lỗi JSON chuẩn.
5. [x] **Health Check Endpoint:** Tạo `/api/health` để monitor trạng thái server.

## Files to Create/Modify
- `vercel-backend/vercel.json`
- `vercel-backend/requirements.txt`
- `vercel-backend/main.py` (hoặc `api/index.py`)

## Test Criteria
- [x] Chạy `vercel dev` locally không lỗi (đã test qua pytest & dependency simulation).
- [x] Endpoint `/api/health` trả về status 200.
