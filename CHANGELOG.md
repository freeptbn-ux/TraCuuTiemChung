# Changelog

## [2026-05-12]
### Added
- **Vercel Backend**: Hệ thống FastAPI hỗ trợ tra cứu và phân tích tiêm chủng từ xa.
- **Upstash Redis Integration**: Tích hợp Redis để cache session Portal, giảm độ trễ và tránh lỗi timeout trên serverless.
- **Production Readiness Tests**: Thêm bộ test kiểm tra điều kiện môi trường trước khi deploy.
- **X-API-KEY Security**: Bảo mật API backend bằng header tùy chỉnh, cấu hình qua biến môi trường.
- **Improved Parsing**: Parser thông minh hỗ trợ bóc tách nhãn phụ (sublabel) cho các loại vắc-xin (Varivax, BCG, MMR...).
- **System Architecture**: Tài liệu kiến trúc hệ thống và hướng dẫn deploy chi tiết.

### Changed
- **Vercel Config**: Chuyển sang sử dụng `Root Directory` trỏ trực tiếp vào `vercel-backend/`.
- **API Headers**: Đổi từ `API_KEY` sang `X-API-KEY` (sử dụng Alias trong FastAPI) để tăng tính bảo mật và chuẩn hóa.
- **README.md**: Viết lại tài liệu hướng dẫn bằng tiếng Việt chuyên nghiệp.

### Fixed
- Lỗi Build trên Vercel do cấu hình `vercel.json` chuẩn cũ và sai đường dẫn đích (destination).
- Lỗi timeout khi login Portal bằng cách chuyển sang mô hình Session Caching với Upstash.

---
*Lưu ý: Phiên bản này đánh dấu sự chuyển dịch sang mô hình Client-Server.*
