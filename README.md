# Tra Cứu Tiêm Chủng (VNCDC Tracker & Analyzer)

![Vercel Deployment](https://img.shields.io/badge/Vercel-Deployment-success?style=flat-square&logo=vercel)
![Android Kotlin](https://img.shields.io/badge/Android-Kotlin-blue?style=flat-square&logo=android)
![FastAPI Python](https://img.shields.io/badge/FastAPI-Python-green?style=flat-square&logo=fastapi)

Ứng dụng hỗ trợ tra cứu lịch sử tiêm chủng cá nhân từ Cổng thông tin Tiêm chủng Quốc gia (VNCDC) và tự động phân tích các mũi tiêm còn thiếu dựa trên các quy tắc y tế phức tạp.

## 🌟 Tính năng chính

- **Tra cứu nhanh**: Tìm kiếm thông tin tiêm chủng chỉ bằng số điện thoại.
*   **Phân tích thông minh**: Engine phân tích (Engine Dispatcher) tự động kiểm tra lịch sử tiêm chủng đối với hơn 12 nhóm quy tắc (Single Series, Age Dependent, MMR Equivalent, Alternative Courses, v.v.).
- **Bảo mật**: Quản lý thông tin đăng nhập an toàn bằng Android DataStore và mã hóa Keystore.
- **Backend Serverless**: Hệ thống backend chạy trên Vercel tối ưu chi phí và hiệu năng, tích hợp Upstash Redis để quản lý session portal.

## 🛠️ Công nghệ sử dụng

### Mobile (Android)
- **Ngôn ngữ**: Kotlin
- **UI Framework**: Jetpack Compose (Material 3)
- **Mạng**: OkHttp / Retrofit
- **Phân tích HTML**: Jsoup
- **Bảo mật**: DataStore + AES Encryption (Keystore)

### Backend
- **Framework**: FastAPI (Python 3.12)
- **Platform**: Vercel
- **Database/Cache**: Upstash Redis
- **Scraping**: BeautifulSoup4 / Requests

## 📂 Cấu trúc dự án

```text
TraCuuTiemChung/
├── app/                        # Module Android chính (Kotlin)
│   ├── src/main/java/.../data/  # Tầng dữ liệu (Portal client, Parsers)
│   ├── src/main/java/.../domain/# Tầng nghiệp vụ (Analysis Engine)
│   └── src/main/java/.../ui/    # Giao diện Jetpack Compose
├── vercel-backend/             # Hệ thống API Backend (Python)
│   ├── api/                     # Điểm cuối API (index.py)
│   ├── services/                # Logic portal client & Redis cache
│   └── core/                    # Engine phân tích chuyển từ Python sang
├── plans/                      # Các tài liệu kế hoạch triển khai (Phase-based)
├── .brain/                     # Dữ liệu tri thức hỗ trợ phát triển AI
└── README.md                   # Hướng dẫn này
```

## 🚀 Hướng dẫn cài đặt

### Backend (Vercel)
1. Cấu hình các biến môi trường trong `.env`:
   - `UPSTASH_REDIS_REST_URL`
   - `UPSTASH_REDIS_REST_TOKEN`
   - `X_API_KEY` (Key bảo mật giữa App và Server)
2. Deploy lên Vercel:
   ```bash
   vercel deploy
   ```

### Android App
1. Mở thư mục `TraCuuTiemChung` bằng Android Studio.
2. Build và cài đặt lên thiết bị Android (API 24+).
3. Cấu hình Base URL của API trong code nếu cần (mặc định trỏ về Vercel production).

## 📖 Cách sử dụng
1. Đăng nhập vào cổng portal bằng tài khoản VNCDC.
2. Nhập số điện thoại để tra cứu danh sách bệnh nhân liên quan.
3. Chọn bệnh nhân để xem chi tiết lịch sử tiêm chủng và kết quả phân tích mũi tiêm còn thiếu/cần tiêm.

## ⚖️ Bản quyền
Copyright 2026 Nguyễn Duy Trường

---
*Dự án được phát triển nhằm hỗ trợ cộng đồng theo dõi sức khỏe tiêm chủng tốt hơn.*
