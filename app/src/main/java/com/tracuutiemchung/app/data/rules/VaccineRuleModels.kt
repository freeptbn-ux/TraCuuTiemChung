package com.tracuutiemchung.app.data.rules

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class DoseSpecificRule(
    @SerialName("alternative_min_age_years") val alternativeMinAgeYears: Int? = null,
    @SerialName("alternative_max_age_years") val alternativeMaxAgeYears: Int? = null,
    @SerialName("min_absolute_age_months") val minAbsoluteAgeMonths: Int? = null,
)

@Serializable
data class AgeBasedRegimen(
    @SerialName("regimen_name") val regimenName: String? = null,
    @SerialName("min_age_at_first_dose_months") val minAgeAtFirstDoseMonths: Int? = null,
    @SerialName("max_age_at_first_dose_months") val maxAgeAtFirstDoseMonths: Int? = null,
    @SerialName("min_age_weeks_at_first_dose") val minAgeWeeksAtFirstDose: Int? = null,
    @SerialName("doses_required") val dosesRequired: Int,
    @SerialName("min_interval_days") val minIntervalDays: List<Int?> = emptyList(),
    @SerialName("dose_specific_rules") val doseSpecificRules: Map<String, DoseSpecificRule> = emptyMap(),
    @SerialName("booster_interval_years") val boosterIntervalYears: Int? = null,
    @SerialName("booster_after_dose_number") val boosterAfterDoseNumber: Int? = null,
    @SerialName("booster_max_age_years") val boosterMaxAgeYears: Int? = null,
    @SerialName("booster") val booster: BoosterConfig? = null,
)

@Serializable
data class BoosterConfig(
    @SerialName("min_age_months") val minAgeMonths: Int? = null,
    @SerialName("min_interval_days_from_last") val minIntervalDaysFromLast: Int? = null,
    @SerialName("description") val description: String? = null,
)

@Serializable
data class CourseConfig(
    @SerialName("raw_names") val rawNames: List<String> = emptyList(),
    @SerialName("doses_required") val dosesRequired: Int,
    @SerialName("display") val display: String,
    @SerialName("min_age_months_at_first_dose") val minAgeMonthsAtFirstDose: Int? = null,
    @SerialName("max_age_years_at_first_dose") val maxAgeYearsAtFirstDose: Int? = null,
    @SerialName("min_interval_days") val minIntervalDays: List<Int?> = emptyList(),
    @SerialName("is_live") val isLive: Boolean = false,
    @SerialName("booster_interval_years") val boosterIntervalYears: Int? = null,
    @SerialName("booster_after_dose_number") val boosterAfterDoseNumber: Int? = null,
    @SerialName("booster_max_age_years") val boosterMaxAgeYears: Int? = null,
)

@Serializable
data class MemberConfig(
    @SerialName("raw_names") val rawNames: List<String> = emptyList(),
    @SerialName("display") val display: String,
    @SerialName("min_age_months_overall") val minAgeMonthsOverall: Int? = null,
    @SerialName("min_age_weeks_overall") val minAgeWeeksOverall: Int? = null,
    @SerialName("rules_by_age") val rulesByAge: List<AgeBasedRegimen> = emptyList(),
)

@Serializable
data class InteractionConfig(
    @SerialName("min_interval_days") val minIntervalDays: Int,
    @SerialName("applies_when_age_months_gte") val appliesWhenAgeMonthsGte: Int? = null,
    @SerialName("direction") val direction: String? = null,
    @SerialName("severity") val severity: String = "warning",
    @SerialName("message") val message: String,
)
