# Phase 04: Basic Series Logic (Mũi Tương Lai Cơ Bản)

## Objective
Chuyển logic của hàm `check_single_vaccine_series` (nằm trong `series.py`) nhưng CẮT BỎ phần Booster và Tương tác. Chỉ tập trung tính: "Đã tiêm X mũi, tính ngày tiêm mũi X+1".

## Tasks
- [ ] Tạo hàm `CheckBasicSeries(rule, adminMap, dob, analysisDate) []MissingItem`.
- [ ] Đếm số mũi hợp lệ (Dùng hàm kiểm tra mũi 1 có hợp lệ tuổi không).
- [ ] Tính `remaining_doses = doses_required - valid_doses`.
- [ ] Trích xuất `min_interval_days[next_dose_idx]` và cộng vào ngày tiêm mũi cuối để ra `earliest_next_dose_date`.
- [ ] So sánh `date_by_interval` với `date_by_abs_min_age_months`. Lấy ngày Max. Trả về item cần tiêm.

## Files
- `engine/checkers/basic_series.go`
- `engine/checkers/basic_series_test.go`

## Detailed Test Cases
### 1. Interval vs Age Constraint (The "MAX" rule)
- **Scenario**: 
    - Rule: Min age 2mo, Min interval 1mo.
    - Patient: DOB: 2024-01-01. Dose 1: 2024-02-15 (at 1.5mo - INVALID).
    - **Expected**: 
        - Dose 1 marked as invalid (too young).
        - Next Dose 1 scheduled for 2024-03-01 (at exactly 2mo).
- **Scenario 2**: Dose 1 valid at 2024-03-01.
    - **Expected**: Dose 2 scheduled for 2024-04-01 (D1 + 1mo).

### 2. Valid Dose Counting
- **Scenario**: Patient has 2 doses. Dose 1 is valid, Dose 2 is too close to Dose 1 (e.g. 15 days).
- **Expected**: `valid_doses` = 1. `earliest_next_dose_date` must be calculated from Dose 1, NOT Dose 2.

### 3. Completion Logic
- **Scenario**: `doses_required` = 3. Patient has 3 valid doses.
- **Expected**: `MissingItem` list should NOT contain this vaccine (unless it has a booster).