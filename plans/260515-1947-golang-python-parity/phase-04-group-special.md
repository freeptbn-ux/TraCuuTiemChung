# Phase 04: Group Special Logic (MMR, Flu, Meningococcal)
Status: ✅ Completed  
Dependencies: Phase 02 (Series Logic)

## Objective
Rewrite `checkMMREquivalentGroup`, `checkFluGroup`, `checkMeningococcalACYWGroup` để match Python logic.

## Phân tích chi tiết sự khác biệt

### 4.1. `check_mmr_equivalent_group` - REWRITE

**Python** (`core/engine/group_special.py:12-163`): 152 lines  
**Go** (`pkg/analyzer/group_special.go:12-117`): 106 lines

#### Khác biệt #1: MVVAC interaction logic (Python chi tiết hơn)
**Python** (line 20-90):
```python
if mvvac_records:
    if num_mmr_group_doses == 0:
        # Chưa tiêm MMR, đã tiêm MVVAC
        earliest_by_mvvac_interval = last_mvvac_date + timedelta(days=84)
        # Kết hợp earliest_by_age và earliest_by_mvvac
        # Status: "due" nếu analysis_date >= earliest, else "info"+"scheduled"
        # Description rất chi tiết
        
    elif num_mmr_group_doses == 1:
        # Đã tiêm 1 mũi MMR sau MVVAC
        # CHECK: interval giữa MVVAC và MMR có đủ 84 ngày không?
        actual_interval = (first_mmr_date - last_mvvac_date).days
        if actual_interval < 84:
            # ⚠️ Warning: khoảng cách không đủ
            
        # Next: mũi 2 MMR sau 3 năm
        
    elif num_mmr_group_doses >= 2:
        return  # Completed
```

**Go**: Có logic tương tự nhưng:
- Thiếu interval violation warning (actual_interval < 84 days)
- Thiếu detailed descriptions
- Thiếu status tags `"info"`, `"scheduled"`, `"interval_violation_mvvac_mmr"`

#### Khác biệt #2: Regimen selection (Python dùng `regimens`, Go dùng `Regimens`)
**Python**: Dùng `check_single_vaccine_series` để delegate → hưởng booster, dose-specific, etc.
**Go**: Inline logic → thiếu features.

#### Khác biệt #3: "Chưa tiêm" message khi không có MVVAC
**Python** (line 92-102): Gọi `get_age_status_and_earliest_date`, sinh description chi tiết.
**Go**: Inline check nhưng thiếu detail.

#### Khác biệt #4: First dose age validity
**Python** (line 121-127): Check `check_first_dose_age_validity` cho group.
**Go**: Không check.

### 4.2. `check_flu_group` - BUG FIX + REWRITE

**Python** (`core/engine/group_special.py:204-290`): 87 lines  
**Go** (`pkg/analyzer/group_special.go:119-215`): 97 lines

#### Khác biệt #1: 🐛 BUG - Keyword matching method
**Python** (line 212-218):
```python
for norm_name, dose_list in administered_map.items():
    raw_name = dose_list[0]["raw_name"].lower()  # ← Match on RAW NAME
    for keyword in recognition_keywords:
        if keyword.lower() in raw_name:
            administered_records.extend(dose_list)
```
**Go** (line 131-148):
```go
for normName, recs := range administeredMap {
    for _, kw := range keywords {
        if strings.Contains(normName, kw) {  // ← Match on NORM NAME ❌
```
**Bug:** Python dùng `raw_name` (tên gốc) để match keywords. Go dùng `normName` (tên đã normalize). Khác biệt này có thể miss vaccines có tên gốc chứa keyword nhưng tên normalize không chứa.

**Ví dụ:** Vaccine tên gốc `"Vaxigrip Tetra (Cúm 4 chủng)"` → normalize bỏ ngoặc → `"vaxigrip tetra"`. Keyword `"cúm"` sẽ match raw_name nhưng KHÔNG match normName!

#### Khác biệt #2: First dose age validation
**Python** (line 233-238): Gọi `check_first_dose_age_validity`, nếu sai → error + return.
**Go**: Không check.

