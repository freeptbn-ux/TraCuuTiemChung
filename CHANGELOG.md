# Changelog

## [2026-05-12]
### Added
- **Vercel Backend**: Hệ thống FastAPI hỗ trợ tra cứu và phân tích tiêm chủng từ xa.
- **Android Integration**: Tích hợp Retrofit để kết nối với Backend Vercel.
- **Improved Parsing**: Parser thông minh hỗ trợ bóc tách nhãn phụ (sublabel) cho các loại vắc-xin (Varivax, BCG, MMR...).
- **API Documentation**: Tài liệu API chi tiết trong `docs/api/`.
- **System Architecture**: Tài liệu kiến trúc hệ thống trong `docs/architecture/`.
- **GitHub Sync**: Đẩy mã nguồn lên repository chính thức.

### Changed
- **MainActivity**: Đơn giản hóa luồng khởi chạy, bỏ qua màn hình Login.
- **Normalization**: Chuyển logic chuẩn hóa quy tắc về phía Backend khi khởi động.

### Fixed
- Lỗi không nhận diện được mũi tiêm Varivax (Thủy đậu) do trùng tên sublabel.
- Lỗi thiếu mũi tiêm Lao (BCG) bằng cách bổ sung brand name `IVACTUBER-BCG`.
- Lỗi biên dịch Unit Test cũ sau khi thay đổi cấu trúc App.

---
*Lưu ý: Phiên bản này đánh dấu sự chuyển dịch sang mô hình Client-Server.*
