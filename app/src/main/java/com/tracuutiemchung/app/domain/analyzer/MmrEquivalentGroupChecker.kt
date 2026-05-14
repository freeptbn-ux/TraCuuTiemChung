package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.AgeBasedRegimen
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Port of Python's check_mmr_equivalent_group()
 * Logic for MVVAC <-> MMR interaction and age-based regimens.
 */
fun checkMmrEquivalentGroup(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>
): List<MissingItem> {
    val displayName = rule.displayName

    // 1. Check MVVAC (Measles Single Dose) interaction
    // MVVAC is the key for Measles vaccine in vaccine_rules.json
    val mvvacRecords = allAdministered["MVVAC"].orEmpty()

    if (mvvacRecords.isNotEmpty()) {
        val latestMvvac = mvvacRecords.last()
        // Filter MMR doses given AFTER MVVAC
        val mmrAfterMvvac = administeredRecords.filter { it.date.isAfter(latestMvvac.date) }

        // Path: 0/1/≥2 MMR doses sau khi có MVVAC, interval 84 ngày.
        if (mmrAfterMvvac.size >= 2) {
            return emptyList() // Completed
        }

        if (mmrAfterMvvac.isEmpty()) {
            // Next dose is 1st MMR after MVVAC
            // Earliest date = max(latestMvvac + 84 days, dob + 12 months)
            val dateByInterval = latestMvvac.date.plusDays(84)
            val dateByMinAge = if (dob != null) VaccineDateUtils.addMonths(dob, 12) else null
            
            val earliestDate = listOfNotNull(dateByInterval, dateByMinAge).maxOrNull() ?: dateByInterval
            val isDue = analysisDate.isAfter(earliestDate) || analysisDate.isEqual(earliestDate)

            return listOf(
                MissingItem(
                    vaccineNameForPopup = displayName,
                    description = "Đã tiêm Sởi (MVVAC). Cần tiêm MMR mũi tiếp theo (cách MVVAC tối thiểu 84 ngày).",
                    earliestNextDoseDate = earliestDate,
                    statusTags = if (isDue) listOf("due") else listOf("scheduled")
                )
            )
        } else {
            // mmrAfterMvvac.size == 1
            // Next dose is 2nd MMR after MVVAC (effectively 3rd measles dose)
            val lastMmrDate = mmrAfterMvvac.last().date
            // Standard interval for 3rd dose in 3-dose regimen is 1095 days (3 years)
            val earliestNext = lastMmrDate.plusDays(1095)
            
            val isDue = analysisDate.isAfter(earliestNext) || analysisDate.isEqual(earliestNext)
            return listOf(
                MissingItem(
                    vaccineNameForPopup = displayName,
                    description = "Đã tiêm Sởi (MVVAC) và 1 mũi MMR. Cần tiêm MMR mũi tiếp theo.",
                    earliestNextDoseDate = earliestNext,
                    statusTags = if (isDue) listOf("due") else listOf("scheduled")
                )
            )
        }
    }

    // 2. No MVVAC: Standard MMR logic (Age Dependent)
    // Use 'regimens' field for MMR_Group
    return checkAgeDependentSeriesWithRegimens(
        ruleKey, rule, rule.regimens, administeredRecords, dob, analysisDate, allRules, allAdministered
    )
}

/**
 * Helper to handle age-dependent series with a specific list of regimens.
 */
private fun checkAgeDependentSeriesWithRegimens(
    ruleKey: String,
    rule: VaccineRule,
    regimens: List<AgeBasedRegimen>,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>
): List<MissingItem> {
    val displayName = rule.displayName

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

    if (administeredRecords.isEmpty()) {
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        val defaultDoses = regimens.firstOrNull()?.dosesRequired ?: rule.requiredDoses ?: 0
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = if (defaultDoses > 0) "${ageStatus.message}. Phác đồ dự kiến $defaultDoses mũi." else ageStatus.message,
                earliestNextDoseDate = ageStatus.earliestDate,
                statusTags = ageStatus.statusTags
            )
        )
    }

    val firstDoseDate = administeredRecords.first().date
    val ageAtFirstDose = VaccineDateUtils.getAgeAtDate(dob, firstDoseDate) ?: return emptyList()
    val totalMonthsAtFirstDose = ageAtFirstDose.totalMonths

    val applicableRegimen = regimens.firstOrNull { reg ->
        val minM = reg.minAgeAtFirstDoseMonths
        val maxM = reg.maxAgeAtFirstDoseMonths
        (minM == null || totalMonthsAtFirstDose >= minM) &&
        (maxM == null || totalMonthsAtFirstDose <= maxM)
    }

    if (applicableRegimen == null) {
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = "$displayName - Không tìm thấy phác đồ cho tuổi mũi đầu ($totalMonthsAtFirstDose tháng).",
                earliestNextDoseDate = null,
                statusTags = listOf("error_no_matching_rule")
            )
        )
    }

    val tempRule = rule.copy(
        displayName = if (applicableRegimen.regimenName != null) "$displayName (${applicableRegimen.regimenName})" else displayName,
        requiredDoses = applicableRegimen.dosesRequired,
        minIntervalDays = applicableRegimen.minIntervalDays,
        doseSpecificRules = applicableRegimen.doseSpecificRules,
        boosterIntervalYears = applicableRegimen.boosterIntervalYears,
        boosterAfterDoseNumber = applicableRegimen.boosterAfterDoseNumber,
        boosterMaxAgeYears = applicableRegimen.boosterMaxAgeYears
    )

    return checkSingleVaccineSeries(ruleKey, tempRule, administeredRecords, dob, analysisDate, allRules, allAdministered)
}
