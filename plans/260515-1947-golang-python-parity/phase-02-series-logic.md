# Phase 02: Series Logic (Single + Age-Dependent)
Status: ✅ Completed  
Dependencies: Phase 01 (Checker Utilities)

## Objective
Rewrite `checkSingleSeries` và `checkAgeDependentSeries` để match logic Python chính xác. Đây là 2 function nền tảng được gọi bởi hầu hết các group checkers.

## Phân tích chi tiết sự khác biệt

### 2.1. `check_single_vaccine_series` - MAJOR REWRITE

**Python** (`core/engine/series.py:12-218`): 207 lines  
**Go** (`pkg/analyzer/engine.go:172-232`): 61 lines

#### Khác biệt #1: Booster logic (Python có, Go thiếu)
**Python** (line 25-61):
```python
if doses_required > 0 and valid_doses_count >= doses_required:
    booster_interval_years = rule_details.get("booster_interval_years")
    booster_applies_after_dose_num = rule_details.get("booster_after_dose_number", doses_required)
    booster_max_age_years = rule_details.get("booster_max_age_years")
    # ... check if booster is due or upcoming
```
**Go**: Khi đủ liều → `return nil`. Hoàn toàn bỏ qua booster!

**Impact:** Các vaccine như Td (uốn ván) cần nhắc lại mỗi 10 năm sẽ không được nhắc.

#### Khác biệt #2: MVVAC ↔ MMR Interaction (Python có, Go thiếu)
**Python** (line 64-100):
```python
if rule_key == "MVVAC" and all_vaccine_rules and dob:
    # Kiểm tra xem sởi đã được cover bởi MMR/Priorix chưa
    measles_covered_by_other_vaccine_info = None
    for other_key, other_rule_details_item in all_vaccine_rules.items():
        if other_rule_details_item.get("provides_measles_protection") or ...
```
**Go**: Hoàn toàn thiếu logic này.

**Impact:** Khi trẻ đã tiêm MMR (có cover sởi), app vẫn recommend tiêm MVVAC riêng → sai.

#### Khác biệt #3: VA-MENGOC-BC Reverse Interaction (Python có, Go thiếu)
**Python** (line 102-119):
```python
if rule_key == "VA-MENGOC-BC" and all_vaccine_rules and dob:
    age_months, _, _ = get_age_at_date(dob, analysis_date)
    if age_months >= 24:
        # Kiểm tra tương tác với MeningococcalACYW
```
**Go**: Thiếu.

#### Khác biệt #4: Khi chưa tiêm mũi nào (Python chi tiết hơn)
**Python** (line 121-131):
```python
if not administered_records:
    if doses_required > 0:
        age_status_msg, earliest_date, age_tags = get_age_status_and_earliest_date(...)
        desc = f"{rule_display_name} (Chưa tiêm - cần {doses_required} liều). {age_status_msg}"
```
**Go** (simplified):
```go
if numDoses == 0 {
    // Check max age → first dose logic
    earliestDate = AddMonths(e.DOB, rule.MinAgeMonthsAtFirstDose)
}
```
**Thiếu:** Go không gọi `getAgeStatusAndEarliestDate`, không kiểm tra `max_age_months_at_first_dose`, description thiếu chi tiết.

#### Khác biệt #5: Check first dose age validity (Python có, Go thiếu)
**Python** (line 133-140):
```python
if not check_first_dose_age_validity(dob, administered_records[0]["date"], ...):
    # Append error về mũi 1 không hợp lệ
    return
```
**Go**: Hoàn toàn thiếu validation này.

#### Khác biệt #6: Interval & dose-specific rules (Python chi tiết hơn)
**Python** (line 142-218): 
- Tính `date_by_interval` từ `min_interval_days`
- Tính `date_by_alt_age_years` từ `dose_specific_rules.alternative_min_age_years`
- Tính `date_by_abs_min_age_months` từ `dose_specific_rules.min_absolute_age_months`
- Lấy max của tất cả → `earliest_next_dose_date`
- Sinh description chi tiết: "Mũi X cách mũi Y tối thiểu Z ngày/tuần/tháng"

**Go** (line 210-217): Chỉ dùng `MinIntervalDays[numDoses]` đơn giản.

### 2.2. `check_age_dependent_series` - REWRITE

**Python** (`core/engine/series.py:220-320`): 100 lines  
**Go** (`pkg/analyzer/engine.go:234-297`): 64 lines

#### Khác biệt #1: Khi chưa tiêm mũi nào
**Python** (line 226-237): Gọi `get_age_status_and_earliest_date`, lấy `default_doses` từ rule.  
**Go** (line 240-247): Chỉ tìm applicable rule rồi check doses.

#### Khác biệt #2: Python delegate xuống `check_single_vaccine_series`
**Python** (line 305-319):
```python
temp_rule_for_check = {
    "display_name": ..., "names_norm": ...,
    "doses_required": applicable_age_rule["doses_required"],
    "min_interval_days": ..., "dose_specific_rules": ...,
    "booster_interval_years": ..., "booster_after_dose_number": ..., ...
}
check_single_vaccine_series(f"{rule_key}_age_specific", temp_rule_for_check, ...)
```
**Go**: Tự inline logic → thiếu booster, thiếu dose_specific_rules.

**Impact:** Code Go không hưởng được tất cả features của `checkSingleSeries` khi xử lý age-dependent.

## Implementation Steps

