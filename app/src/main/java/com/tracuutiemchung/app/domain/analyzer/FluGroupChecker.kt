package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Port of Python's check_flu_group()
 * Logic:
 * - First dose age < 9 years: 2 doses separated by 30 days.
 * - Then, annual booster (1 year from last dose).
 * - First dose age >= 9 years: 1 dose, then annual booster.
 */
fun checkFluGroup(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>,
): List<MissingItem> {
    val displayName = rule.displayName
    val lastDoseDate = administeredRecords.lastOrNull()?.date

    // 1. No doses
    if (administeredRecords.isEmpty()) {
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = "Chưa tiêm vắc-xin Cúm. ${ageStatus.message}.",
                earliestNextDoseDate = ageStatus.earliestDate,
                statusTags = ageStatus.statusTags
            )
        )
    }

    // 2. Initial series completion check
    val firstDoseDate = administeredRecords.first().date
    val ageAtFirstDose = dob?.let { VaccineDateUtils.getAgeAtDate(it, firstDoseDate) }
    
    // Trẻ dưới 9 tuổi lần đầu tiêm cần 2 mũi
    val needsTwoInitialDoses = ageAtFirstDose != null && ageAtFirstDose.totalYears < 9
    val initialSeriesComplete = if (needsTwoInitialDoses) {
        administeredRecords.size >= 2
    } else {
        administeredRecords.size >= 1
    }

    if (!initialSeriesComplete) {
        // Trường hợp trẻ < 9 tuổi mới tiêm 1 mũi
        val nextDate = firstDoseDate.plusDays(30)
        val isDue = analysisDate.isAfter(nextDate) || analysisDate.isEqual(nextDate)
        
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = "Trẻ dưới 9 tuổi lần đầu tiêm Cúm cần 2 mũi cách nhau tối thiểu 30 ngày. Đã tiêm 1 mũi.",
                earliestNextDoseDate = nextDate,
                statusTags = if (isDue) listOf("due") else listOf("scheduled")
            )
        )
    }

    // 3. Annual booster
    if (lastDoseDate != null) {
        val nextBoosterDate = VaccineDateUtils.addYears(lastDoseDate, 1)
        val isDue = analysisDate.isAfter(nextBoosterDate) || analysisDate.isEqual(nextBoosterDate)
        
        val statusTags = if (isDue) {
            listOf("due", "booster_due")
        } else {
            listOf("info", "booster_upcoming")
        }
        
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = "Đã hoàn thành phác đồ cơ bản. Cần tiêm nhắc lại hàng năm để duy trì miễn dịch.",
                earliestNextDoseDate = nextBoosterDate,
                statusTags = statusTags
            )
        )
    }

    return emptyList()
}
