# Plan: Vercel Backend Integration
Created: 2026-05-12T16:17
Status: 🟡 In Progress

## Overview
Dự án di chuyển core logic (Portal Integration + Vaccine Analysis Engine) từ source code cũ (`VaccineAnalyzer-Pro-main`) lên hạ tầng **Vercel Serverless (Python 3.9+)**. Mục tiêu là để ứng dụng Android (`TraCuuTiemChung`) chỉ cần đóng vai trò là Client giao tiếp qua REST API (gửi SĐT -> nhận danh sách -> nhận dự báo), giúp app nhẹ, bảo mật và dễ cập nhật rule logic ở 1 nơi duy nhất.

## Tech Stack
- **Backend Framework:** FastAPI (Python)
- **Deployment:** Vercel (Serverless Functions)
- **Database/Cache:** Upstash Redis (lưu Session Cookie của portal để tối ưu thời gian phản hồi)
- **Web Scraping:** `requests` + `BeautifulSoup4` (không dùng headless browser)
- **Client Integration:** Kotlin Coroutines & Retrofit (Android)

## Kiến trúc (Architecture)
- **App (Android)** gọi `GET /api/lookup?phone=098xxx` -> **Vercel**
- **Vercel** check Session từ Redis. Nếu hết hạn -> login lại portal -> lưu Redis.
- **Vercel** dùng Session gửi HTTP Request qua cổng CDC.
- **Vercel** bóc tách HTML (`html_parser.py`) -> Trả JSON về Android.
- Tương tự cho bước Analyze (`GET /api/analyze?id=...`).

## Phases

| Phase | Name | Status | Progress |
|-------|------|--------|----------|
| 01 | Setup Environment & Vercel Config | ⬜ Pending | 0% |
| 02 | Porting Core Logic (Engine & Parser) | ⬜ Pending | 0% |
| 03 | FastAPI Endpoints | ⬜ Pending | 0% |
| 04 | Redis Session Integration | ⬜ Pending | 0% |
| 05 | Android Integration & Testing | ⬜ Pending | 0% |

## Quick Commands
- Bắt đầu Phase 1: `/code phase-01`
- Kiểm tra tiến độ: `/next`
- Lưu context: `/save-brain`
