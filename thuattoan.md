# Báo cáo: So sánh thuật toán Login và cấu trúc trả kết quả về App Kotlin

Sau khi phân tích mã nguồn của hai hệ thống `vercel-backend` (Golang) và `vercel-backend-python-legacy` (Python), dưới đây là báo cáo so sánh chi tiết:

## 1. Thuật toán đăng nhập (Portal Login)
**Kết luận: Không giống nhau. Golang ưu việt hơn nhờ xử lý Race Condition.**

* **Python Legacy (`services/auth_service.py`):**
  - Lưu cache cookie trong Redis (TTL = 30 phút).
  - Khi cache hết hạn hoặc cần đăng nhập lại, tiến trình trực tiếp gửi GET request lấy CSRF Token (`__RequestVerificationToken`) và POST request để lấy cookie có `.ASPXAUTH`.
  - **Nhược điểm:** Khi có nhiều request tới cùng một lúc và cache vừa hết hạn, tất cả các request này sẽ đồng loạt gọi hàm đăng nhập. Điều này tạo ra tình trạng **Race Condition** (nhiều tiến trình đăng nhập đồng thời), làm chậm hệ thống và dễ bị VNVC chặn IP hoặc khóa tài khoản do spam request login.

* **Golang (`pkg/portal/client.go`):**
  - Cũng sử dụng Redis để lưu cookie.
  - **Điểm khác biệt:** Golang áp dụng **Distributed Lock (Khóa phân tán)** trong Redis.
    - Khi một request cần đăng nhập lại, nó sẽ cố gắng tạo một lock (`portal:lock:login:<username>`).
    - Nếu thành công (lấy được lock), nó thực hiện quy trình đăng nhập.
    - Nếu thất bại (có một request khác đang đăng nhập), nó **không** gọi đăng nhập nữa mà sẽ chuyển sang trạng thái chờ (polling) mỗi 500ms để đợi request kia lưu cookie xong. Nếu quá 10 giây sẽ timeout.
  - **Ưu điểm:** Khắc phục triệt để lỗi đăng nhập đè (Race Condition), hoạt động an toàn và cực kỳ ổn định trong môi trường Serverless (nơi auto-scale sinh ra rất nhiều instance cùng lúc).

## 2. Cấu trúc trả kết quả về cho App Kotlin (Android)
**Kết luận: Không giống nhau. Golang được thiết kế "đo ni đóng giày" cho UI của App Kotlin.**

* **Python Legacy (`api/index.py`):**
  - Endpoint: `POST /api/analyze`
  - Đóng gói: Trả về trực tiếp mảng các object kết quả (raw result) ngay sau khi Engine chạy xong mà không qua bước xử lý nào.
  - App Kotlin nếu dùng API này sẽ phải tự viết logic phức tạp để gom nhóm các mũi tiêm đã tiêm, chưa tiêm và phân tích các `status_tags` thành mã màu/trạng thái hiển thị.

* **Golang (`api/index.go`):**
  - Endpoint: `POST /api/analyze` (Có đi kèm Rate Limiting 50 req/phút).
  - Đóng gói: Go Engine không trả kết quả thô, mà có hẳn một vòng lặp map dữ liệu (Formatting for Android) để App Kotlin dễ dàng render UI nhất có thể.
  - Go Backend xử lý gom nhóm dữ liệu thành `missing_vaccines` và `administered_vaccines`.
  - **Quan trọng:** Go Backend tự động quy đổi các `StatusTags` của Engine (ví dụ: `due`, `overdue`, `error`, `warning`) thành **UI Status Enumerations** rất dễ sử dụng cho App Kotlin:
    - `DUE_NOW` (Cần tiêm ngay)
    - `OVERDUE` (Đã quá hạn)
    - `NEEDS_REVIEW` (Có cảnh báo/lỗi tương tác)
    - `COMPLETED` (Đã hoàn thành)
    - `DUE_LATER` (Chưa đến hạn)

## Tổng kết
Phiên bản Golang là bản nâng cấp toàn diện:
1. **Login:** Bổ sung cơ chế bảo vệ Distributed Lock với Redis để chống Race condition.
2. **API Output:** Dữ liệu trả về được chuẩn hóa thành Data Transfer Object (DTO) thân thiện và đóng gói các trạng thái chuẩn xác để App Kotlin hiển thị trực tiếp mà không cần tính toán lại.