#### Khác biệt #3: Dedup logic
**Python** (line 253-258): Check `already_missing_second_dose` trước khi append mũi 2.
**Python** (line 280-285): Check `already_missing_annual_booster` trước khi append booster.
**Go**: Không dedup → có thể duplicate items.

#### Khác biệt #4: Description chi tiết
**Python**: "Cần mũi 2 (do <9 tuổi lần đầu tiêm) cách mũi 1 khoảng 4 tuần."
**Go**: "Tiêm mũi 2 khởi đầu"

### 4.3. `check_meningococcal_acyw_group` - REWRITE

**Python** (`core/engine/group_special.py:292-505`): 214 lines  
**Go** (`pkg/analyzer/group_special.go:217-357`): 141 lines

#### Khác biệt #1: Separate Menactra vs MenQuadfi handling
**Python**: Logic hoàn toàn tách biệt:
- Nếu có Menactra records → xử lý Menactra (delegate to `check_age_dependent_series`)
- Nếu có MenQuadfi records → xử lý MenQuadfi (inline with booster)
- Nếu không có records nào → gợi ý cả 2

**Go**: Logic hợp nhất thành 1 flow dựa trên member của mũi 1.

#### Khác biệt #2: `apply_interactions_and_append` function
**Python** (line 306-358): Function nội bộ xử lý tương tác với VA-MENGOC-BC và 6in1:
- Kiểm tra mỗi interaction riêng biệt
- Append warning items riêng biệt cho mỗi interaction
- Adjust `earliest_next_dose_date` dựa trên interaction constraints

**Go**: Interaction warnings merged vào description string, không tạo separate items.

#### Khác biệt #3: MenQuadfi booster logic
**Python** (line 440-455): Khi MenQuadfi đủ liều cơ bản → check booster config:
```python
booster_config = applicable_mq_rule.get("booster")
if booster_config:
    min_booster_age_m = booster_config["min_age_months"]
    min_interval_days = booster_config["min_interval_days_from_last"]
    earliest_by_age = add_months(dob, min_booster_age_m)
    earliest_by_interval = last_dose_date + timedelta(days=min_interval_days)
```
**Go**: Có tương tự nhưng logic tính toán slightly different.

#### Khác biệt #4: Khi chưa tiêm mũi nào - gợi ý cả 2 members
**Python** (line 457-504): 
- Gợi ý MenQuadfi (ưu tiên, từ 6 tuần)
- PLUS: Nếu trẻ đủ tuổi Menactra (9 tháng) → gợi ý thêm Menactra
- Tức là append 2 items!

**Go**: Chỉ gợi ý MenQuadfi.

#### Khác biệt #5: 60-day interval = add_months(2) trong Python
**Python**: `if min_interval == 60: earliest_date = add_months(last_dose_date, 2)` → calendar months.
**Go**: `AddDate(0, 0, 60)` → 60 days.

**Ví dụ:** `2025-01-31 + 2 months = 2025-03-31` vs `2025-01-31 + 60 days = 2025-04-01`. Khác 1 ngày!

## Implementation Steps

### 4.1: Rewrite `checkMMREquivalentGroup`
1. [x] Add MVVAC interval violation warning (actual_interval < 84)
2. [x] Add status tags: `"info"`, `"scheduled"`, `"interval_violation_mvvac_mmr"`
3. [x] Add first dose age validity check cho group
4. [x] Delegate to `checkSingleSeries` for regimen completion
5. [x] Add detailed Vietnamese descriptions
6. [x] Change return type to `[]AnalysisResult`

### 4.2: Fix `checkFluGroup`
1. [x] 🐛 FIX: Match keywords on **raw_name** instead of normName
2. [x] Add `checkFirstDoseAgeValidity` call
3. [x] Add dedup check (`already_missing_second_dose`, `already_missing_annual_booster`)
4. [x] Add detailed Vietnamese descriptions
5. [x] Change return type to `[]AnalysisResult`
6. [x] Need access to raw_name in VaccineRecord → verify `VaccineName` field is raw name

