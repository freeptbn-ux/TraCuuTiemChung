# Phase 05: Final Build & Verification
Status: ⬜ Pending
Dependencies: Phase 01-04

## Objective
Tạo bản build hoàn chỉnh và kiểm tra tổng thể hệ thống.

## Implementation Steps
1. [ ] **Create Signing Key:** Tạo Keystore mới để ký app (nếu chưa có).
2. [ ] **Generate Release Bundle:** Build -> Generate Signed Bundle / APK.
3. [ ] **Final Smoke Test:**
    - Cài APK lên máy thật.
    - Test luồng: Tra cứu -> Lấy dữ liệu từ cổng -> Backend xử lý -> App hiển thị kết quả.
4. [ ] **Documentation Update:** Cập nhật README.md với hướng dẫn cài đặt bản Production.

## Test Criteria
- [ ] App chạy mượt mà trên máy thật.
- [ ] Thời gian phản hồi API chấp nhận được trên mạng di động (4G/5G).
