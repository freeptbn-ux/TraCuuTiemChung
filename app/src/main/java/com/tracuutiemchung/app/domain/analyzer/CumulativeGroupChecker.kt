package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Port of Python's check_cumulative_group_doses()
 * Applied to rule types: group_cumulative_unique_doses, group_cumulative_unique_doses_min_age
 */
fun checkCumulativeGroupDoses(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>,
): List<MissingItem> {
    val displayName = rule.displayName
    val dosesRequired = rule.requiredDoses ?: 0
    val validDosesCount = administeredRecords.size
    val lastDoseDate = administeredRecords.lastOrNull()?.date

    // 1. Series Complete Check
    if (validDosesCount >= dosesRequired) {
        return emptyList()
    }

    // 2. No doses - Check age
    if (administeredRecords.isEmpty()) {
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        
        // Use custom description for better UX in group context
        val description = if (ageStatus.statusTags.contains("too_young")) {
            "Chưa tiêm liều nào trong nhóm $displayName. Chờ đến ${ageStatus.message}."
        } else {
            "Chưa tiêm liều nào trong nhóm $displayName. ${ageStatus.message}."
        }

        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = description,
                earliestNextDoseDate = ageStatus.earliestDate,
                statusTags = ageStatus.statusTags
            )
        )
    }

    // 3. Next dose logic
    val nextDoseNumber = validDosesCount + 1
    val interval = rule.minIntervalDays.getOrNull(validDosesCount) ?: rule.intervalDays
    
    val dateByInterval = if (interval != null && lastDoseDate != null) {
        lastDoseDate.plusDays(interval.toLong())
    } else null
    
    val earliestDate = dateByInterval ?: analysisDate
    val isDue = analysisDate.isAfter(earliestDate) || analysisDate.isEqual(earliestDate)
    val statusTags = if (isDue) listOf("due") else listOf("scheduled")

    return listOf(
        MissingItem(
            vaccineNameForPopup = displayName,
            description = buildString {
                append("Đã tiêm $validDosesCount/$dosesRequired liều trong nhóm. Cần tiêm liều tiếp theo.")
                if (interval != null) {
                    append(" Cách liều trước tối thiểu $interval ngày.")
                }
            },
            earliestNextDoseDate = earliestDate,
            statusTags = statusTags
        )
    )
}
