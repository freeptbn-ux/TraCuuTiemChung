# Phase 03: Group Alternative Logic
Status: ✅ Completed
Dependencies: Phase 02 (Series Logic)

## Objective
Rewrite `checkAlternativeCoursesMinAgeGroup` và `checkAlternativeCoursesAgeRangeGroup` để match Python logic, đặc biệt cho Rota, JE_Group, và HepA.

## Phân tích chi tiết sự khác biệt

### 3.1. `check_alternative_courses_group` (Rota) - MAJOR REWRITE

**Python** (`core/engine/group_alternative.py:10-141`): 132 lines  
**Go** (`pkg/analyzer/group_alternative.go:11-123`): 113 lines

#### Khác biệt #1: Python delegates to `check_single_vaccine_series` cho mỗi course
**Python**: Tạo `temp_course_rule` rồi gọi `check_single_vaccine_series()` → hưởng đầy đủ booster, age validity, interval descriptions.

**Go**: Tự tính interval inline → thiếu tất cả features ở trên.

#### Khác biệt #2: `max_age_months_for_completion_group` logic (Rota specific)
**Python** (line 104-123):
```python
if is_rota_rule_key and current_age_months is not None and \
   max_completion_age_months_rota is not None and \
   current_age_months > max_completion_age_months_rota:
    if course_specific_administered_records and current_course_messages:
        # Transform "due" messages → "too_old_to_complete"
```
**Go** (line 28-39): Chỉ check ở đầu function cho tất cả courses → không transform từng course.

#### Khác biệt #3: Multi-course tracking
**Python**: Track `any_course_completed` và `potential_incomplete_messages` riêng biệt. Nếu không course nào completed, append tất cả incomplete messages.

**Go**: Khi tìm được started course thì chỉ theo dõi 1 course.

#### Khác biệt #4: `max_age_months_to_start_first_dose_group` handling
**Python** (line 46-52): Nếu đã quá tuổi bắt đầu → return message "Đã qua X tháng tuổi, không còn chỉ định".

**Go** (line 43-49): Có nhưng logic message khác.

#### Khác biệt #5: Group first dose age validity
**Python** (line 63-82): Check `check_first_dose_age_validity` cho overall group trước khi check từng course.

**Go**: Không có check này.

### 3.2. `check_alternative_courses_age_range_group` (JE, HepA) - MAJOR REWRITE

**Python** (`core/engine/group_alternative.py:143-414`): 272 lines  
**Go** (`pkg/analyzer/group_alternative.go:125-247`): 123 lines

#### Khác biệt #1: JE_Group special logic (Python dành ~80 lines cho JE)
**Python** (line 148-236): Logic phức tạp cho 3 loại vaccine VNNB:
- **Jevax (inactivated)**: 3 mũi cơ bản + booster 3 năm/lần đến 15 tuổi
- **Imojev (live)**: 2 mũi
- **JEEV (live)**: 2 mũi

Cross-logic:
- Jevax → Imojev switch: Nếu đã 1 mũi Jevax, có thể chuyển sang 2 mũi Imojev
- Jevax → JEEV switch: Cảnh báo mixed
- Imojev trước Jevax: Error - không nên tiêm Jevax sau Imojev
- Jevax 3 mũi + 1 Imojev = completed
- Jevax 3 mũi completed → recommend Jevax booster HOẶC 1 mũi Imojev
- Jevax 1-2 mũi → recommend options để chuyển sang Imojev

**Go**: Chỉ có basic course matching, không có JE cross-logic.

#### Khác biệt #2: Per-course age range validation
**Python**: Mỗi course có `min_age_months_at_first_dose` VÀ `max_age_years_at_first_dose`. Python kiểm tra:
- Trẻ có đủ tuổi bắt đầu course này không?
- Trẻ có quá tuổi cho course này không?
- Mũi 1 đã tiêm có nằm trong range tuổi cho phép không?

**Go**: Chỉ check `MinAgeMonthsAtFirstDose` và `MaxAgeYearsAtFirstDose` cơ bản.

#### Khác biệt #3: "Chưa tiêm" với multiple options
**Python** (line 351-413): Khi chưa tiêm và có nhiều course phù hợp tuổi:
- HepA/JE: Liệt kê các lựa chọn có thể tiêm
- Other: Show "Lựa chọn có thể tiêm: Option1 Hoặc Option2"
- Tìm earliest date từ options

**Go**: Chỉ tìm bestCourse đầu tiên phù hợp → miss options khác.

#### Khác biệt #4: courses_to_skip mechanism
**Python**: Dùng `courses_to_skip` set để skip courses đã xử lý bởi JE special logic.

