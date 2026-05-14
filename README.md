# Tra Cứu Tiêm Chủng (Vaccination Lookup)

Hệ thống tra cứu và phân tích lịch sử tiêm chủng cá nhân từ Cổng thông tin tiêm chủng quốc gia.

## 🚀 Giới thiệu
Dự án bao gồm ứng dụng Android và hệ thống backend hiệu năng cao được viết bằng Go, cho phép người dùng tra cứu thông tin tiêm chủng chỉ bằng số điện thoại và nhận được phân tích chi tiết về các mũi tiêm còn thiếu hoặc cần tiêm trong tương lai.

## ✨ Tính năng chính
- **Tra cứu nhanh**: Chỉ cần số điện thoại để tìm kiếm thông tin các thành viên trong gia đình.
- **Phân tích thông minh**: Engine phân tích tự động kiểm tra lịch sử tiêm dựa trên các quy tắc y tế (Rota, Nhật Bản B, Phế cầu, v.v.).
- **Giao diện hiện đại**: Ứng dụng Android sử dụng Jetpack Compose và Material 3, hỗ trợ các tính năng mới nhất như Predictive Back Gesture.
- **Backend hiệu năng cao**: Chuyển đổi từ Python sang Go để tối ưu hóa tốc độ và khả năng xử lý trên Vercel.

## 🛠 Công nghệ sử dụng
- **Frontend**: Android (Kotlin, Jetpack Compose, Material 3).
- **Backend**: Go (Gin Web Framework, Goquery, Resty).
- **Deployment**: Vercel (Serverless Functions).
- **Quy tắc tiêm chủng**: Engine xử lý dựa trên dữ liệu JSON linh hoạt.

## 📁 Cấu trúc thư mục
- `app/`: Mã nguồn ứng dụng Android.
- `vercel-backend/`: Mã nguồn API Backend (Go).
  - `api/`: Entry points cho Vercel.
  - `internal/`: Logic xử lý (Analyzer, Portal Client).
  - `assets/`: Chứa các quy tắc tiêm chủng (`vaccine_rules.json`).
- `.brain/`: Lưu trữ ngữ cảnh và tri thức của dự án.

## ⚙️ Cài đặt & Sử dụng

### Backend
1. Chuyển vào thư mục backend: `cd vercel-backend`
2. Cài đặt dependency: `go mod download`
3. Chạy local: `go run api/index.go` (Cần cấu hình biến môi trường trong `.env`)

### Android
1. Mở dự án bằng Android Studio.
2. Build và chạy trên thiết bị Android hoặc Emulator.

## 📝 Bản quyền
Copyright 2026 Nguyễn Duy Trường

---
*Dự án này được phát triển và tối ưu hóa bởi Antigravity AI.*
