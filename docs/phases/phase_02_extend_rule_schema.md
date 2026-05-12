# Phase 02 — Mở rộng Rule Schema Kotlin

## Mục tiêu
Parse đầy đủ dữ liệu từ `vaccine_rules.json` để engine có nguyên liệu xử lý. **Chỉ mở rộng schema, chưa đổi behavior engine.**

## Phụ thuộc
- Phase 01 hoàn thành (biết gap rõ ràng).

## Việc cần làm

### 1. Mở rộng `RawVaccineRule` trong `VaccineRuleRepository.kt`

Thêm các field JSON chưa parse:

```kotlin
// Trong RawVaccineRule, thêm:
@SerialName("dose_specific_rules") val doseSpecificRules: Map<String, JsonObject> = emptyMap(),
@SerialName("booster_interval_years") val boosterIntervalYears: Int? = null,
@SerialName("booster_after_dose_number") val boosterAfterDoseNumber: Int? = null,
@SerialName("booster_max_age_years") val boosterMaxAgeYears: Int? = null,
@SerialName("rules_by_age") val rulesByAge: List<JsonObject> = emptyList(),
@SerialName("regimens") val regimens: List<JsonObject> = emptyList(),
@SerialName("members") val members: Map<String, JsonObject> = emptyMap(),
@SerialName("courses") val courses: List<JsonObject> = emptyList(),
@SerialName("interactions") val interactions: Map<String, JsonObject> = emptyMap(),
@SerialName("provides_measles_protection_group") val providesMeaslesProtectionGroup: Boolean = false,
@SerialName("max_age_months_to_start_first_dose_group") val maxAgeMonthsToStartGroup: Int? = null,
@SerialName("max_age_months_for_completion_group") val maxAgeMonthsForCompletionGroup: Int? = null,
@SerialName("initial_series_interval_days") val initialSeriesIntervalDays: Int? = null,
@SerialName("is_live") val isLive: Boolean = false,
@SerialName("min_age_months_overall") val minAgeMonthsOverall: Int? = null,
@SerialName("min_age_years_overall") val minAgeYearsOverall: Int? = null,
@SerialName("min_age_days_at_first_dose") val minAgeDaysAtFirstDose: Int? = null,
```

### 2. Mở rộng `VaccineRule` data class

Thêm field mới vào public model:

```kotlin
data class VaccineRule(
    // ... existing fields giữ nguyên ...
    val minIntervalDays: List<Int?> = emptyList(),       // full interval list
    val doseSpecificRules: Map<String, DoseSpecificRule> = emptyMap(),
    val boosterIntervalYears: Int? = null,
    val boosterAfterDoseNumber: Int? = null,
    val boosterMaxAgeYears: Int? = null,
    val rulesByAge: List<AgeBasedRegimen> = emptyList(),  // for age_dependent_series
    val regimens: List<AgeBasedRegimen> = emptyList(),    // for mmr_equivalent_group
    val members: Map<String, MemberConfig> = emptyMap(),  // for meningococcal_acyw_group
    val courses: List<CourseConfig> = emptyList(),         // for alternative/age_range groups
    val interactions: Map<String, InteractionConfig> = emptyMap(),
    val providesMeaslesProtection: Boolean = false,
    val isLive: Boolean = false,
    val maxAgeMonthsToStart: Int? = null,
    val maxAgeMonthsForCompletion: Int? = null,
    val initialSeriesIntervalDays: Int? = null,
)
```

### 3. Tạo data class phụ trợ

