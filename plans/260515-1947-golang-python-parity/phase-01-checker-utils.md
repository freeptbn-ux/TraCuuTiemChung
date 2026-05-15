# Phase 01: Checker Utilities & Rules Preprocessing
Status: ✅ Completed  
Dependencies: None (First phase)

## Objective
Port các utility functions thiếu từ Python sang Go và chuẩn hóa rules preprocessing.

## Phân tích chi tiết sự khác biệt

### 1.1. `NormalizeVaccineName` - Minor Fix
**Python** (`core/utils.py:4-22`):
```python
name = re.sub(r'\s*\(.*?\)\s*', '', name)  # Xóa ngoặc đơn → replace bằng empty string
```
**Go** (`pkg/analyzer/utils.go:18-20`):
```go
reParen := regexp.MustCompile(`\s*\(.*?\)\s*`)
name = reParen.ReplaceAllString(name, " ")  // ← Replace bằng SPACE, không phải empty string
```
**Go thêm** (Go có, Python không):
```go
reSuffix := regexp.MustCompile(`(?i)-(TCDV|TCMR)$`)
name = reSuffix.ReplaceAllString(name, "")
```

**Action:** 
- [ ] Fix Go: thay `" "` thành `""` trong regex ngoặc đơn
- [ ] Quyết định: Giữ hay bỏ regex `-TCDV/-TCMR` (keep vì không gây hại, nhưng đánh dấu là Go-only extension)

### 1.2. `get_administered_dose_records` - MISSING in Go
**Python** (`core/engine/checker_utils.py:29-38`):
```python
def get_administered_dose_records(names_norm_list, administered_map):
    records = []
    for norm_name in names_norm_list:
        records.extend(administered_map.get(norm_name, []))
    records.sort(key=lambda x: x["date"])
    return records
```
**Go** (`pkg/analyzer/engine.go:78-94`): Có `getMatchingRecords` nhưng khác ở chỗ nó **dedup by name** (dùng `seen` map).

**Vấn đề:** Python `extend` mà KHÔNG dedup by name. Nếu cùng 1 norm_name xuất hiện 2 lần trong `names_norm_list`, Python sẽ duplicate records. Go sẽ skip.

**Action:**
- [ ] Xác minh: `names_norm_list` có bao giờ chứa tên trùng không? → Thường KHÔNG, vì rules đã `list(set(...))` → Behavior giống nhau
- [ ] Kết luận: `getMatchingRecords` OK, nhưng cần document

### 1.3. `get_age_status_and_earliest_date` - MISSING in Go
**Python** (`core/engine/checker_utils.py:40-95`): 55 lines - Phức tạp!
- Check `min_age_months_overall_group`, `min_age_months_at_first_dose`, `min_age_months_overall`
- Check `min_age_weeks_overall_group`, `min_age_weeks_at_first_dose`
- Check `min_age_years_overall_group`, `min_age_years_at_first_dose`
- Check `min_age_days_overall_group`, `min_age_days_at_first_dose`
- Return tuple: `(status_message, earliest_date, status_tags)`
- Dùng `GRACE_PERIOD_DAYS = 0` (không ảnh hưởng nhưng cần giữ structure)

**Go:** Không có function tương đương. Logic inline trong `checkSingleSeries` nhưng rất đơn giản.

**Action:**
- [ ] Tạo file `pkg/analyzer/checker.go`
- [ ] Port `getAgeStatusAndEarliestDate()` trả về `(msg string, earliestDate *time.Time, tags []string)`
- [ ] Xử lý priority: days > weeks > months > years (theo Python)

### 1.4. `check_first_dose_age_validity` - MISSING in Go
**Python** (`core/engine/checker_utils.py:97-149`): 53 lines
- Kiểm tra tuổi mũi 1 có đúng quy định không
- Return True/False, append error vào `missing_items_list` nếu sai
- Check theo thứ tự: days → weeks → months → years

**Go:** Không có. Go không kiểm tra mũi 1 có hợp lệ về tuổi hay không.

**Action:**
- [ ] Port `checkFirstDoseAgeValidity()` trả về `(valid bool, errorResult *AnalysisResult)`
- [ ] Đảm bảo error messages giống Python ("Mũi 1 tiêm quá sớm...")

