# Phase 07: Processor Assembly & Sorting

## Objective
Ghép toàn bộ các hàm Checkers lại thành vòng lặp Main Engine giống file `processor.py` và chạy Hậu xử lý (Spacing & Sort) giống `post_processor.py`.

## Tasks
- [ ] Viết hàm `ProcessAllRules()` lặp qua toàn bộ `vaccine_rules.json`.
- [ ] Switch-case ánh xạ `rule_type` string sang đúng checker function (ví dụ: `"age_dependent_series"` -> gọi `CheckAgeDependentSeries`).
- [ ] Viết hàm `ApplySpacingAndSort()`: Giãn cách các vắc xin sống giảm độc lực (live vaccines) tối thiểu 28 ngày nếu không tiêm cùng ngày. Sort danh sách cuối cùng theo ngày `earliest_next_dose_date`.

## Files
- `engine/processor.go`
- `engine/post_processor.go`

## Detailed Test Cases
### 1. Live Vaccine Spacing (The 28-day rule)
- **Scenario**: 
    - Recommendation A: MMR (Live) - Earliest 2024-06-01.
    - Recommendation B: Varicella (Live) - Earliest 2024-06-01.
- **Expected After Post-Processing**: 
    - One stays at 2024-06-01.
    - The other is pushed to 2024-06-29 (28 days later).
- **Verification**: `assert.True(t, math.Abs(dateA - dateB) >= 28 days || dateA == dateB)` (If tiêm cùng ngày is allowed, otherwise just 28 days).

### 2. Output Sorting
- **Scenario**: Recommendations come out as: `Varicella (2024-10-01)`, `BCG (2024-01-01)`, `6-in-1 (2024-02-01)`.
- **Expected**: Final list must be ordered: `BCG`, `6-in-1`, `Varicella`.
- **Verification**: `assert.True(t, results[i].Date.Before(results[i+1].Date))`.

### 3. Rule Dispatcher Switch
- **Scenario**: Ensure all 4 `rule_type` from JSON are handled.
- **Verification**: If a rule has unknown type, it should return a clear error or log a warning instead of silent failure.