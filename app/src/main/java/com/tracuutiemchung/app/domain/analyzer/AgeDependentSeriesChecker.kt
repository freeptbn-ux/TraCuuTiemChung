package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Port of Python's check_age_dependent_series()
 */
fun checkAgeDependentSeries(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>,
): List<MissingItem> {
    val displayName = rule.displayName

    // 1. No DOB
    if (dob == null) {
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = "$displayName - Thiếu ngày sinh, không thể kiểm tra phác đồ.",
                earliestNextDoseDate = null,
                statusTags = listOf("error_dob")
            )
        )
    }

    // 2. No doses
    if (administeredRecords.isEmpty()) {
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        // Default doses if no doses administered (usually the first regimen)
        // Python: default_doses = rules_by_age[0].doses_required if available
        val defaultDoses = rule.rulesByAge.firstOrNull()?.dosesRequired ?: rule.requiredDoses ?: 0
        
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = if (defaultDoses > 0) "${ageStatus.message}. Phác đồ dự kiến $defaultDoses mũi." else ageStatus.message,
                earliestNextDoseDate = ageStatus.earliestDate,
                statusTags = ageStatus.statusTags
            )
        )
    }

    // 3. Overall age check (min_age_*_overall)
    val firstDoseDate = administeredRecords.first().date
    val (isFirstDoseValid, errorItem) = AnalysisRuleUtils.checkFirstDoseAgeValidity(
        dob, firstDoseDate, rule, displayName
    )
    if (!isFirstDoseValid && errorItem != null) {
        return listOf(errorItem)
    }

    // 4. Find matching rule by age at first dose
    val ageAtFirstDose = VaccineDateUtils.getAgeAtDate(dob, firstDoseDate) ?: return listOf(
        MissingItem(
            vaccineNameForPopup = displayName,
            description = "$displayName - Lỗi tính tuổi cho mũi đầu.",
            earliestNextDoseDate = null,
            statusTags = listOf("error_age_calculation")
        )
    )

    val totalMonthsAtFirstDose = ageAtFirstDose.totalMonths
    val applicableRule = rule.rulesByAge.firstOrNull { ageRule ->
        val minM = ageRule.minAgeAtFirstDoseMonths
        val maxM = ageRule.maxAgeAtFirstDoseMonths
        (minM == null || totalMonthsAtFirstDose >= minM) &&
        (maxM == null || totalMonthsAtFirstDose <= maxM)
    }

    if (applicableRule == null) {
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = "$displayName - Không tìm thấy phác đồ phù hợp cho độ tuổi tiêm mũi đầu ($totalMonthsAtFirstDose tháng).",
                earliestNextDoseDate = null,
                statusTags = listOf("error_no_matching_rule")
            )
        )
    }

    // 5. Delegate to SingleSeriesChecker with temp rule
    val tempRule = rule.copy(
        displayName = if (applicableRule.regimenName != null) "$displayName (${applicableRule.regimenName})" else displayName,
        requiredDoses = applicableRule.dosesRequired,
        minIntervalDays = applicableRule.minIntervalDays,
        doseSpecificRules = applicableRule.doseSpecificRules,
        boosterIntervalYears = applicableRule.boosterIntervalYears,
        boosterAfterDoseNumber = applicableRule.boosterAfterDoseNumber,
        boosterMaxAgeYears = applicableRule.boosterMaxAgeYears
    )

    return checkSingleVaccineSeries(ruleKey, tempRule, administeredRecords, dob, analysisDate, allRules, allAdministered)
}
