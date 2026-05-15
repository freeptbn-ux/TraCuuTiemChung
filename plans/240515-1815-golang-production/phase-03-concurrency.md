# Phase 03: Concurrency & Locking
Status: ✅ Completed
Dependencies: [Phase 02]

## Objective
Gia cố cơ chế Login để an toàn khi có nhiều request đồng thời (Race condition).

## Requirements
### Functional
- [x] Sử dụng Distributed Lock (Redis SET NX) để đảm bảo tại một thời điểm chỉ có 1 worker login.
- [x] Cơ chế `WaitAndRetry`: Các request đến sau sẽ đợi cho đến khi Lock được release và lấy cookie mới từ Redis.

## Implementation Steps
1. [x] Tạo `pkg/portal/lock.go` xử lý logic Distributed Lock (Dùng `SET NX` với TTL).
2. [x] Bọc hàm `pc.Login()` bằng lock logic.
3. [x] Thêm timeout cho Lock (ví dụ 10s) để tránh treo hệ thống nếu login bị crash.
4. [x] Đảm bảo việc Release Lock phải an toàn (dùng Lua script check-and-delete để tránh xóa nhầm lock của process khác).

## Files to Create/Modify
- `pkg/portal/lock.go`
- `pkg/portal/client.go`

## Test Criteria
- [x] Chạy script `tests/load_test.go` giả lập 50 goroutines cùng gọi lookup. (Đã chạy concurrency_test.go)
- [x] Kiểm tra log: "Acquiring login lock..." chỉ được xuất hiện 1 lần.

---
Next Phase: [API Hardening & Middleware](./phase-04-hardening.md)