1. [ ] **Rewrite `checkSingleSeries`** → match Python `check_single_vaccine_series`:
   - [ ] Add booster logic khi đã đủ liều
   - [ ] Add MVVAC ↔ MMR interaction khi `ruleKey == "MVVAC"`
   - [ ] Add VA-MENGOC-BC reverse interaction
   - [ ] Add `checkFirstDoseAgeValidity` call khi đã có records
   - [ ] Add `getAgeStatusAndEarliestDate` khi chưa tiêm
   - [ ] Add `max_age_months_at_first_dose` check (quá tuổi tiêm)
   - [ ] Add interval description tiếng Việt (ngày/tuần/tháng/năm conversion)
   - [ ] Add dose-specific rules: `alternative_min_age_years`, `min_absolute_age_months`
   - [ ] Add `earliest_next_dose_date < analysis_date` → clamp to analysis_date

2. [ ] **Rewrite `checkAgeDependentSeries`** → match Python:
   - [ ] Khi chưa tiêm: gọi `getAgeStatusAndEarliestDate`
   - [ ] Khi không tìm được applicable rule: append error
   - [ ] Delegate xuống `checkSingleSeries` (tạo temp rule)
   - [ ] Add overall age check cho first dose validity

3. [ ] **Thay đổi function signature**: 
   - `checkSingleSeries` cần nhận `allRules map[string]VaccineRule` để check interactions
   - `checkSingleSeries` cần trả về `[]AnalysisResult` (list thay vì single) vì Python append nhiều items

4. [ ] Update callers trong `engine.go` `Analyze()` function

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/analyzer/engine.go` | MAJOR MODIFY | Rewrite `checkSingleSeries`, `checkAgeDependentSeries` |
| `pkg/analyzer/rules.go` | MODIFY | Add booster fields to VaccineRule if missing |
| `pkg/analyzer/series_test.go` | **CREATE** | Comprehensive tests |

## Test Criteria

### Unit Tests (`pkg/analyzer/series_test.go`)

#### Booster Tests
- [ ] `TestSingleSeries_Booster_Due` - 3/3 liều Td đã tiêm, mũi cuối 10+ năm → `booster_due`
- [ ] `TestSingleSeries_Booster_Upcoming` - 3/3 liều Td, mũi cuối 5 năm → `booster_upcoming`
- [ ] `TestSingleSeries_Booster_MaxAge` - Booster max age exceeded → no result
- [ ] `TestSingleSeries_Completed_NoBooster` - 3/3 liều BCG (no booster rule) → nil

#### MVVAC Interaction Tests
- [ ] `TestSingleSeries_MVVAC_CoveredByMMR` - Đã tiêm MMR lúc 12 tháng → MVVAC shows `coverage_by_other`
- [ ] `TestSingleSeries_MVVAC_NotCoveredMMR_TooYoung` - MMR tiêm lúc 7 tháng → MVVAC vẫn recommend

#### VA-MENGOC-BC Tests
- [ ] `TestSingleSeries_Mengoc_ReverseInteraction` - Trẻ >24 tháng, đã tiêm MenQuadfi → warning

#### First Dose Validity Tests
- [ ] `TestSingleSeries_FirstDose_TooEarly_Months` - BCG tiêm lúc 0 tháng, rule yêu cầu 2 tháng → error
- [ ] `TestSingleSeries_FirstDose_Valid` - BCG tiêm lúc 2 tháng, rule yêu cầu 2 tháng → continue

#### Age Status Tests (0 doses)
- [ ] `TestSingleSeries_NoDoses_Eligible` - Trẻ 12 tháng, rule min 6 tháng → `eligible`, earliest=analysis_date
- [ ] `TestSingleSeries_NoDoses_TooYoung` - Trẻ 1 tháng, rule min 6 tháng → `too_young`
- [ ] `TestSingleSeries_NoDoses_TooOld` - Trẻ 96 tháng, max_age 32 tuần → `too_old`

#### Dose Specific Rules Tests
- [ ] `TestSingleSeries_DoseSpecific_AltAge` - Mũi 2 DPT có `alternative_min_age_years: 4` → earliest = max(interval, age_4)
- [ ] `TestSingleSeries_DoseSpecific_AbsMinAge` - Mũi 3 cần `min_absolute_age_months: 12` → clamp

#### Age Dependent Tests
- [ ] `TestAgeDep_NoDoses_Status` - Prevenar13 chưa tiêm, trẻ 1 tháng → description chi tiết
- [ ] `TestAgeDep_FirstDose_At2Months` - Prevenar13 mũi 1 lúc 2 tháng → 4-dose schedule
- [ ] `TestAgeDep_FirstDose_At7Months` - Prevenar13 mũi 1 lúc 7 tháng → 3-dose schedule
- [ ] `TestAgeDep_Delegate_Booster` - Delegates correctly to checkSingleSeries with booster

### Cách chạy test
```bash
cd vercel-backend
go test ./pkg/analyzer/ -run "TestSingleSeries|TestAgeDep" -v
```

## Notes
- **CRITICAL**: Python `check_single_vaccine_series` trả về void, append vào list truyền vào. Go hiện trả về single `*AnalysisResult`. Cần đổi thành `[]AnalysisResult` vì Python có thể append NHIỀU items (warning + recommendation).
- `all_vaccine_rules` parameter cần được truyền vào để check cross-vaccine interactions.

---
Previous Phase: [phase-01-checker-utils.md](./phase-01-checker-utils.md)  
Next Phase: [phase-03-group-alternative.md](./phase-03-group-alternative.md)
