# 💉 Tra Cứu Tiêm Chủng

Ứng dụng Android giúp tra cứu lịch sử tiêm chủng và nhận khuyến nghị vaccine cá nhân hóa, tích hợp với cổng thông tin VNCDC.

---

## 🏗️ Kiến trúc hệ thống

```
┌─────────────────────┐        ┌──────────────────────┐        ┌─────────────────┐
│   Android App       │ ──────▶│  Vercel Backend (Go) │ ──────▶│  Cổng VNCDC    │
│   (Kotlin/Compose)  │ JSON   │  tracuutiemchung      │  HTTP  │  tiemchung.     │
│                     │◀────── │  .vercel.app          │◀────── │  vncdc.gov.vn   │
└─────────────────────┘        └──────────────────────┘        └─────────────────┘
```

### Thành phần chính

| Thư mục | Mô tả |
|---|---|
| `app/` | Ứng dụng Android (Kotlin + Jetpack Compose) |
| `vercel-backend/` | Backend Go chạy trên Vercel Serverless |
| `engine/` | Engine phân tích lịch tiêm chủng (Go) |

---

## 📱 Ứng dụng Android

### Tính năng
- Tra cứu bệnh nhân theo số điện thoại
- Hiển thị lịch sử các mũi tiêm đã thực hiện
- Phân tích và đề xuất các vaccine còn thiếu
- Cảnh báo tương tác giữa các loại vaccine

### Yêu cầu
- Android 8.0 (API 26) trở lên
- Kết nối Internet

### Cấu hình (`local.properties`)

```properties
sdk.dir=C:\Users\...\AppData\Local\Android\Sdk
X_API_KEY=your_api_key_here
```

### Build & Chạy

```bash
# Build debug APK
./gradlew assembleDebug

# Cài lên thiết bị
adb install app/build/outputs/apk/debug/app-debug.apk
```

---

## ⚙️ Backend Vercel

Backend viết bằng Go, chạy dưới dạng Serverless Function trên Vercel. Mỗi request tự đăng nhập vào cổng VNCDC, scrape dữ liệu và phân tích.

### API Endpoints

| Method | Endpoint | Mô tả |
|---|---|---|
| `GET` | `/api/health` | Kiểm tra trạng thái |
| `POST` | `/api/lookup` | Tìm bệnh nhân theo SĐT |
| `POST` | `/api/analyze` | Phân tích lịch tiêm của bệnh nhân |

### Xác thực

Tất cả request tới `/api/lookup` và `/api/analyze` phải có header:

```
X-API-KEY: <your_key>
```

### Biến môi trường (Vercel)

| Tên | Mô tả |
|---|---|
| `X_API_KEY` | Khóa bảo mật cho mobile app |
| `PORTAL_USERNAME` | Tên đăng nhập cổng VNCDC |
| `PORTAL_PASSWORD` | Mật khẩu cổng VNCDC |
| `UPSTASH_REDIS_REST_URL` | URL Redis (Upstash) để cache session |
| `UPSTASH_REDIS_REST_TOKEN` | Token Redis (Upstash) |

### Chạy local

```bash
cd vercel-backend
cp .env.example .env
# Điền thông tin vào .env

go run cmd/server/main.go
```

### Chạy test

```bash
cd vercel-backend

# Unit tests
go test ./...

# Integration test với dữ liệu thật (cần .env)
go test ./api/ -run TestLiveVNVC -v -timeout 120s
```

### Deploy

Push code lên nhánh `main` → Vercel tự động deploy.

```bash
git push origin main
```

---

## 🔬 Engine Phân Tích

Engine phân tích vaccine nằm tại `vercel-backend/pkg/analyzer/`, được thiết kế để:

- Áp dụng các quy tắc tiêm chủng từ file `assets/vaccine_rules.json`
- Tính toán khoảng cách tối thiểu giữa các mũi tiêm
- Phát hiện tương tác giữa các loại vaccine (VD: MMR và vaccine sống)
- Hỗ trợ nhiều loại quy tắc: chuỗi liều, phụ thuộc tuổi, nhóm đặc biệt

---

## 🔐 Bảo mật

- `local.properties` và `.env` **không được commit** lên Git
- API key được mã hóa trong `BuildConfig` (không hard-code trong source)
- Backend sử dụng Redis lock để tránh race condition khi nhiều request đăng nhập đồng thời

---

## 📝 Ghi chú phát triển

- Backend là **stateless serverless** — mỗi invocation phải tự đăng nhập lại
- Redis (Upstash) được dùng để cache cookie session và distributed lock
- Engine phân tích đã được kiểm tra parity với implementation Python legacy

---

## 📄 License

Dự án nội bộ. Không dùng cho mục đích thương mại.