### 4.3: Rewrite `checkMeningococcalACYWGroup`
1. [x] Separate Menactra and MenQuadfi handling
2. [x] Port `apply_interactions_and_append` as method
3. [x] Add MenQuadfi booster logic with calendar months for 60-day interval
4. [x] Khi chưa tiêm: gợi ý cả MenQuadfi VÀ Menactra (2 items)
5. [x] Fix 60-day = 2 calendar months conversion
6. [x] Delegate Menactra to `checkAgeDependentSeries`
7. [x] Change return type to `[]AnalysisResult`

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/analyzer/group_special.go` | MAJOR REWRITE | Match all 3 Python functions |
| `pkg/analyzer/engine.go` | MODIFY | Update callers |
| `pkg/analyzer/group_special_test.go` | REWRITE | Comprehensive tests |

## Test Criteria

### MMR Tests
- [ ] `TestMMR_MVVAC_Then_0MMR` - MVVAC tiêm lúc 9 tháng, chưa tiêm MMR → recommend MMR sau 84 ngày
- [ ] `TestMMR_MVVAC_Then_1MMR_Good_Interval` - MVVAC→MMR cách 90 ngày → no warning, recommend mũi 2
- [ ] `TestMMR_MVVAC_Then_1MMR_Bad_Interval` - MVVAC→MMR cách 60 ngày → warning + recommend mũi 2
- [ ] `TestMMR_MVVAC_Then_2MMR` - MVVAC→2 MMR → completed, no result
- [ ] `TestMMR_No_MVVAC_NoDoses` - Chưa tiêm gì → recommend MMR mũi 1
- [ ] `TestMMR_No_MVVAC_1Dose_12m` - 1 MMR lúc 12 tháng → recommend mũi 2
- [ ] `TestMMR_FirstDose_TooEarly` - MMR lúc 7 tháng → error

### Flu Tests
- [ ] `TestFlu_NoDoses` - Trẻ 8 tháng → recommend mũi 1
- [ ] `TestFlu_1Dose_Under9` - 1 mũi lúc 2 tuổi → recommend mũi 2 (4 tuần)
- [ ] `TestFlu_2Doses_Under9` - 2 mũi → recommend annual booster
- [ ] `TestFlu_1Dose_Over9` - 1 mũi lúc 10 tuổi → recommend annual booster (no mũi 2)
- [ ] `TestFlu_Annual_Due` - Mũi cuối 13 tháng trước → `due`
- [ ] `TestFlu_Annual_Upcoming` - Mũi cuối 6 tháng trước → `booster_upcoming`
- [ ] `TestFlu_KeywordMatch_RawName` - Vaccine tên "Cúm A/B" phải match keyword "cúm"

### Meningococcal Tests
- [ ] `TestMeninACYW_NoDoses_Infant` - Trẻ 3 tháng → MenQuadfi gợi ý (3 mũi + nhắc)
- [ ] `TestMeninACYW_NoDoses_9m` - Trẻ 9 tháng → MenQuadfi + Menactra gợi ý
- [ ] `TestMeninACYW_Menactra_1Dose` - 1 mũi Menactra → recommend mũi 2 (delegate to age-dep)
- [ ] `TestMeninACYW_MenQuadfi_BasicComplete` - 3 mũi MenQuadfi → booster recommendation
- [ ] `TestMeninACYW_Interaction_MengocBC` - VA-Mengoc tiêm 1 tháng trước → interaction warning item
- [ ] `TestMeninACYW_Interaction_6in1` - 6in1 tiêm 20 ngày trước → interaction warning item
- [ ] `TestMeninACYW_60days_Calendar` - Verify 60-day interval = add_months(2)

### Cách chạy test
```bash
cd vercel-backend
go test ./pkg/analyzer/ -run "TestMMR|TestFlu|TestMeninACYW" -v
```

## Notes
- **BUG quan trọng:** Flu keyword matching cần dùng raw_name. Verify `models.VaccineRecord.VaccineName` stores raw name (it does based on parser logic).
- Meningococcal ACYW là phức tạp nhất trong phase này. Recommend: tách thành `handleMenactra()` và `handleMenQuadfi()` private functions.
- 60-day = 2 calendar months: Cần dùng `AddMonths(date, 2)` thay vì `AddDate(0, 0, 60)`.

---
Previous Phase: [phase-03-group-alternative.md](./phase-03-group-alternative.md)  
Next Phase: [phase-05-pneumo-cumulative-postproc.md](./phase-05-pneumo-cumulative-postproc.md)
