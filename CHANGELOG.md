# Changelog

Tất cả các thay đổi quan trọng đối với dự án Tra Cứu Tiêm Chủng sẽ được ghi lại tại đây.

## [2026-05-16] - Ổn định Backend & Tích hợp Android

### Added
- Thêm file `gradlew.bat` ở gốc project để hỗ trợ build trên Windows.
- Thêm cơ chế `go:embed` trong Go backend để đóng gói `vaccine_rules.json` vào file thực thi.
- Tạo `README.md` (Tiếng Việt) hướng dẫn chi tiết về dự án.
- Tạo `vercel-backend/assets/embed.go` để quản lý dữ liệu nhúng.

### Changed
- Cập nhật `VercelPortalRepository.kt` để trỏ đúng vào Vercel production API.
- Refactor `analyzer.Engine` hỗ trợ khởi tạo từ byte memory.
- Cấu hình `X-API-KEY` trong `local.properties` cho Android.
- Buộc Login portal VNCDC trước mỗi lần phân tích để đảm bảo session trong môi trường stateless.

### Fixed
- Sửa lỗi `no such file or directory` khi nạp quy tắc tiêm chủng trên Vercel.
- Sửa lỗi không tìm thấy dữ liệu tiêm chủng do mất session cookie giữa các lần gọi serverless function.
- Sửa lỗi `path/filepath` unused trong backend.
