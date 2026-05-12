# Phase 01: Setup Environment & Vercel Config
Status: ✅ Completed
Dependencies: None

## Objective
Khởi tạo cấu trúc thư mục chuẩn cho backend FastAPI chạy trên Vercel. Thiết lập các thư viện cần thiết và cấu hình file `vercel.json` để tránh lỗi timeout.

## Requirements
### Functional
- [x] Khởi tạo thư mục `vercel-backend/` trong dự án.
- [x] Thiết lập môi trường ảo (virtual environment) Python nếu code local.
- [x] Cài đặt các thư viện: `fastapi`, `uvicorn`, `requests`, `beautifulsoup4`, `upstash-redis`.
- [x] Tạo file cấu hình `.env.example` chứa các key cần thiết (Upstash URL, Token, Portal Credentials).
- [x] Cấu hình `vercel.json`.

### Non-Functional
- [x] Đảm bảo `vercel.json` cấu hình `maxDuration` hợp lý (vd: 30s) để tránh lỗi Timeout khi scrape portal lần đầu.

## Implementation Steps
1. [x] Tạo cấu trúc thư mục:
   ```
   vercel-backend/
   ├── api/
   │   └── index.py
   ├── core/
   ├── requirements.txt
   ├── vercel.json
   └── .env.example
   ```
2. [x] Khai báo dependencies trong `requirements.txt`.
3. [x] Viết `vercel.json` chỉ định `api/index.py` là entry point cho Vercel serverless.
4. [x] Khởi tạo file `api/index.py` với 1 endpoint GET `/health` đơn giản để test.

## Files to Create/Modify
- `vercel-backend/requirements.txt`
- `vercel-backend/vercel.json`
- `vercel-backend/api/index.py`
- `vercel-backend/.env.example`

## Test Criteria
- [x] Chạy `pip install -r requirements.txt` thành công.
- [x] Chạy local bằng `uvicorn api.index:app --reload` và ping `/health` trả về 200 OK.

---
Next Phase: `phase-02-core-logic.md`