**Go**: Không có mechanism này.

## Implementation Steps

### 3.1: Rewrite `checkAlternativeCoursesMinAgeGroup`
1. [x] Change return type: `*AnalysisResult` → `[]AnalysisResult`
2. [x] Delegate mỗi course xuống `checkSingleSeries` (tạo temp rule)
3. [x] Track `anyCourseCompleted` + `potentialIncompleteMessages` riêng
4. [x] Add `max_age_months_to_start_first_dose_group` check
5. [x] Add `max_age_months_for_completion_group` + Rota transform logic
6. [x] Add group first dose age validity check
7. [x] Build options string cho "Chưa tiêm" message

### 3.2: Rewrite `checkAlternativeCoursesAgeRangeGroup`
1. [x] Change return type: `*AnalysisResult` → `[]AnalysisResult`
2. [x] Add JE_Group special logic:
   - [x] Jevax + Imojev completion check (3 Jevax + 1 Imojev = done)
   - [x] Jevax → JEEV mixed warning
   - [x] Imojev trước Jevax → error
   - [x] Jevax 3 mũi → booster recommendation
   - [x] Jevax 1-2 mũi → Imojev switch options
3. [x] Add per-course age range validation (min + max)
4. [x] Add multi-option "Chưa tiêm" message
5. [x] Add `courses_to_skip` set mechanism
6. [x] Delegate mỗi course xuống `checkSingleSeries`

### 3.3: Update callers
1. [x] Update `Analyze()` in `engine.go` to handle `[]AnalysisResult` return

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/analyzer/group_alternative.go` | MAJOR REWRITE | Match Python logic |
| `pkg/analyzer/engine.go` | MODIFY | Update callers for new return types |
| `pkg/analyzer/group_alternative_test.go` | REWRITE | Comprehensive tests |

## Test Criteria

### Rota Tests
- [x] `TestAltMinAge_Rota_NoDoses_Eligible` - Trẻ 2 tháng → show options (Rotarix, RotaTeq)
- [x] `TestAltMinAge_Rota_NoDoses_TooOld` - Trẻ 8 tháng, max_start 32 tuần → `too_old_to_start`
- [x] `TestAltMinAge_Rota_1Dose_Rotarix` - 1 mũi Rotarix → recommend mũi 2 Rotarix
- [x] `TestAltMinAge_Rota_2Doses_Rotarix` - 2 mũi Rotarix → completed, no result
- [x] `TestAltMinAge_Rota_TooOldToComplete` - 1 mũi Rotarix, trẻ 9 tháng, max_completion 8 tháng → `too_old_to_complete`
- [x] `TestAltMinAge_Rota_FirstDose_TooEarly` - Mũi 1 Rotarix lúc 4 tuần → error

### JE Tests
- [x] `TestAltAgeRange_JE_NoDoses` - Trẻ 12 tháng → show Jevax + Imojev options
- [x] `TestAltAgeRange_JE_Jevax3_Complete` - 3 mũi Jevax → recommend booster hoặc Imojev
- [x] `TestAltAgeRange_JE_Jevax3_Imojev1` - 3 Jevax + 1 Imojev → completed
- [x] `TestAltAgeRange_JE_Jevax1_SwitchOption` - 1 Jevax → show "chuyển sang Imojev" option
- [x] `TestAltAgeRange_JE_Jevax_Then_JEEV` - Jevax rồi JEEV → mixed warning
- [x] `TestAltAgeRange_JE_Imojev_Then_Jevax` - Imojev trước Jevax → error_interchange

### HepA Tests
- [x] `TestAltAgeRange_HepA_NoDoses_12m` - Trẻ 12 tháng → show eligible courses
- [x] `TestAltAgeRange_HepA_NoDoses_6m` - Trẻ 6 tháng → `too_young`
- [x] `TestAltAgeRange_HepA_1Dose_Havax` - 1 mũi Havax → recommend mũi 2

### Cách chạy test
```bash
cd vercel-backend
go test ./pkg/analyzer/ -run "TestAltMinAge|TestAltAgeRange" -v
```

## Notes
- **CRITICAL DESIGN CHANGE**: Cả 2 functions cần đổi từ trả `*AnalysisResult` sang `[]AnalysisResult` vì Python append nhiều items.
- JE_Group logic là phức tạp nhất. Recommend: tách thành private function `handleJEGroupSpecialLogic()`.
- Cần verify vaccine_rules.json có đầy đủ data cho courses (min_age, max_age, doses_required) cho JE, HepA.

---
Previous Phase: [phase-02-series-logic.md](./phase-02-series-logic.md)  
Next Phase: [phase-04-group-special.md](./phase-04-group-special.md)
