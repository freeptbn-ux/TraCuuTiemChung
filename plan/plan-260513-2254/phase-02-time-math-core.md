# Phase 02: Time Math Core

## Objective
Giải quyết triệt để sự khác biệt trong tính ngày tháng giữa Python (`dateutil.relativedelta` hoặc custom math) và Golang (`time.AddDate`). Cốt lõi để tính đúng mũi tương lai.

## Tasks
- [ ] Viết hàm `AddMonths(date time.Time, months int) time.Time`: Phải xử lý tay trường hợp cộng tháng mà ngày bị lố (Ví dụ: 31/01 + 1 tháng -> phải ra ngày cuối tháng 2 là 28 hoặc 29, chứ không phải 02/03 hay 03/03 như mặc định của Go).
- [ ] Viết hàm `AddYears(date time.Time, years int) time.Time` (Xử lý năm nhuận 29/02).
- [ ] Viết hàm `GetAgeAtDate(dob, target time.Time) (months, weeks, years int)` giống hệt file `checker_utils.py` của Python.
- [ ] Viết hàm chuẩn hóa tên `NormalizeVaccineName`.

## Files
- `engine/utils/time_math.go`
- `engine/utils/time_math_test.go`

## Detailed Test Cases
### 1. Month Overflow (The "31st" Rule)
- **Case A**: `AddMonths("2024-01-31", 1)` -> `2024-02-29` (Leap year).
- **Case B**: `AddMonths("2023-01-31", 1)` -> `2023-02-28` (Non-leap).
- **Case C**: `AddMonths("2024-08-31", 1)` -> `2024-09-30`.

### 2. Leap Year Anniversary
- **Case**: `AddYears("2024-02-29", 1)` -> `2025-02-28`.

### 3. Pediatric Age Calculation (CDSi Standard)
- **Case A (Boundary)**: DOB: `2024-01-15`, Target: `2024-03-14`.
    - **Expected**: 1 month, 28 days (NOT 2 months).
- **Case B (Exact)**: DOB: `2024-01-15`, Target: `2024-03-15`.
    - **Expected**: 2 months, 0 days.

### 4. Normalization Edge Cases
- **Case**: `NormalizeVaccineName("  Vắc-xin   6 trong 1  ")`.
    - **Expected**: `vac-xin 6 trong 1` (Remove accents, trim, collapse spaces).