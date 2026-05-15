# Tra Cứu Tiêm Chủng (VNCDC Tracker & Analyzer)

![Vercel Deployment](https://img.shields.io/badge/Vercel-Deployment-success?style=flat-square&logo=vercel)
![Android Kotlin](https://img.shields.io/badge/Android-Kotlin-blue?style=flat-square&logo=android)
![Go Backend](https://img.shields.io/badge/Go-Backend-00ADD8?style=flat-square&logo=go)
![Upstash Redis](https://img.shields.io/badge/Redis-Upstash-red?style=flat-square&logo=redis)

**Tra Cứu Tiêm Chủng** là giải pháp toàn diện hỗ trợ người dân theo dõi lịch sử tiêm chủng cá nhân. Hệ thống tự động kết nối với Cổng thông tin Tiêm chủng Quốc gia (VNCDC), trích xuất dữ liệu và sử dụng Engine phân tích thông minh để đưa ra các gợi ý tiêm chủng chính xác theo quy định y tế.

## 🚀 Giới thiệu
Dự án bao gồm ứng dụng Android và hệ thống backend hiệu năng cao được viết bằng **Go**, cho phép người dùng tra cứu thông tin tiêm chủng chỉ bằng số điện thoại và nhận được phân tích chi tiết về các mũi tiêm còn thiếu hoặc cần tiêm trong tương lai.

## ✨ Tính năng nổi bật
- **Tra cứu nhanh**: Chỉ cần số điện thoại để tìm kiếm thông tin các thành viên trong gia đình từ hệ thống VNCDC.
- **Engine phân tích thông minh**: 
    - Tự động nhận diện hơn 12 nhóm quy tắc tiêm chủng (Single Series, Age-Dependent, MMR interaction, Alternative Courses...).
    - Cảnh báo các mũi tiêm bị thiếu hoặc đến hạn tiêm tiếp theo.
    - Xử lý các trường hợp tiêm trộn (mixing) vaccine phức tạp (ví dụ: Phế cầu, Nhật Bản B).
- **Giao diện hiện đại**: Ứng dụng Android sử dụng Jetpack Compose và Material 3, hỗ trợ các tính năng mới nhất như Predictive Back Gesture.
- **Backend hiệu năng cao**: Hệ thống được xây dựng bằng Go (Gin framework) để tối ưu hóa tốc độ xử lý và khả năng mở rộng trên Vercel.

## 🛠 Công nghệ sử dụng
- **Frontend**: Android (Kotlin, Jetpack Compose, Material 3).
- **Backend**: Go (Gin Web Framework, Goquery, Resty).
- **Deployment**: Vercel (Serverless Functions).
- **Caching/Session**: Upstash Redis (Serverless-optimized).

## 📁 Cấu trúc thư mục
- `app/`: Mã nguồn ứng dụng Android.
- `vercel-backend/`: Mã nguồn API Backend (Go).
  - `api/`: Entry points cho Vercel.
  - `pkg/`: Logic xử lý (Analyzer, Portal Client).
  - `assets/`: Chứa các quy tắc tiêm chủng (`vaccine_rules.json`).
- `.brain/`: Lưu trữ ngữ cảnh và tri thức của dự án.
- `docs/`: Tài liệu thiết kế và đặc tả kỹ thuật.

## ⚙️ Hướng dẫn cài đặt & Triển khai

### 1. Triển khai Backend (Vercel)
1. Truy cập [Vercel](https://vercel.com) và tạo project mới từ repository này.
2. Thiết lập **Root Directory** là `vercel-backend`.
3. Thêm các **Environment Variables** sau:
   - `UPSTASH_REDIS_REST_URL`: URL kết nối Redis từ Upstash.
   - `UPSTASH_REDIS_REST_TOKEN`: REST Token từ Upstash.
   - `X_API_KEY`: Mã bảo mật để App Android có thể gọi API (phải khớp với `BuildConfig.X_API_KEY`).
   - `PORTAL_USERNAME`: Tài khoản đăng nhập Cổng tiêm chủng.
   - `PORTAL_PASSWORD`: Mật khẩu đăng nhập Cổng tiêm chủng.
4. Deploy và cấu hình Custom Domain nếu cần.

### 2. Cấu hình Redis
Hệ thống sử dụng Redis để:
- **Session Persistence**: Duy trì trạng thái đăng nhập giữa các Lambda instances của Vercel.
- **Distributed Locking**: Đảm bảo chỉ có một process thực hiện đăng nhập tại một thời điểm, tránh bị Portal chặn (Bypass Rate Limit).
- **API Rate Limiting**: Giới hạn tần suất gọi API từ App (mặc định 50 req/phút).

### 2. Chạy Backend Local
1. Chuyển vào thư mục backend: `cd vercel-backend`
2. Cài đặt dependency: `go mod download`
3. Chạy: `go run api/index.go` (Cần cấu hình `.env`)

### 3. Cài đặt App Android
1. Mở dự án bằng Android Studio (Koala+).
2. Build và chạy trên thiết bị Android hoặc Emulator.

## 🛠 Cấu hình cho AI Agent (MCP)
Dự án này hỗ trợ **ADB MCP Server** để hỗ trợ AI Agent tự động thao tác trên thiết bị thật:
- Cài đặt ADB Platform Tools.
- Thêm cấu hình sau vào `mcp_config.json`:
```json
"adb": {
  "command": "npx",
  "args": ["-y", "adb-mcp"]
}
```

## 📝 Cập nhật gần đây (v1.1.0)
- **Hệ thống Backend (Go)**: Hoàn tất chuyển đổi từ Python sang Go, tối ưu hóa tốc độ và bảo mật.
- **Middleware Hardening**: Tích hợp Rate Limiting, Request ID, và Centralized Error Handling.
- **Session Management**: Triển khai Redis Cookie Jar và Distributed Lock để xử lý concurrency.
- **E2E Verified**: Đã kiểm thử thành công trên thiết bị thật với dữ liệu thực tế từ VNCDC.

---
Copyright 2026 Nguyễn Duy Trường.
*Dự án được tối ưu hóa bởi Antigravity AI.*