```kotlin
@Serializable
data class DoseSpecificRule(
    val alternativeMinAgeYears: Int? = null,
    val alternativeMaxAgeYears: Int? = null,
    val minAbsoluteAgeMonths: Int? = null,
)

@Serializable
data class AgeBasedRegimen(
    val regimenName: String? = null,
    val minAgeAtFirstDoseMonths: Int? = null,
    val maxAgeAtFirstDoseMonths: Int? = null,
    val minAgeWeeksAtFirstDose: Int? = null,
    val dosesRequired: Int,
    val minIntervalDays: List<Int?> = emptyList(),
    val doseSpecificRules: Map<String, DoseSpecificRule> = emptyMap(),
    val boosterIntervalYears: Int? = null,
    val boosterAfterDoseNumber: Int? = null,
    val boosterMaxAgeYears: Int? = null,
    val booster: BoosterConfig? = null,    // cho MenQuadfi
)

@Serializable
data class BoosterConfig(
    val minAgeMonths: Int? = null,
    val minIntervalDaysFromLast: Int? = null,
    val description: String? = null,
)

@Serializable
data class CourseConfig(
    val rawNames: List<String> = emptyList(),
    val dosesRequired: Int,
    val display: String,
    val minAgeMonthsAtFirstDose: Int? = null,
    val maxAgeYearsAtFirstDose: Int? = null,
    val minIntervalDays: List<Int?> = emptyList(),
    val isLive: Boolean = false,
    val boosterIntervalYears: Int? = null,
    val boosterAfterDoseNumber: Int? = null,
    val boosterMaxAgeYears: Int? = null,
)

@Serializable
data class MemberConfig(
    val rawNames: List<String> = emptyList(),
    val display: String,
    val minAgeMonthsOverall: Int? = null,
    val minAgeWeeksOverall: Int? = null,
    val rulesByAge: List<AgeBasedRegimen> = emptyList(),
)

@Serializable
data class InteractionConfig(
    val minIntervalDays: Int,
    val appliesWhenAgeMonthsGte: Int? = null,
    val direction: String? = null,
    val severity: String = "warning",
    val message: String,
)
```

### 4. Update `toVaccineRule()` conversion

Cập nhật hàm `RawVaccineRule.toVaccineRule(key)` để map các field mới:
- `minIntervalDays` giữ nguyên list thay vì lấy phần tử đầu.
- Parse `rulesByAge`, `regimens`, `courses`, `members` từ JsonObject sang typed class.
- Giữ backward-compatible: `intervalDays` vẫn lấy first non-null element từ list.

### 5. Update `VaccineNameNormalizer`

Mở rộng alias từ `courses` và `members`:
- Với mỗi `CourseConfig`, thêm `rawNames` vào alias list.
- Với mỗi `MemberConfig`, thêm `rawNames` vào alias list.

## File dự kiến chạm
- `app/src/main/java/.../data/rules/VaccineRuleRepository.kt` — chính
- Có thể tạo file mới: `VaccineRuleModels.kt` cho các data class phụ trợ

## Tests cần thêm/sửa trong `VaccineRuleRepositoryTest.kt`

```
@Test loadRules_parsesMinIntervalDaysAsList()
  - Six_In_One_Combined phải có minIntervalDays = [null, 30, 30, 360]

@Test loadRules_parsesRulesByAge()
  - Prevenar13 phải có 4 age rules với doses_required [4, 3, 2, 1]

@Test loadRules_parsesCourses()  
  - JE_Group phải có 3 courses (Imojev, JEEV, Jevax)
  - Rota phải có 3 courses

@Test loadRules_parsesMembers()
  - MeningococcalACYW_Group phải có 2 members (MENACTRA, MENQUADFI)

@Test loadRules_parsesRegimens()
  - MMR_Group phải có 3 regimens

@Test loadRules_parsesBoosterConfig()
  - Jevax course trong JE_Group phải có boosterIntervalYears = 3

@Test loadRules_parsesInteractions()
  - MeningococcalACYW_Group phải có interaction VA-MENGOC-BC

@Test normalizer_mapsCourseNamesToGroupKey()
  - "Imojev" → JE_Group
  - "Rota Teq" → Rota
  - "HAVAX" → HepA
```

## Tiêu chí xong
- [x] Tất cả field JSON quan trọng được parse sang typed Kotlin.
- [x] Tests schema pass (không thay đổi behavior engine).
- [x] Normalizer vẫn map đúng + thêm mapping từ courses/members.
- [x] Test cũ không bị break.

## Thời gian ước tính
~2-3 giờ (schema design + tests + refactor converter).
