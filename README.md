# TraCuuTiemChung (Vaccination Analyzer Pro) 🛡️

**TraCuuTiemChung** là một hệ thống hiện đại giúp tra cứu và phân tích lịch sử tiêm chủng từ cổng thông tin tiêm chủng quốc gia (VNCDC). Hệ thống bao gồm ứng dụng Android (Jetpack Compose) và Backend (FastAPI/Vercel) được thiết kế tối ưu cho hiệu năng, bảo mật và trải nghiệm người dùng cao cấp.

## 🏗️ Kiến trúc Hệ thống

Dự án được chuyển đổi từ mô hình xử lý cục bộ sang kiến trúc Client-Server để tối ưu hóa tài nguyên di động và đảm bảo tính ổn định:

1.  **Android App (Client)**: Giao diện hiện đại với Jetpack Compose, xử lý giao tiếp API qua Retrofit.
2.  **Vercel Backend (Server)**: FastAPI (Python) chịu trách nhiệm cào dữ liệu (scraping), xử lý session và chạy engine phân tích mũi tiêm chuyên sâu.
3.  **Redis Cache (Optional)**: Sử dụng Upstash Redis để duy trì session auth, giúp giảm độ trễ khi tra cứu.

## ✨ Tính năng nổi bật

- **Tra cứu siêu tốc**: Tìm kiếm hồ sơ tiêm chủng theo số điện thoại qua hệ thống Backend tối ưu.
- **Phân tích thông minh**: Engine phân tích đạt độ chính xác cao, tự động khuyến nghị các mũi tiêm tiếp theo, liều nhắc lại và phát hiện lỗi lịch tiêm dựa trên tiêu chuẩn y tế.
- **Giao diện Premium**: Sử dụng Material 3, hỗ trợ Dark Mode, Skeleton loading và các micro-animations mượt mà.
- **Bảo mật**: Kết nối API được bảo vệ bằng X-API-KEY và lưu trữ thông tin an toàn.
- **Quản lý đa bệnh nhân**: Dễ dàng chuyển đổi và xem lịch sử tiêm của nhiều người trong cùng một số điện thoại.

## 🚀 Công nghệ sử dụng

### Android Side
- **UI**: Jetpack Compose (Material 3)
- **Networking**: Retrofit 2 & OkHttp
- **Architecture**: MVVM + Clean Architecture
- **DI**: Manual Dependency Injection (Simple & Robust)
- **Serialization**: Kotlinx Serialization

### Backend Side (FastAPI)
- **Framework**: FastAPI (Python 3.12)
- **Parser**: BeautifulSoup4 & LXML
- **Deployment**: Vercel Serverless Functions
- **Cache**: Upstash Redis (Session Management)

## 🛠️ Hướng dẫn cài đặt

### 1. Backend (Vercel)
- Yêu cầu: Python 3.9+, Vercel CLI.
- Cấu hình file `.env` trong thư mục `vercel-backend`:
  ```env
  API_KEY=your_secret_api_key
  REDIS_URL=your_upstash_redis_url
  ```
- Deploy: `vercel deploy`

### 2. Android App
- Yêu cầu: Android Studio Ladybug+, JDK 17.
- Cấu hình `VERCEL_API_URL` và `VERCEL_API_KEY` trong `app/build.gradle.kts` hoặc thông qua biến môi trường.
- Build và chạy trực tiếp từ Android Studio.

## 🧪 Kiểm thử
- **Android**: `./gradlew test` (Bao gồm các bài test tích hợp với MockWebServer).
- **Backend**: Chạy script test trong `vercel-backend/tests`.

---
*Phát triển bởi đội ngũ tâm huyết nhằm mục tiêu cải thiện sức khỏe cộng đồng thông qua công nghệ.*
