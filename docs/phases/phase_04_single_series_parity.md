# Phase 04 — Single Series Parity

## Mục tiêu
Rewrite logic xử lý `single_series`, `single_dose_min_age`, `single_series_min_age` trong engine Kotlin để khớp chính xác Python `check_single_vaccine_series()`.

## Phụ thuộc
- Phase 02 (schema mở rộng) + Phase 03 (date/age helpers) hoàn thành.

## Tham chiếu Python: `series_checkers.py` → `check_single_vaccine_series()`

### Luồng chính (tóm tắt)
1. Lấy tất cả `administered_records` theo `names_norm`.
2. `valid_doses_count = len(administered_records)` — đếm tất cả mũi, KHÔNG loại vì khoảng cách sai.
3. Nếu `valid_doses_count >= doses_required`:
   - Kiểm tra **booster** nếu có `booster_interval_years`.
   - Kiểm tra `booster_max_age_years` — nếu vượt tuổi → dừng, không gợi ý.
   - Nếu đến hạn booster → `status_tags = ["due", "booster_due"]`.
   - Nếu chưa đến hạn → `status_tags = ["info", "booster_upcoming"]`.
   - Return (series complete).
4. Logic **MVVAC ↔ MMR interaction** (chỉ áp dụng `rule_key == "MVVAC"`):
   - Nếu vaccine khác cung cấp measles protection và đã tiêm mũi sau 9 tháng tuổi → MVVAC không cần.
   - Trả `status_tags = ["info", "coverage_by_other"]`.
5. Logic **VA-MENGOC-BC** interaction (chỉ khi `rule_key == "VA-MENGOC-BC"` và trẻ ≥ 24 tháng).
6. Nếu `administered_records` rỗng:
   - Dùng `_get_age_status_and_earliest_date()` → trả status tuổi.
7. Kiểm tra `_check_first_dose_age_validity()`:
   - Nếu mũi đầu quá sớm → trả `error_age_first_dose` + restart.
8. Nếu `valid_doses_count < doses_required`:
   - Tính `next_dose_number = valid_doses_count + 1`.
   - Lấy interval từ `min_interval_days[next_dose_rule_idx]`.
   - Áp dụng `dose_specific_rules` nếu có (alternative age, min_absolute_age_months).
   - `earliest_next_dose_date = max(date_by_interval, date_by_abs_min_age, date_by_alt_age, analysis_date)`.

## Việc cần làm

### 1. Tạo `SingleSeriesChecker.kt`
Đặt tại: `app/src/main/java/.../domain/analyzer/SingleSeriesChecker.kt`

```kotlin
fun checkSingleVaccineSeries(
    ruleKey: String,
    rule: VaccineRule,  // hoặc config subset
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
): List<MissingItem>
```

### 2. Thay logic engine hiện tại

Trong `VaccineAnalysisEngine.analyze()`:
- Với các rule type `single_series`, `single_dose_min_age`, `single_series_min_age`:
  - Gọi `checkSingleVaccineSeries()` thay vì logic inline hiện tại.
  - Map kết quả `List<MissingItem>` → `VaccineRecommendation`.

### 3. Xử lý per-dose interval

```kotlin
// Thay vì: rule.intervalDays?.let { administered.last().date.plusDays(it.toLong()) }
// Port thành:
val nextDoseIndex = validDosesCount  // 0-indexed
val intervalForNext = rule.minIntervalDays.getOrNull(nextDoseIndex)
val dateByInterval = if (intervalForNext != null && lastDoseDate != null) {
    lastDoseDate.plusDays(intervalForNext.toLong())
} else null
```

### 4. Xử lý dose_specific_rules

```kotlin
val doseSpecificRule = rule.doseSpecificRules[nextDoseNumber.toString()]
// Nếu có alternative_min_age_years + alternative_max_age_years:
//   → dateByAltAge = addYears(dob, altMinAgeYears)
//   → lấy min giữa dateByInterval và dateByAltAge
// Nếu có min_absolute_age_months:
//   → dateByAbsMinAge = addMonths(dob, minAbsAgeMonths)
//   → earliestDate phải >= dateByAbsMinAge
```

### 5. Xử lý booster recurring

```kotlin
if (validDosesCount >= dosesRequired && rule.boosterIntervalYears != null) {
    val baseDateForBooster = administeredRecords.last().date
    val nextBoosterDue = addYears(baseDateForBooster, rule.boosterIntervalYears)
    
    // Check max age
    if (dob != null && rule.boosterMaxAgeYears != null) {
        val ageAtBooster = getAgeAtDate(dob, nextBoosterDue)
        if (ageAtBooster != null && ageAtBooster.totalYears >= rule.boosterMaxAgeYears) {
            return emptyList()  // Ngừng booster
        }
    }
    
    // Tạo MissingItem cho booster
}
```

### 6. MVVAC ↔ MMR interaction (optional cho phase này)

Chỉ áp dụng khi `ruleKey == "MVVAC"`. Kiểm tra trong `allRules` có rule nào
`providesMeaslesProtection == true` và đã có records → trả info coverage.

> **Lưu ý**: Logic này phức tạp và liên quan đến MMR group. Có thể defer sang Phase 06 nếu muốn.

## File dự kiến chạm
- Tạo mới: `SingleSeriesChecker.kt`
- Sửa: `VaccineAnalysisEngine.kt` — dùng dispatcher thay logic inline
- Sửa: `Models.kt` — thêm `statusTags: List<String>` vào `VaccineRecommendation`

## Tests cần thêm

### `SingleSeriesCheckerTest.kt`

```
// --- Booster tests ---
@Test completedSeries_noBooster_returnsEmpty()
@Test completedSeries_withBooster_due_returnsDueItem()
@Test completedSeries_withBooster_upcoming_returnsInfoItem()
@Test completedSeries_withBooster_maxAgeReached_returnsEmpty()

// --- First dose age validation ---
@Test noDoses_tooYoung_returnsWithEarliestDate()
@Test noDoses_eligible_returnsDue()
@Test noDoses_noDob_withMinAge_returnsNotEnoughData()
@Test firstDose_tooEarly_returnsError()

// --- Per-dose interval ---
@Test oneDoseOf3_usesCorrectIntervalForDose2()
  - Six_In_One: min_interval_days=[null,30,30,360]
  - Sau mũi 1, interval cho mũi 2 = 30 ngày

@Test twoDosesOf4_usesCorrectIntervalForDose3()
  - Sau mũi 2, interval cho mũi 3 = 30 ngày

@Test threeDosesOf4_usesCorrectIntervalForDose4()
  - Sau mũi 3, interval cho mũi 4 = 360 ngày

// --- Dose specific rules ---
@Test doseSpecific_alternativeAgeRange_usesEarlierDate()
@Test doseSpecific_minAbsoluteAgeMonths_constrainsDate()

// --- Backward compatibility ---
@Test existingGoldenTests_stillPass()
```

## Tiêu chí xong
- [x] `SingleSeriesChecker` xử lý đúng per-dose interval.
- [x] Booster recurring hoạt động với max age limit.
- [x] First dose age validation match Python.
- [x] Dose specific rules (alternative age, min absolute age) hoạt động.
- [x] **Tất cả test cũ (EngineTest + GoldenTest) pass hoặc update có lý do rõ ràng**.
- [x] Nếu GoldenTest cần update: ghi chú lý do trong commit message.

## Thời gian ước tính
~3-4 giờ.
