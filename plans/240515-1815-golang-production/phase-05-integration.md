# Phase 05: Full Integration & Testing
Status: ⬜ Pending
Dependencies: [Phase 04]

## Objective
Kiểm thử toàn diện và bàn giao.

## Requirements
- [ ] Kiểm thử E2E (Android -> Vercel -> Portal).
- [ ] Đảm bảo không có rò rỉ bộ nhớ (Memory leak) trong long-running.
- [ ] Tài liệu hóa các biến môi trường cần thiết.

## Implementation Steps
1. [ ] Chạy full test suite.
2. [ ] Kiểm tra logs trên Vercel Dashboard (giả lập bằng local log).
3. [ ] Cập nhật README.md với hướng dẫn deploy và config Redis.

## Test Criteria
- [ ] Search một số điện thoại thật -> Trả về kết quả chính xác.
- [ ] Đổi mật khẩu Portal -> Backend phải tự động detect và báo lỗi login fail (không được crash).

---
🏁 End of Plan
