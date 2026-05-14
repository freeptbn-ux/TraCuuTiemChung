# Changelog

Tất cả những thay đổi quan trọng của dự án sẽ được lưu vết tại đây.

## [2026-05-14] - Backend Go Migration & Cleanup

### Added
- **Backend Go**: Khởi tạo và hoàn thiện module backend bằng Go thay thế cho Python.
- **API Handlers**: Triển khai `api/index.go` với các endpoint `/api/lookup` và `/api/analyze`.
- **Analyzer Engine**: Porting thành công engine phân tích quy tắc tiêm chủng từ Python sang Go.
- **Portal Client**: Triển khai client cào dữ liệu từ Cổng thông tin Tiêm chủng Quốc gia với hỗ trợ cookie và session.
- **Unit Tests**: Hệ thống test toàn diện cho logic analyzer, portal client và cấu hình.
- **Vercel Config**: Cấu hình `vercel.json` và `.vercelignore` tối ưu cho Go runtime.
- **Documentation**: Tạo file `README.md` chi tiết bằng tiếng Việt.

### Changed
- Cấu trúc thư mục: Chuyển toàn bộ logic backend vào thư mục `vercel-backend/`.
- Cải thiện logic phân tích nhóm vaccine Phế cầu (Pneumo) và Nhật Bản B (JE).
- Tăng cường bảo mật API bằng middleware kiểm tra `X-API-KEY`.

### Fixed
- Lỗi xung đột hàm helper trong các unit test của Go.
- Lỗi phân tích phác đồ trộn (mixed series) trong engine.

### Removed
- **Legacy Python**: Di chuyển toàn bộ code Python cũ vào thư mục `python_backup/` để dọn dẹp không gian làm việc chính.

---
*Cập nhật bởi Antigravity AI.*
