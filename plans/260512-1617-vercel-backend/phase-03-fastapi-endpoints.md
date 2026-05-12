# Phase 03: FastAPI Endpoints & Portal Client
Status: ✅ Completed
Dependencies: phase-02

## Objective
Xây dựng lớp giao tiếp (API Client) để nối với CDC Portal (từ `api_client.py` cũ) và bọc nó bằng các endpoint FastAPI để App Android có thể gọi.

## Requirements
### Functional
- [x] Implement `PortalClient`: Class chứa hàm HTTP POST/GET, tự động attach token và cookie.
- [x] Endpoint `POST /api/lookup`: Nhận JSON `{"phone": "098xxx"}` -> gọi portal -> gọi parser -> trả mảng JSON thông tin người dùng cơ bản (tên, ngày sinh, id nội bộ nếu có).
- [x] Endpoint `POST /api/analyze`: Nhận JSON `{"patient_id": "xxx", "phone": "098xxx"}` -> gọi portal lấy lịch sử -> gọi `AnalyzerService` -> trả về báo cáo phân tích.

### Non-Functional
- [x] Handle lỗi lịch sự: Trả về HTTP 400/500 kèm thông điệp tiếng Việt dễ hiểu nếu portal sập, sai SĐT, hoặc parse lỗi.
- [x] Add bảo mật cơ bản: Kiểm tra `X-API-KEY` header khớp với key quy định trước để tránh việc API bị public lạm dụng.

## Implementation Steps
1. [x] Tạo `vercel-backend/services/portal_client.py`.
2. [x] Sửa `vercel-backend/api/index.py`, thêm 2 endpoint `POST`.
3. [x] Định nghĩa Pydantic Models cho Request (vd: `class LookupRequest(BaseModel): phone: str`) và Response.
4. [x] Thêm logic middleware/dependency injection trong FastAPI để check `X-API-KEY`.

## Files to Create/Modify
- `vercel-backend/services/portal_client.py`
- `vercel-backend/api/index.py`

## Test Criteria
- [x] Gọi POST `/api/lookup` bằng Postman với SĐT thật (hoặc mock) trả về JSON chuẩn.
- [x] Nếu không truyền header `X-API-KEY`, API chặn (HTTP 401/403).

---
Next Phase: `phase-04-redis-integration.md`
