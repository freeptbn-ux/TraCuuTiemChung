# Tra Cứu Tiêm Chủng (VNCDC Portal Integration)

Hệ thống tra cứu lịch sử tiêm chủng và phân tích kế hoạch tiêm chủng tự động, tích hợp dữ liệu từ Cổng thông tin tiêm chủng Quốc gia (VNCDC).

## 🌟 Tính năng chính

- **Tra cứu nhanh:** Tìm kiếm thông tin bệnh nhân và lịch sử tiêm chủng qua số điện thoại từ hệ thống VNCDC.
- **Phân tích thông minh:** Tự động đối chiếu lịch sử tiêm với bộ quy tắc (`vaccine_rules.json`) để đưa ra khuyến nghị các mũi tiêm tiếp theo.
- **Quản lý phiên (Session Management):** Cơ chế quản lý Cookie Jar qua Redis (Upstash) giúp duy trì kết nối ổn định và tối ưu hóa tốc độ truy cập.
- **Đa nền tảng:** 
  - **Backend:** Chạy trên môi trường Serverless của Vercel (Go & Python Legacy).
  - **Mobile:** Ứng dụng Android viết bằng Kotlin (Jetpack Compose).

## 🛠 Công nghệ sử dụng

### Backend (Golang)
- **Framework:** Gin Gonic (Web Framework).
- **Library:** Resty (HTTP Client), Goquery (HTML Parsing).
- **Database/Cache:** Redis (Upstash) để lưu trữ Cookie Jar.
- **Deployment:** Vercel Serverless Functions.

### Android App
- **Language:** Kotlin.
- **UI Framework:** Jetpack Compose.
- **Network:** Retrofit + OkHttp.
- **Architecture:** MVVM.

### Legacy Backend (Python)
- **Framework:** Flask-like (Serverless).
- **Library:** Requests, BeautifulSoup4.

## 📁 Cấu trúc thư mục

- `app/`: Mã nguồn ứng dụng Android.
- `vercel-backend/`: Mã nguồn Backend hiện tại (Golang).
- `vercel-backend-python-legacy/`: Mã nguồn Backend cũ (Python) dùng để tham chiếu/backup.
- `.brain/`: Thư mục lưu trữ kiến thức và ngữ cảnh làm việc của AI Assistant.
- `docs/`: Tài liệu kiến trúc, hướng dẫn API và các giai đoạn phát triển.

## ⚙️ Hướng dẫn cài đặt

### Backend (Go)
1. Cài đặt [Go](https://golang.org/dl/) (phiên bản 1.22 trở lên).
2. Di chuyển vào thư mục `vercel-backend`.
3. Tạo file `.env` từ `.env.example` và điền các thông tin:
   - `PORTAL_USERNAME`, `PORTAL_PASSWORD`: Tài khoản portal VNCDC.
   - `UPSTASH_REDIS_REST_URL`, `UPSTASH_REDIS_REST_TOKEN`: Thông tin kết nối Redis.
   - `X_API_KEY`: Key bảo mật để App kết nối.
4. Chạy server local: `go run cmd/server/main.go`.

### Android
1. Mở project bằng Android Studio.
2. Cấu hình `X_API_KEY` trong file `local.properties`.
3. Build và chạy trên thiết bị giả lập hoặc thật.

## 📝 Cách sử dụng

- Sử dụng ứng dụng Android để nhập số điện thoại cần tra cứu.
- Backend sẽ kết nối với Portal VNCDC, lấy dữ liệu HTML, parse thông tin và trả về kết quả JSON cho App.
- Hệ thống tự động phân tích các mũi đã tiêm và hiển thị các mũi còn thiếu kèm thời gian tiêm dự kiến.

## ⚖️ Bản quyền

Copyright 2026 Nguyễn Duy Trường

---
*Dự án được phát triển và tối ưu hóa bởi Antigravity AI Assistant.*
