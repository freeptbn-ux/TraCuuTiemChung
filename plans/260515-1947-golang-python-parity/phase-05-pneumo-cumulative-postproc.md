# Phase 05: Pneumococcal + Cumulative + Post-processor
Status: ✅ Completed  
Dependencies: Phase 02 (Series Logic), Phase 01 (Checker Utils)

## Objective
Rewrite `processPneumoRules`, implement missing `check_cumulative_group_doses`, và port `apply_spacing_and_sort` post-processor.

## Phân tích chi tiết sự khác biệt

### 5.1. `processPneumoRules` - REWRITE

**Python** (`core/engine/processor.py:30-81`): 52 lines  
**Go** (`pkg/analyzer/group_pneumo.go:12-114`): 103 lines

#### Khác biệt #1: Mixed series warning (logic giống, message khác)
**Python** (line 48-53):
```python
mixed_names = [pneumo_rules[k].get('display_name') for k in active_series_keys]
current_missing_items.append({
    "description": f"Cảnh báo: Đã ghi nhận tiêm xen kẽ các loại phế cầu ({' và '.join(mixed_names)})...",
    "vaccine_name_for_popup": "Phế cầu (nhiều loại)"
})
```
**Go** (line 39-43): Dùng `"Phế cầu (Mixed)"` và joined by `" + "`.

#### Khác biệt #2: `patient_age_years >= 2` logic (Python có, Go khác)
**Python** (line 56-80): Logic phức tạp khi trẻ > 2 tuổi:
- Nếu primary series < 3 mũi → gợi ý Pneumovax23 "Có thể tiêm 1 mũi để hoàn thành"
- Nếu primary series == 3 mũi → check if mũi 4 needed, gợi ý Pneumovax23 "thay thế cho mũi 4"

**Go** (line 96-111): Logic đơn giản hơn:
- Chỉ check `monthsNow >= 24` cho Pneumovax23
- Không có primary series count logic

#### Khác biệt #3: Skip logic
**Python**: Dùng `pneumo_rules_to_skip` set, skip tất cả khi Pneumovax23 đã tiêm HOẶC mixed series.
**Go**: Logic khác - show tất cả 3 PCV khi chưa bắt đầu loại nào.

#### Khác biệt #4: Delegate to `check_age_dependent_series`
**Python** (line 70): Gọi `check_age_dependent_series(primary_key, ...)` để check nếu mũi 4 cần thiết.
**Go**: Gọi `checkAgeDependentSeries` nhưng logic khác.

### 5.2. `check_cumulative_group_doses` - MISSING in Go

**Python** (`core/engine/group_special.py:166-202`): 37 lines

Xử lý `group_cumulative_unique_doses` và `group_cumulative_unique_doses_min_age`:
- Count tổng số liều unique đã tiêm trong group
- Check first dose age validity
- Nếu thiếu liều → recommend
- Dùng `required_total_unique_doses` từ rule

**Go**: Không có handler. `Analyze()` function không xử lý `RuleTypeGroupCumulativeUnique`.

**Action:**
- [ ] Tạo function `checkCumulativeGroupDoses`
- [ ] Thêm vào switch case trong `Analyze()`

### 5.3. Post-processor `apply_spacing_and_sort` - MISSING in Go

**Python** (`core/engine/post_processor.py:1-93`): 93 lines

Chức năng:
1. **Live vaccine spacing**: Nếu mũi thiếu là vaccine sống VÀ mũi gần nhất đã tiêm cũng là vaccine sống → cách tối thiểu 28 ngày
2. **General spacing**: Tất cả mũi thiếu cách mũi gần nhất đã tiêm tối thiểu 14 ngày
3. **Sort**: Sort kết quả theo `earliest_next_dose_date`, rồi theo `description`
4. **Helper functions**: `get_vaccine_live_status_by_norm_name()`, `is_missing_item_live()`

**Go**: Không có post-processor nào. Kết quả trả về unsorted, không có spacing logic.

**Quan trọng:** Post-processor cần biết vaccine nào là "live" (sống). Thông tin này nằm trong `VaccineRule.IsLive` và `Course.IsLive`.

## Implementation Steps

### 5.1: Rewrite `processPneumoRules`
1. [x] Port skip logic: Pneumovax23 đã tiêm → skip tất cả PCV
2. [x] Port `patient_age_years >= 2` logic:
   - [x] Primary < 3 mũi → Pneumovax23 alternative completion
   - [x] Primary == 3 mũi → check mũi 4, Pneumovax23 alternative booster
3. [x] Match Python's `pneumo_rules_to_skip` set behavior
4. [x] Match description messages
5. [x] Delegate to `checkAgeDependentSeries` (updated Phase 02 version)
6. [x] Fix mixed series: dùng `display_name` thay vì hardcoded names

