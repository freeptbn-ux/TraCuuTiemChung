# System Architecture Overview

## 🏗️ Tổng quan
Hệ thống **TraCuuTiemChung** được thiết kế theo kiến trúc Client-Server hiện đại, chuyển dịch toàn bộ logic xử lý nặng (Scraping, Parsing, Analysis) lên Cloud để tối ưu hóa ứng dụng di động.

## 📱 Android Client
- **Tech Stack**: Kotlin, Jetpack Compose, Retrofit.
- **Vai trò**: 
    - Hiển thị giao diện người dùng.
    - Nhận input (SĐT) và gửi yêu cầu tới Backend.
    - Hiển thị kết quả phân tích dưới dạng Dashboard/Report.
- **Security**: Lưu trữ API Key an toàn, không chứa logic cào dữ liệu nhạy cảm.

## ☁️ Vercel Backend (FastAPI)
- **Tech Stack**: Python 3.12, FastAPI, BeautifulSoup4.
- **Vai trò**:
    - **PortalClient**: Giả lập trình duyệt, xử lý đăng nhập và duy trì session tới VNCDC Portal.
    - **HTML Parser**: Bóc tách dữ liệu thô từ HTML Portal thành cấu trúc dữ liệu chuẩn.
    - **Analysis Engine**: Trái tim của hệ thống, chạy các quy tắc y tế để đưa ra khuyến nghị tiêm chủng.
- **Cache**: Tích hợp Redis để tối ưu hóa tốc độ truy cập session.

## 🔄 Luồng dữ liệu (Data Flow)
1. Người dùng nhập SĐT trên Android.
2. Android gửi request tới `/api/lookup`.
3. Backend cào Portal VNCDC -> Trả về danh sách bệnh nhân.
4. Người dùng chọn 1 bệnh nhân -> Request tới `/api/analyze`.
5. Backend cào Detail Page -> Parse -> Chạy Engine phân tích -> Trả về JSON khuyến nghị.
6. Android hiển thị Report.

## 🛠️ Infrastructure
- **Hosting**: Vercel (Serverless).
- **Database/Cache**: Upstash Redis.
- **CI/CD**: Tích hợp Github Action / Vercel Auto-deploy.
