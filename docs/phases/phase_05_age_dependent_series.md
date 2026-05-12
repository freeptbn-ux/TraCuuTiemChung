# Phase 05 — Age-Dependent Series

## Mục tiêu
Port `check_age_dependent_series()` cho các vaccine có phác đồ phụ thuộc tuổi mũi đầu. Áp dụng cho: **Prevenar13**, **Vaxneuvance**, **Synflorix**.

## Phụ thuộc
- Phase 04 hoàn thành (`SingleSeriesChecker` hoạt động).

## Tham chiếu Python: `series_checkers.py` → `check_age_dependent_series()`

### Luồng chính
1. Lấy `administered_records` theo `names_norm`.
2. Nếu **chưa tiêm mũi nào**:
   - Dùng `_get_age_status_and_earliest_date()` → gợi ý theo tuổi hiện tại.
   - Lấy `default_doses` từ `rules_by_age[0].doses_required`.
3. Nếu **thiếu DOB** → trả `error_dob`.
4. Tính `age_at_first_dose_months` bằng `get_age_at_date(dob, first_dose_date)`.
5. Kiểm tra tuổi overall cho nhóm (`min_age_*_overall`).
6. Duyệt `rules_by_age` → tìm rule phù hợp theo:
   - `min_age_at_first_dose_months <= age <= max_age_at_first_dose_months`.
7. Nếu **không tìm thấy rule phù hợp** → phân biệt:
   - Mũi đầu quá sớm (không hợp lệ overall) → `error_age_first_dose`.
   - Tuổi mũi đầu ngoài tất cả range → `error_no_matching_rule`.
8. Nếu **tìm thấy rule phù hợp** → tạo temp rule và delegate cho `checkSingleVaccineSeries()`.

### Ví dụ: Prevenar 13
```json
"rules_by_age": [
    {"max_age_at_first_dose_months": 6, "doses_required": 4, "min_interval_days": [null, 30, 30, 240]},
    {"min_age_at_first_dose_months": 7, "max_age_at_first_dose_months": 11, "doses_required": 3, "min_interval_days": [null, 30, 180]},
    {"min_age_at_first_dose_months": 12, "max_age_at_first_dose_months": 23, "doses_required": 2, "min_interval_days": [null, 60]},
    {"min_age_at_first_dose_months": 24, "doses_required": 1, "min_interval_days": [null]}
]
```
- Mũi đầu lúc 4 tháng → phác đồ 4 mũi.
- Mũi đầu lúc 14 tháng → phác đồ 2 mũi.
- Mũi đầu lúc 24 tháng → chỉ cần 1 mũi.

## Việc cần làm

### 1. Tạo `AgeDependentSeriesChecker.kt`

```kotlin
fun checkAgeDependentSeries(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
): List<MissingItem>
```

### 2. Logic chọn age rule

```kotlin
val ageAtFirstDose = getAgeAtDate(dob, firstDoseDate) ?: return errorItem(...)

val applicableRule = rule.rulesByAge.firstOrNull { ageRule ->
    val minM = ageRule.minAgeAtFirstDoseMonths
    val maxM = ageRule.maxAgeAtFirstDoseMonths
    (minM == null || ageAtFirstDose.totalMonths >= minM) &&
    (maxM == null || ageAtFirstDose.totalMonths <= maxM)
}
```

### 3. Delegate sang SingleSeriesChecker

```kotlin
if (applicableRule != null) {
    val tempRule = rule.copy(
        displayName = "${rule.displayName} (${applicableRule.regimenName ?: "..."})",
        requiredDoses = applicableRule.dosesRequired,
        minIntervalDays = applicableRule.minIntervalDays,
        doseSpecificRules = applicableRule.doseSpecificRules,
        boosterIntervalYears = applicableRule.boosterIntervalYears,
        // ...
    )
    return checkSingleVaccineSeries(ruleKey, tempRule, administeredRecords, dob, analysisDate, allRules)
}
```

### 4. Update engine dispatcher

Trong `VaccineAnalysisEngine`:
- Khi `rule.type == "age_dependent_series"` → gọi `checkAgeDependentSeries()`.

## Tests cần thêm — `AgeDependentSeriesCheckerTest.kt`

```
// --- Prevenar 13 ---
@Test prevenar13_firstDoseAt4Months_needs4Doses()
@Test prevenar13_firstDoseAt9Months_needs3Doses()
@Test prevenar13_firstDoseAt14Months_needs2Doses()
@Test prevenar13_firstDoseAt24Months_needs1Dose()
@Test prevenar13_3of4DosesCompleted_nextDoseInterval240Days()

// --- Synflorix special ---
@Test synflorix_7to11Months_dose2RequiresMinAbsAge12Months()
  - dose_specific_rules: {"2": {"min_absolute_age_months": 12}}

// --- Edge cases ---
@Test noDoses_returnsGuidanceBasedOnCurrentAge()
@Test noDob_returnsErrorDob()
@Test firstDose_tooEarlyOverall_returnsError()
@Test firstDose_outsideAllRanges_returnsNoMatchingRule()
@Test allDosesCompleted_returnsEmpty()
```

## File dự kiến chạm
- Tạo mới: `AgeDependentSeriesChecker.kt`
- Sửa: `VaccineAnalysisEngine.kt` — thêm dispatch cho `age_dependent_series`

## Tiêu chí xong
- [x] Prevenar13/Vaxneuvance/Synflorix trả đúng phác đồ theo tuổi mũi đầu.
- [x] Synflorix `dose_specific_rules` min_absolute_age_months hoạt động.
- [x] Edge cases (no doses, no DOB, too early) xử lý đúng.
- [x] SingleSeriesChecker tests vẫn pass.

## Thời gian ước tính
~2 giờ.