### 1.5. Rules Preprocessing - Partial
**Python** (`core/rules.py:34-78`):
```python
# Xử lý names_norm_group cho các group types
if "raw_names_members" in new_details:
    all_member_names_norm = []
    for member_key, raw_names_list in new_details["raw_names_members"].items():
        all_member_names_norm.extend(...)
    new_details["names_norm_group"] = list(set(all_member_names_norm))
```

**Go** (`pkg/analyzer/engine.go:33-68`): Chỉ xử lý `NamesNorm` cơ bản. Không tạo `NamesNormGroup`.

**Action:**
- [ ] Thêm field `NamesNormGroup []string` vào `VaccineRule` struct
- [ ] Trong `NewEngine`, xử lý `raw_names_members` → `NamesNormGroup`
- [ ] Xử lý courses → `NamesNorm` collective (khi rule-level `names_norm` empty)

### 1.6. AnalysisResult struct - Cần mở rộng
**Python** engine trả về dict với nhiều fields hơn Go.

**Action:**
- [ ] Đảm bảo `AnalysisResult` đủ fields cho tất cả use cases

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/analyzer/checker.go` | **CREATE** | Port `get_age_status_and_earliest_date`, `check_first_dose_age_validity` |
| `pkg/analyzer/utils.go` | MODIFY | Fix normalize regex (parentheses) |
| `pkg/analyzer/rules.go` | MODIFY | Add `NamesNormGroup` field |
| `pkg/analyzer/engine.go` | MODIFY | Add NamesNormGroup preprocessing |
| `pkg/analyzer/checker_test.go` | **CREATE** | Unit tests |

## Implementation Steps
1. [ ] Fix `NormalizeVaccineName` regex (parentheses → empty string)
2. [ ] Add `NamesNormGroup` field to `VaccineRule`
3. [ ] Create `pkg/analyzer/checker.go` with:
   - [ ] `getAgeStatusAndEarliestDate()`
   - [ ] `checkFirstDoseAgeValidity()`
   - [ ] Constants: `GracePeriodDays = 0`
4. [ ] Update `NewEngine` to process `NamesNormGroup` from `raw_names_members`
5. [ ] Update `NewEngine` to build collective `NamesNorm` from courses when rule-level empty
6. [ ] Write unit tests for all new functions

## Test Criteria

### Unit Tests (`pkg/analyzer/checker_test.go`)
- [ ] `TestGetAgeStatusAndEarliestDate_TooYoung_Months` - Trẻ 2 tháng, rule yêu cầu 6 tháng → `too_young`
- [ ] `TestGetAgeStatusAndEarliestDate_Eligible_Months` - Trẻ 12 tháng, rule yêu cầu 6 tháng → `eligible`
- [ ] `TestGetAgeStatusAndEarliestDate_TooYoung_Weeks` - Trẻ 4 tuần, rule yêu cầu 6 tuần → `too_young`
- [ ] `TestGetAgeStatusAndEarliestDate_TooYoung_Years` - Trẻ 3 tuổi, rule yêu cầu 4 tuổi → `too_young`
- [ ] `TestGetAgeStatusAndEarliestDate_NoDOB` - Không có DOB → `error_dob`
- [ ] `TestCheckFirstDoseAgeValidity_Valid` - Mũi 1 lúc 6 tháng, rule 2 tháng → true
- [ ] `TestCheckFirstDoseAgeValidity_TooEarly_Months` - Mũi 1 lúc 1 tháng, rule 2 tháng → false
- [ ] `TestCheckFirstDoseAgeValidity_TooEarly_Weeks` - Mũi 1 lúc 3 tuần, rule 6 tuần → false
- [ ] `TestNormalizeVaccineName_Parentheses` - Input `"Varivax (Thủy đậu)"` → `"varivax"` (không phải `"varivax  "`)
- [ ] `TestRulesPreprocessing_NamesNormGroup` - Load rules file, verify `NamesNormGroup` populated correctly

### Cách chạy test
```bash
cd vercel-backend
go test ./pkg/analyzer/ -run "TestGetAgeStatus|TestCheckFirstDose|TestNormalizeVaccineName_Parentheses|TestRulesPreprocessing" -v
```

## Notes
- `GRACE_PERIOD_DAYS = 0` trong Python nghĩa là không có grace period. Giữ constant cho tương lai.
- Priority check của Python: `min_age_months` → `min_age_years` → `min_age_weeks` → `min_age_days`. Go cần follow đúng thứ tự này.

---
Next Phase: [phase-02-series-logic.md](./phase-02-series-logic.md)
