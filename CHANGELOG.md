# Changelog

Tất cả các thay đổi quan trọng đối với dự án Tra Cứu Tiêm Chủng sẽ được ghi lại tại đây.

## [2026-05-16] - Ổn định Backend & Tích hợp Android

### Fixed
- **Security:** Đổi mã lỗi xác thực sai API Key từ `403` sang `401 Unauthorized` để đồng bộ với App Android.
- **Session:** Sửa lỗi `RedisCookieJar` không gửi cookie do thiếu Domain/Path, cải thiện độ ổn định của phiên đăng nhập.
- **Search:** Bổ sung header `Accept-Language`, `X-Requested-With` và tham số `remember_me=true` để khớp hoàn toàn với bản legacy.
- **Robustness:** Sửa lỗi panic khi log snippet HTML của trang lỗi quá ngắn.
- **Verified:** Chạy bộ test `qa_test_report.py` trên local thành công, xác nhận kết quả khớp 100% với mã Python cũ (tìm thấy 2 bệnh nhân cho số 0388634123).

- Thêm file `gradlew.bat` ở gốc project để hỗ trợ build trên Windows.
- Thêm cơ chế `go:embed` trong Go backend để đóng gói `vaccine_rules.json` vào file thực thi.
- Tạo `README.md` (Tiếng Việt) hướng dẫn chi tiết về dự án.
- Tạo `vercel-backend/assets/embed.go` để quản lý dữ liệu nhúng.

### Changed
- Cập nhật `VercelPortalRepository.kt` để trỏ đúng vào Vercel production API.
- Refactor `analyzer.Engine` hỗ trợ khởi tạo từ byte memory.
- Cấu hình `X-API-KEY` trong `local.properties` cho Android.
- Buộc Login portal VNCDC trước mỗi lần phân tích để đảm bảo session trong môi trường stateless.

### [2026-05-16] - Phục hồi Tra cứu Portal VNCDC
- **Fixed:** Sửa lỗi tra cứu trả về 0 kết quả bằng cách chuyển từ `POST` sang `GET` (đồng bộ với code Python cũ).
- **Fixed:** Sửa lỗi `cookie_jar` bị chặn do lọc domain quá nghiêm ngặt.
- **Improved:** Thêm cơ chế tự động Login lại khi gặp trang lỗi ASP.NET ("Sorry, an error occurred").
- **Improved:** Rotate Redis session key (`v2`) để giải quyết tình trạng session bị corrupted.
- **Verified:** Kiểm tra thực tế bằng `check_block.py` trả về đúng 2 bệnh nhân trên Vercel Production.
