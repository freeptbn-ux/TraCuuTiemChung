# Tra Cứu Tiêm Chủng (VNCDC Tracker & Analyzer)

![Vercel Deployment](https://img.shields.io/badge/Vercel-Deployment-success?style=flat-square&logo=vercel)
![Android Kotlin](https://img.shields.io/badge/Android-Kotlin-blue?style=flat-square&logo=android)
![FastAPI Python](https://img.shields.io/badge/FastAPI-Python-green?style=flat-square&logo=fastapi)
![Upstash Redis](https://img.shields.io/badge/Redis-Upstash-red?style=flat-square&logo=redis)

**Tra Cứu Tiêm Chủng** là giải pháp toàn diện hỗ trợ người dân theo dõi lịch sử tiêm chủng cá nhân. Ứng dụng tự động kết nối với Cổng thông tin Tiêm chủng Quốc gia (VNCDC), trích xuất dữ liệu và sử dụng Engine phân tích thông minh để đưa ra các gợi ý tiêm chủng chính xác theo quy định y tế.

## 🌟 Tính năng nổi bật

- **Tra cứu thời gian thực**: Chỉ cần số điện thoại để lấy danh sách bệnh nhân và lịch sử tiêm chủng từ hệ thống VNCDC.
- **Engine phân tích thông minh**: 
    - Tự động nhận diện hơn 12 nhóm quy tắc tiêm chủng (Single Series, Age-Dependent, MMR interaction, Alternative Courses...).
    - Cảnh báo các mũi tiêm bị thiếu hoặc đến hạn tiêm tiếp theo.
    - Xử lý các trường hợp tiêm trộn (mixing) vaccine phức tạp (ví dụ: Phế cầu, 6-trong-1).
- **Bảo mật tối đa**:
    - Sử dụng mã hóa AES với Android Keystore để bảo vệ thông tin đăng nhập.
    - Token và session được quản lý qua Upstash Redis trên Cloud.
- **Giao diện hiện đại**: Xây dựng hoàn toàn bằng Jetpack Compose với phong cách Material 3.

## 🛠️ Công nghệ sử dụng

### Mobile (Android App)
- **Ngôn ngữ**: Kotlin (Modern Android Development)
- **UI Framework**: Jetpack Compose
- **Kiến trúc**: Clean Architecture + MVVM
- **Network**: Retrofit 2 + OkHttp
- **Local Security**: DataStore + Tink/Keystore Encryption
- **HTML Parsing**: Jsoup

### Backend (Cloud API)
- **Framework**: FastAPI (Python 3.12)
- **Runtime**: Vercel Serverless Functions
- **Caching**: Upstash Redis (Serverless-optimized)
- **Security**: Custom API Key (X-API-KEY) Header

## 📂 Cấu trúc dự án

```text
TraCuuTiemChung/
├── app/                        # Module Android (Kotlin)
│   ├── src/main/java/          # Source code theo Clean Architecture
│   │   ├── data/               # Tầng dữ liệu (Auth, Portal Client, Parsers)
│   │   ├── domain/             # Tầng nghiệp vụ (Engine phân tích, Rule Dispatcher)
│   │   └── ui/                 # Tầng giao diện (Compose Screens, ViewModels)
│   └── build.gradle.kts        # Cấu hình build Android
├── vercel-backend/             # Hệ thống API Backend (Python)
│   ├── api/                    # Serverless Functions (FastAPI entrypoint)
│   ├── core/                   # Shared logic engine chuyển từ Python sang
│   ├── services/               # Kết nối Portal & Upstash Redis
│   ├── requirements.txt        # Danh sách thư viện Python
│   └── vercel.json             # Cấu hình deployment Cloud
├── .brain/                     # Dữ liệu tri thức hỗ trợ phát triển (Project Context)
├── docs/                       # Tài liệu thiết kế và đặc tả kỹ thuật
├── plans/                      # Kế hoạch phát triển theo từng giai đoạn
└── README.md                   # Hướng dẫn dự án
```

## 🚀 Hướng dẫn cài đặt & Triển khai

### 1. Triển khai Backend (Vercel)
1. Truy cập [Vercel](https://vercel.com) và tạo project mới từ repository này.
2. Thiết lập **Root Directory** là `vercel-backend`.
3. Thêm các Environment Variables:
   - `UPSTASH_REDIS_REST_URL`: URL từ Upstash Console.
   - `UPSTASH_REDIS_REST_TOKEN`: Token từ Upstash Console.
   - `X_API_KEY`: Key tự định nghĩa để bảo mật API.
4. Deploy và lấy URL của backend.

### 2. Cài đặt App Android
1. Mở project bằng Android Studio (Koala hoặc mới hơn).
2. Cập nhật URL backend vào cấu hình API trong app.
3. Chạy ứng dụng trên máy ảo hoặc thiết bị thật (Android 7.0+).

## 📖 Cách sử dụng
1. Mở app, đăng nhập bằng tài khoản Cổng Tiêm Chủng (nếu cần).
2. Nhập số điện thoại để tra cứu.
3. Chọn người cần xem, hệ thống sẽ hiển thị danh sách các mũi đã tiêm và **Phân tích kết quả** (các mũi cần tiêm tiếp theo).

## ⚖️ Bản quyền
Copyright 2026 Nguyễn Duy Trường

---
*Dự án được xây dựng với mục tiêu nâng cao nhận thức cộng đồng về tiêm chủng phòng bệnh.*
