# Phase 05: Android Integration & Final Testing
Status: ✅ Completed
Dependencies: phase-04

## Objective
Kết nối ứng dụng Android (`TraCuuTiemChung`) với Backend Vercel vừa tạo. Thay thế hoàn toàn các logic Engine và Portal cũ dưới app bằng các cuộc gọi API.

## Requirements
### Functional
- [x] Cấu hình URL Vercel và `X-API-KEY` trong Android app (`local.properties` hoặc `BuildConfig`).
- [x] Định nghĩa Retrofit API interface cho `POST /api/lookup` và `POST /api/analyze`.
- [x] Cập nhật `LookupVaccinationByPhoneUseCase` để gọi API thay vì chạy cục bộ.
- [x] Cập nhật UI Android để hỗ trợ luồng: Nhập SĐT -> Chờ -> Hiện list Tên -> Chọn 1 người -> Chờ -> Hiện kết quả phân tích. (Trước đây nếu làm gộp sẽ tốn thời gian, luồng tách biệt này sẽ tốt hơn).

### Non-Functional
- [x] Hiển thị Skeleton/Loading state rõ ràng trong lúc chờ API.
- [x] Bắt lỗi Network, Timeout từ Vercel và hiển thị cảnh báo (Snackbar/Toast) thân thiện.

## Implementation Steps
1. [x] Deploy code Vercel (bằng CLI `vercel` hoặc push lên Github và auto deploy).
2. [x] Sửa `RetrofitClient.kt` trong app Android để trỏ tới domain của Vercel.
3. [x] Viết Data Models (Kotlin data classes) khớp với JSON mà Vercel trả về.
4. [x] Khởi chạy App, thực hiện một ca kiểm thử thực tế từ đầu đến cuối.

## Files to Create/Modify (Android Side)
- `app/src/main/java/com/tracuutiemchung/app/data/remote/VercelApiService.kt`
- `app/src/main/java/com/tracuutiemchung/app/domain/usecase/LookupVaccinationByPhoneUseCase.kt`
- Các file View/ViewModel liên quan để hiển thị List Danh sách thay vì phi thẳng vào Report.

## Test Criteria
- [x] App không bị crash nếu mạng chậm.
- [x] Hiển thị đúng kết quả phân tích mũi tiêm cho một bệnh nhân thực tế.

---
**END OF PLAN** - Báo cáo với User để duyệt.
