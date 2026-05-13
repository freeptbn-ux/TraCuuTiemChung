# Plan: Phân rã Logic Tính Mũi Vắc-xin Tương lai (Python -> Go)

## Overview
Chia nhỏ engine tính toán ngày tiêm vắc-xin tương lai thành nhiều lớp. Do đặc thù tính phức tạp của phác đồ tiêm chủng (tuổi tối thiểu, khoảng cách tối thiểu, logic bù mũi, tiêm xen kẽ), quá trình port code sang Golang phải được test strict bằng dữ liệu thực tế (Minh Khoi patient data) để tránh "off-by-one errors" (lệch ngày).

## Tech Stack
- **Language**: Golang (1.21+)
- **Testing**: `testify/assert` (cho Unit Test)
- **Parser**: thư viện `PuerkitoBio/goquery` (để parse file HTML trong file test).

## Phases

| Phase | Name | Status |
|------|------|--------|
| 01 | Setup Models | Complete |
| 02 | Time Math Core | Complete (Bugfixing) |
| 03 | Rules Loader | Complete |
| 04 | Basic Series Logic | Complete (Bugfixing) |
| 05 | Complex Series Logic | Complete (Bugfixing) |
| 06 | Special Groups (Pneumo/Flu) | Complete (Bugfixing) |
| 07 | Processor Assembly | Complete (Bugfixing) |
| 08 | HTML Parity Testing | Pending |

> [!IMPORTANT]
> **Update 2026-05-13**: Đã bổ sung "Detailed Test Cases" vào từng file phase. Cần thực hiện kiểm tra lại toàn bộ logic dựa trên các case này để sửa lỗi "lệch ngày" và logic bù mũi.