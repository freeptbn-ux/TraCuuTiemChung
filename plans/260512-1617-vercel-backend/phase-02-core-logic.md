# Phase 02: Porting Core Logic (Engine & Parser) - COMPLETED ✅
Status: ✅ Completed
Dependencies: phase-01

## Objective
Di chuyển logic "cốt lõi" từ mã nguồn Python cũ (`VaccineAnalyzer-Pro-main`) sang thư mục `core/` của backend mới. Đảm bảo logic dự báo (Analysis Engine) và logic bóc tách HTML (Parser) chạy độc lập, không dính dáng tới UI.

## Requirements
### Functional
- [x] Copy `html_parser.py` sang `vercel-backend/core/`.
- [x] Copy và chuẩn hóa toàn bộ file thuộc **Vaccine Analysis Engine** (`rule_processor.py`, `series_checkers.py`, `group_checkers_*.py`, v.v.).
- [x] Đảm bảo engine đọc file `vaccine_rules.json` (tái sử dụng bộ quy tắc đã làm cho bản Kotlin).
- [x] Đảm bảo không còn các thư viện UI (`PySide6`) hoặc hàm print console thừa trong các file này.

### Non-Functional
- [x] Cấu trúc code sạch, sử dụng Type Hinting của Python.
- [x] Quản lý đường dẫn tới file JSON an toàn trong môi trường Vercel (dùng `os.path`).

## Implementation Steps
1. [x] Tạo `vercel-backend/core/parser.py` (từ `html_parser.py`).
2. [x] Tạo thư mục `vercel-backend/core/engine/` và copy các file checker cũ vào.
3. [x] Tạo `vercel-backend/core/rules.py` để load `vaccine_rules.json`.
4. [x] Tạo wrapper class `AnalyzerService` nhận input là raw data (đã parse) và gọi Engine trả ra dict/json kết quả.

## Files to Create/Modify
- `vercel-backend/core/parser.py`
- `vercel-backend/core/engine/*`
- `vercel-backend/core/rules.py`
- `vercel-backend/assets/vaccine_rules.json`

## Test Criteria
- [x] Viết một đoạn script test nhỏ (hoặc Unit Test) feed HTML giả vào Parser -> Engine.
- [x] Kết quả trả về đúng định dạng JSON mong muốn.

---
Next Phase: `phase-03-fastapi-endpoints.md`
