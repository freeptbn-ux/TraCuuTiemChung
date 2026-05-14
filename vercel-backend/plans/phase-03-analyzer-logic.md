# Phase 03: Analyzer Logic Migration

## Objective
Di chuyển engine phân tích vaccine cực kỳ phức tạp từ Python sang Go.

## Tasks:
- [x] Porting `vaccine_rules.json` sang folder `assets/`.
- [x] Implement logic chuẩn hóa tên vaccine (`normalize_vaccine_name`).
- [x] Viết bộ Engine xử lý các `RULE_TYPE` (Single, Age Dependent, v.v.).
- [x] Chú ý: Go `time` package xử lý ngày tháng rất nghiêm ngặt, cần viết utility chuyển đổi `DD/MM/YYYY`.

## 🧪 Testing for this phase:
- **File:** `pkg/analyzer/engine_test.go`
- **Mục tiêu:**
    - Tạo các bảng Test Cases (Table-driven tests) với Input là lịch tiêm giả định và Output là kết quả mong đợi.
    - Test các trường hợp đặc biệt: Tiêm thiếu mũi, tiêm sai khoảng cách, phế cầu tiêm xen kẽ.
- **Parity Test (So sánh song song):**
    - **Dữ liệu mẫu:** Sử dụng file `test/Gia-Han.html`.
    - **Quy trình:** 
        1. Parse file HTML này sang lịch tiêm (Administered list).
        2. Chạy Engine Go để dự đoán các mũi tiêm tương lai.
        3. So sánh kết quả với Output của code Python hiện tại (đã được export ra file `test/Gia-Han_expected.json`).
        4. Đảm bảo `next_dose` và `status` khớp 100%.
- **Lệnh chạy:** `go test -v ./pkg/analyzer/...`
