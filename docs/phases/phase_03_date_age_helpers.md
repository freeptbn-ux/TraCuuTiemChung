# Phase 03 — Port Helper Ngày/Tuổi & Model Nội Bộ

## Mục tiêu
Tạo nền tảng utility tính ngày/tuổi có parity chính xác với Python. Tạo model nội bộ `MissingItem` để giữ parity output trước khi map sang `VaccineRecommendation`.

## Phụ thuộc
- Phase 02 hoàn thành (schema mở rộng xong).

## Tham chiếu Python

### `utils.py` — `VaccineAnalysisUtils.get_age_at_date(dob, target_date)`
```python
def get_age_at_date(dob, target_date):
    total_days = (target_date - dob).days
    if total_days < 0: return None, None, None

    years = target_date.year - dob.year
    if (target_date.month, target_date.day) < (dob.month, dob.day):
        years -= 1
    years = max(0, years)

    months_total = (target_date.year - dob.year) * 12 + (target_date.month - dob.month)
    if target_date.day < dob.day:
        months_total -= 1
    months_total = max(0, months_total)

    return months_total, total_days, years
```

### `rule_checker_utils.py` — `_add_months(source_date, months)`
```python
def _add_months(source_date, months):
    month = source_date.month - 1 + months
    year = source_date.year + month // 12
    month = month % 12 + 1
    day = source_date.day
    try:
        return date(year, month, day)
    except ValueError:
        last_day_of_month = calendar.monthrange(year, month)[1]
        return date(year, month, last_day_of_month)
```

### `rule_checker_utils.py` — `_add_years(source_date, years_to_add)`
```python
def _add_years(source_date, years_to_add):
    new_year = source_date.year + years_to_add
    try:
        return date(new_year, source_date.month, source_date.day)
    except ValueError:
        return date(new_year, source_date.month, 28)
```

## Việc cần làm

### 1. Tạo file `VaccineDateUtils.kt`

Đặt tại: `app/src/main/java/.../domain/analyzer/VaccineDateUtils.kt`

```kotlin
data class AgeAtDate(
    val totalMonths: Int,
    val totalDays: Int,
    val totalYears: Int,
)

fun getAgeAtDate(dob: LocalDate, atDate: LocalDate): AgeAtDate?
fun addMonths(source: LocalDate, months: Int): LocalDate
fun addYears(source: LocalDate, years: Int): LocalDate
```

Lưu ý quan trọng khi port:
- `addMonths(2024-01-31, 1)` phải = `2024-02-29` (leap year).
- `addMonths(2025-01-31, 1)` phải = `2025-02-28` (non-leap year).
- `addYears(2024-02-29, 1)` phải = `2025-02-28`.
- `getAgeAtDate` khi `atDate < dob` trả `null`.
- Months calculation: nếu `atDate.dayOfMonth < dob.dayOfMonth` → trừ 1 tháng.

### 2. Tạo file `AnalysisRuleUtils.kt`

Port `_get_age_status_and_earliest_date()` và `_check_first_dose_age_validity()`.

```kotlin
data class AgeStatus(
    val message: String,
    val earliestDate: LocalDate?,
    val statusTags: List<String>,
)

fun getAgeStatusAndEarliestDate(
    dob: LocalDate?,
    analysisDate: LocalDate,
    rule: VaccineRule,  // hoặc subset cần thiết
    displayName: String,
): AgeStatus

fun checkFirstDoseAgeValidity(
    dob: LocalDate?,
    firstDoseDate: LocalDate,
    rule: VaccineRule,  // hoặc subset fields min_age_*
    displayName: String,
): Pair<Boolean, MissingItem?>  // valid? + optional error item
```

Priority order cho min age (giống Python):
1. `min_age_days_at_first_dose`
2. `min_age_weeks_at_first_dose`
3. `min_age_months_at_first_dose` / `min_age_months_overall`
4. `min_age_years_at_first_dose`

### 3. Tạo model nội bộ `MissingItem`

```kotlin
data class MissingItem(
    val vaccineNameForPopup: String,
    val description: String,
    val earliestNextDoseDate: LocalDate?,
    val statusTags: List<String>,
)
```

### 4. Tạo hàm mapping `MissingItem` → `VaccineRecommendation`

```kotlin
fun MissingItem.toRecommendation(analysisDate: LocalDate): VaccineRecommendation
```

Mapping statusTags → RecommendationStatus:
- `"error_*"` → `NEEDS_REVIEW`
- `"too_young"`, `"too_old_*"` → `DUE_LATER`
- `"booster_upcoming"`, `"info"`, `"scheduled"` → `DUE_LATER`
- `"due"`, `"booster_due"` → `DUE_NOW`
- `"error_dob"` → `NOT_ENOUGH_DATA`

## Tests cần thêm — `VaccineDateUtilsTest.kt`

```
@Test addMonths_januaryThirtyFirstToFebruaryLeapYear()
  addMonths(2024-01-31, 1) == 2024-02-29

@Test addMonths_januaryThirtyFirstToFebruaryNonLeapYear()
  addMonths(2025-01-31, 1) == 2025-02-28

@Test addMonths_thirtyFirstToThirtyDayMonth()
  addMonths(2025-01-31, 3) == 2025-04-30

@Test addMonths_normalCase()
  addMonths(2025-03-15, 2) == 2025-05-15

@Test addYears_leapDayToNonLeapYear()
  addYears(2024-02-29, 1) == 2025-02-28

@Test addYears_normalCase()
  addYears(2025-06-15, 3) == 2028-06-15

@Test getAgeAtDate_beforeDobReturnsNull()
  getAgeAtDate(2025-01-01, 2024-12-31) == null

@Test getAgeAtDate_sameDateReturnsZero()
  getAgeAtDate(2025-01-01, 2025-01-01) == AgeAtDate(0, 0, 0)

@Test getAgeAtDate_thirteenMonthsOld()
  getAgeAtDate(2024-01-15, 2025-02-20) == AgeAtDate(months=13, days=402, years=1)

@Test getAgeAtDate_dayBeforeBirthday()
  getAgeAtDate(2024-03-15, 2025-03-14)
  → months=11, years=0 (chưa đủ 12 tháng và chưa đủ 1 tuổi)

@Test getAgeAtDate_exactlyOneYear()
  getAgeAtDate(2024-03-15, 2025-03-15) == AgeAtDate(months=12, days=365, years=1)
```

## Tests — `AnalysisRuleUtilsTest.kt`

```
@Test ageStatus_tooYoungByMonths()
@Test ageStatus_eligibleByMonths()
@Test ageStatus_noDob()
@Test firstDoseValidity_tooEarlyByWeeks()
@Test firstDoseValidity_validAge()
@Test missingItemToRecommendation_dueStatus()
@Test missingItemToRecommendation_errorStatus()
```

## File dự kiến chạm
- Tạo mới: `VaccineDateUtils.kt`
- Tạo mới: `AnalysisRuleUtils.kt`
- Tạo mới: `VaccineDateUtilsTest.kt`
- Tạo mới: `AnalysisRuleUtilsTest.kt`
- Chưa sửa `VaccineAnalysisEngine.kt` (để Phase 04).

## Tiêu chí xong
- [x] Tất cả date/age helper tests pass.
- [x] `addMonths` / `addYears` khớp Python edge cases.
- [x] `getAgeAtDate` khớp Python `VaccineAnalysisUtils.get_age_at_date`.
- [x] `MissingItem → VaccineRecommendation` mapping logic test pass.
- [x] Test cũ không bị break.

## Thời gian ước tính
~2 giờ.
