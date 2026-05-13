# Phase 05: Complex Series Logic (Tuổi Động & Mũi Nhắc)

## Objective
Giải quyết các phần khó nhất của `series.py`: Phác đồ chọn theo tuổi mũi 1 (`rules_by_age`) và Mũi nhắc (Booster).

## Tasks
- [ ] Bổ sung logic Booster vào `CheckBasicSeries`: Nếu đã tiêm đủ `doses_required`, kiểm tra `booster_interval_years`. Nếu thoả mãn, báo "Cần tiêm nhắc lại".
- [ ] Viết hàm `CheckAgeDependentSeries`: Lấy ngày tiêm mũi 1, tính số tháng tuổi -> quét mảng `rules_by_age` để chọn ra sub-rule áp dụng -> Gọi lại `CheckBasicSeries` bằng sub-rule đó.
- [ ] Chuyển các hardcode "MVVAC vs MMR" (miễn dịch sởi bao phủ) và "VA-MENGOC-BC vs MenQuadfi" (reverse interaction).

## Files
- `engine/checkers/complex_series.go`
- `engine/checkers/age_dependent.go`

## Detailed Test Cases
### 1. Rules By Age (Selection Logic)
- **Scenario**: 
    - Rule A: First dose at 0-6 months -> 3 doses required.
    - Rule B: First dose at 7-11 months -> 2 doses required.
    - Patient 1: Dose 1 at 2 months. Expected: 3 total doses.
    - Patient 2: Dose 1 at 8 months. Expected: 2 total doses.
- **Verification**: `assert.Equal(t, 2, selectedRule.DosesRequired)` for Patient 2.

### 2. Booster Interval (Years)
- **Scenario**: 
    - Series finished at age 1. Booster required 4 years later.
    - Patient age 3.
- **Expected**: `MissingItem` exists but `Status` is "Chưa đến lịch" or date is in the future.
- **Expected (at age 5)**: `Status` becomes "Cần tiêm" (Due).

### 3. Interaction logic (Reverse dependency)
- **Scenario**: VA-MENGOC-BC vs MenQuadfi.
- **Rule**: If MenQuadfi (better) is already injected, do NOT recommend VA-MENGOC-BC (inferior).
- **Verification**: `assert.Empty(t, results)` for VA-MENGOC-BC if MenQuadfi is present in `AdministeredMap`.