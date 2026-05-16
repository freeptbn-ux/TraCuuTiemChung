# 💉 Tra Cứu Tiêm Chủng (Vaccination Lookup)

Ứng dụng Android giúp tra cứu lịch sử tiêm chủng và nhận khuyến nghị vaccine cá nhân hóa bằng cách tích hợp trực tiếp với cổng thông tin VNCDC. Hệ thống sử dụng engine phân tích thông minh để đưa ra các đề xuất dựa trên độ tuổi, lịch sử tiêm và các quy tắc y tế phức tạp.

---

## 🏗️ Kiến trúc hệ thống

```mermaid
graph TD
    A[Android App (Kotlin/Compose)] -- JSON (X-API-KEY) --> B[Vercel Backend (Go)]
    B -- Scraping (Cookie Jar) --> C[Cổng VNCDC (vncdc.gov.vn)]
    B -- Cache/Lock --> D[Upstash Redis]
    B -- Engine --> E[Vaccine Rules JSON]
```

---

## 🚀 Tính năng chính

### 📱 Ứng dụng Android
- **Tra cứu nhanh**: Tìm kiếm hồ sơ bệnh nhân qua số điện thoại.
- **Lịch sử chi tiết**: Hiển thị danh sách các mũi tiêm đã thực hiện từ cổng VNCDC.
- **Phân tích thông minh**: Đề xuất các mũi tiêm còn thiếu hoặc đến lịch.
- **Cảnh báo an toàn**: Phát hiện tương tác giữa các loại vaccine (VD: khoảng cách giữa các vaccine sống).
- **Giao diện hiện đại**: Sử dụng Jetpack Compose với Material Design 3.

### ⚙️ Backend (Vercel Go)
- **Scraping Engine**: Tự động đăng nhập và trích xuất dữ liệu từ cổng VNCDC.
- **High Performance**: Viết bằng Go, tối ưu hóa tốc độ xử lý và bộ nhớ.
- **Distributed Locking**: Sử dụng Redis để quản lý session và tránh xung đột khi đăng nhập đồng thời.
- **Stateless Architecture**: Chạy hoàn hảo trên môi trường Serverless của Vercel.

---

## 🛠️ Công nghệ sử dụng

| Thành phần | Công nghệ |
|---|---|
| **Mobile** | Kotlin, Jetpack Compose, Retrofit, OkHttp, Coroutines, Kotlinx Serialization |
| **Backend** | Go (Golang), Gin Gonic, Resty, GoQuery, slog |
| **Dữ liệu** | Upstash Redis (Session & Rate Limiting) |
| **Infrastructure** | Vercel (Serverless Functions) |

---

## 📁 Cấu trúc thư mục

```
.
├── app/                # Mã nguồn ứng dụng Android (Kotlin)
├── vercel-backend/      # Backend Go chạy trên Vercel
│   ├── api/            # Các endpoint API (lookup, analyze)
│   ├── pkg/            # Thư viện logic (portal, analyzer, middleware)
│   └── assets/         # Quy tắc tiêm chủng (vaccine_rules.json)
├── engine/             # Engine phân tích dùng chung (Go)
├── models/             # Định nghĩa cấu trúc dữ liệu (Go)
└── .brain/             # Dữ liệu tri thức dự án (AI Context)
```

---

## 🔧 Hướng dẫn cài đặt

### 1. Cấu hình Backend (Vercel)
Tạo dự án trên Vercel và thiết lập các biến môi trường sau:
- `X_API_KEY`: Khóa bí mật dùng để xác thực giữa App và Backend.
- `PORTAL_USERNAME`: Tài khoản cổng VNCDC.
- `PORTAL_PASSWORD`: Mật khẩu cổng VNCDC.
- `UPSTASH_REDIS_REST_URL`: URL Redis từ Upstash.
- `UPSTASH_REDIS_REST_TOKEN`: Token Redis từ Upstash.

### 2. Cấu hình Android App
Tạo file `local.properties` trong thư mục gốc:
```properties
X_API_KEY=your_api_key_here
BASE_URL=https://your-project.vercel.app/
```

### 3. Build & Chạy
- **Android**: Mở bằng Android Studio và nhấn Run, hoặc dùng lệnh:
  ```bash
  ./gradlew assembleDebug
  ```
- **Backend (Local)**:
  ```bash
  cd vercel-backend
  go run cmd/server/main.go
  ```

---

## 🔐 Bảo mật & Quy tắc
- Tuyệt đối không commit các file chứa API Key (`.env`, `local.properties`).
- Backend sử dụng Middleware xác thực `X-API-KEY` cho mọi request.
- Mọi dữ liệu nhạy cảm được xử lý qua HTTPS.

---

## 📜 Bản quyền
Copyright 2026 Nguyễn Duy Trường

---
*Dự án được phát triển và tối ưu hóa bởi Antigravity Workflow Framework.*
