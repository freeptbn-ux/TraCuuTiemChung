# TraCuuTiemChung (Tra Cứu Tiêm Chủng)

Dự án cung cấp một hệ thống phân tích và tra cứu lịch tiêm chủng thông minh. Dựa trên dữ liệu tiêm chủng của bệnh nhân (như từ VNVC), hệ thống có khả năng phân tích phác đồ, tính toán mũi tiêm tiếp theo, cảnh báo khoảng cách, giới hạn độ tuổi, và kiểm tra các tương tác giữa các loại vắc xin với độ chính xác cao.

## 🚀 Công nghệ sử dụng
- **Golang**: Xây dựng Engine phân tích lịch tiêm chủng hiệu năng cao (`vercel-backend`). Đóng vai trò là backend chính triển khai dưới dạng serverless.
- **Android (Kotlin)**: Ứng dụng di động giúp người dùng tiện lợi tra cứu và nhận thông báo về lịch tiêm chủng.
- **Python**: Các script hỗ trợ kiểm thử tính chính xác (parity test) so với hệ thống cũ và xử lý dữ liệu.
- **Vercel**: Nền tảng Serverless để triển khai Golang Backend một cách tự động và linh hoạt.
- **Redis**: Quản lý phiên đăng nhập (session persistence) và Rate Limit để bảo vệ API.

## 📂 Cấu trúc thư mục
- `/vercel-backend`: Mã nguồn chính của hệ thống backend viết bằng Go, sẵn sàng deploy lên Vercel. Chứa engine phân tích (`pkg/analyzer`).
- `/app`: Mã nguồn ứng dụng Android viết bằng Kotlin.
- `/tests`: Chứa các script Python để crawl dữ liệu giả lập và sinh golden files cho unit testing.
- `/docs` và `/plans`: Tài liệu thiết kế hệ thống và theo dõi quá trình phát triển (với AWF).
- `/.brain`: Thư mục chứa cấu hình trí tuệ nhân tạo AWF của dự án (luôn được đồng bộ).

## 🛠 Hướng dẫn cài đặt và chạy thử (Backend)
1. **Yêu cầu hệ thống**:
   - Golang 1.20+
   - Redis (nếu chạy local hoàn chỉnh với tính năng portal)

2. **Cài đặt dependencies**:
   ```bash
   cd vercel-backend
   go mod download
   ```

3. **Chạy các Unit Tests (Bao gồm Parity Tests)**:
   ```bash
   cd vercel-backend
   go test ./... -v
   ```

4. **Chạy server local**:
   ```bash
   cd vercel-backend
   go run api/index.go
   ```

## 📝 Bản quyền
Copyright 2026 Nguyễn Duy Trường.
Mọi quyền được bảo lưu.