### 5.2: Implement `checkCumulativeGroupDoses`
1. [x] Create function matching Python `check_cumulative_group_doses`
2. [x] Count total unique doses across group names
3. [x] Check first dose age validity
4. [x] Return appropriate results
5. [x] Add to switch case in `Analyze()`:
   ```go
   case RuleTypeGroupCumulativeUnique:
       res := e.checkCumulativeGroupDoses(ruleKey, rule, administeredMap)
       results = append(results, res...)
   ```

### 5.3: Port Post-processor
1. [x] Create `pkg/analyzer/post_processor.go`
2. [x] Port `getVaccineLiveStatusByNormName()` - check if norm_name belongs to a live vaccine
3. [x] Port `isMissingItemLive()` - check if a result item is for a live vaccine
4. [x] Port `ApplySpacingAndSort()`:
   - [x] Find last live vaccine date from administered records
   - [x] Find last overall vaccine date
   - [x] For each result: apply 14-day general spacing
   - [x] For live results: apply 28-day live-to-live spacing
   - [x] Clamp to analysis_date
   - [x] Sort results
5. [x] Call `ApplySpacingAndSort()` at the end of `Analyze()` in engine.go

### 5.4: Wire everything together
1. [x] Add `RuleTypeGroupCumulativeUnique` to Analyze() switch
2. [x] Call post-processor at end of Analyze()
3. [x] Need to pass `administeredMap` to post-processor

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/analyzer/group_pneumo.go` | REWRITE | Match Python pneumo logic |
| `pkg/analyzer/group_cumulative.go` | **CREATE** | New: cumulative group doses |
| `pkg/analyzer/post_processor.go` | **CREATE** | New: spacing, sorting, live vaccine logic |
| `pkg/analyzer/engine.go` | MODIFY | Add cumulative handler, call post-processor |
| `pkg/analyzer/post_processor_test.go` | **CREATE** | Post-processor tests |
| `pkg/analyzer/group_pneumo_test.go` | **CREATE** | Pneumo tests |

## Test Criteria

### Pneumococcal Tests
- [ ] `TestPneumo_NoDoses` - Chưa tiêm → gợi ý cả 3 PCV (Prevenar, Vaxneuvance, Synflorix)
- [ ] `TestPneumo_Mixed_Prevenar_Synflorix` - Tiêm cả 2 → mixed warning
- [ ] `TestPneumo_Pneumovax23_Done` - Pneumovax23 đã tiêm → skip tất cả
- [ ] `TestPneumo_Prevenar_2Doses_Over2Years` - 2 mũi Prevenar, trẻ > 2 tuổi → Pneumovax23 alternative
- [ ] `TestPneumo_Prevenar_3Doses_Over2Years` - 3 mũi Prevenar, trẻ > 2 tuổi → check mũi 4 + Pneumovax23 option
- [ ] `TestPneumo_Prevenar_4Doses` - 4 mũi Prevenar → completed, skip tất cả

### Cumulative Tests
- [ ] `TestCumulative_NoDoses` - 0/3 liều DPT combined → recommend
- [ ] `TestCumulative_Partial` - 2/3 liều → "Cần thêm 1 liều"
- [ ] `TestCumulative_Complete` - 3/3 liều → no result
- [ ] `TestCumulative_FirstDose_TooEarly` - Mũi 1 quá sớm → error

### Post-processor Tests
- [ ] `TestPostProc_LiveVaccineSpacing` - 2 results: live vaccine, last administered live 15 days ago → adjust to 28 days
- [ ] `TestPostProc_GeneralSpacing` - Result has date 5 days after last administered → adjust to 14 days
- [ ] `TestPostProc_NoAdmin` - No administered vaccines → no spacing adjustment
- [ ] `TestPostProc_Sort` - 3 results with dates → sorted by date ascending
- [ ] `TestPostProc_NilDate_SortLast` - Result with nil date → sorted last
- [ ] `TestIsLive_Imojev` - Imojev is live → true
- [ ] `TestIsLive_Jevax` - Jevax is inactivated → false
- [ ] `TestIsLive_JEGroup_Description` - JE item mentioning "imojev" in description → true

### Cách chạy test
```bash
cd vercel-backend
go test ./pkg/analyzer/ -run "TestPneumo|TestCumulative|TestPostProc|TestIsLive" -v
```

## Notes
- Post-processor cần access vào `administeredMap` (for finding last dose dates) VÀ `Rules` (for is_live lookups). Consider adding it as Engine method.
- `is_missing_item_live` trong Python khá hacky (checks description string for "imojev", "jevax" etc). Port as-is for parity.
- `VaccineRule.IsLive` field already exists in Go rules.go. Verify it's populated from JSON.

---
Previous Phase: [phase-04-group-special.md](./phase-04-group-special.md)  
Next Phase: [phase-06-integration-test.md](./phase-06-integration-test.md)
