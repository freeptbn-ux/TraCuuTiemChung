package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Port of Python's check_single_vaccine_series()
 */
fun checkSingleVaccineSeries(
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

    // 2. MVVAC coverage check
    // If checking MVVAC but the patient has MMR or other measles-protecting vaccine, mark as complete.
    if (ruleKey == "MVVAC") {
        val coveredByOther = allAdministered.any { (key, records) ->
            key != "MVVAC" && records.isNotEmpty() && allRules[key]?.providesMeaslesProtection == true
        }
        if (coveredByOther) return emptyList()
    }

    // 1. Series Complete Check
    if (validDosesCount >= dosesRequired) {
        // Booster logic
        if (rule.boosterIntervalYears != null && lastDoseDate != null) {
            val nextBoosterDue = VaccineDateUtils.addYears(lastDoseDate, rule.boosterIntervalYears)
            
            // Check max age
            if (dob != null && rule.boosterMaxAgeYears != null) {
                val ageAtBooster = VaccineDateUtils.getAgeAtDate(dob, nextBoosterDue)
                if (ageAtBooster != null && ageAtBooster.totalYears >= rule.boosterMaxAgeYears) {
                    return emptyList() // Ngừng gợi ý booster nếu vượt tuổi
                }
            }

            val isDue = analysisDate.isAfter(nextBoosterDue) || analysisDate.isEqual(nextBoosterDue)
            val statusTags = if (isDue) {
                listOf("due", "booster_due")
            } else {
                listOf("info", "booster_upcoming")
            }

            return listOf(
                MissingItem(
                    vaccineNameForPopup = displayName,
                    description = "Đã hoàn thành phác đồ cơ bản. Đến hạn tiêm nhắc (booster).",
                    earliestNextDoseDate = nextBoosterDue,
                    statusTags = statusTags
                )
            )
        }
        return emptyList() // Đã xong và không có booster
    }

    // 2. MVVAC ↔ MMR interaction (Deferred to Phase 06 as per plan)

    // 3. No doses
    if (administeredRecords.isEmpty()) {
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        return listOf(
            MissingItem(
                vaccineNameForPopup = displayName,
                description = ageStatus.message,
                earliestNextDoseDate = ageStatus.earliestDate,
                statusTags = ageStatus.statusTags
            )
        )
    }

    // 4. First dose age validity
    val (isFirstDoseValid, errorItem) = AnalysisRuleUtils.checkFirstDoseAgeValidity(
        dob, administeredRecords.first().date, rule, displayName
    )
    if (!isFirstDoseValid && errorItem != null) {
        return listOf(errorItem)
    }

    // 5. Next dose logic
    val nextDoseNumber = validDosesCount + 1
    val nextDoseIndex = validDosesCount // 0-indexed

    // Interval from min_interval_days
    val intervalForNext = rule.minIntervalDays.getOrNull(nextDoseIndex) ?: rule.intervalDays
    val dateByInterval = if (intervalForNext != null && lastDoseDate != null) {
        lastDoseDate.plusDays(intervalForNext.toLong())
    } else null

    // Start with dateByInterval or analysisDate
    var earliestDate = dateByInterval ?: analysisDate

    // Dose Specific Rules
    val doseSpecificRule = rule.doseSpecificRules[nextDoseNumber.toString()]
    if (doseSpecificRule != null) {
        // alternative_min_age_years
        if (dob != null && doseSpecificRule.alternativeMinAgeYears != null) {
            val dateByAltAge = VaccineDateUtils.addYears(dob, doseSpecificRule.alternativeMinAgeYears)
            // Python: min(date_by_interval, date_by_alt_age) if date_by_interval exists
            if (dateByInterval != null && dateByAltAge.isBefore(dateByInterval)) {
                earliestDate = dateByAltAge
            }
        }
        
        // min_absolute_age_months
        if (dob != null && doseSpecificRule.minAbsoluteAgeMonths != null) {
            val dateByAbsMinAge = VaccineDateUtils.addMonths(dob, doseSpecificRule.minAbsoluteAgeMonths)
            // earliest_next_dose_date = max(earliest_next_dose_date, date_by_abs_min_age)
            if (earliestDate.isBefore(dateByAbsMinAge)) {
                earliestDate = dateByAbsMinAge
            }
        }
    }

    val isDue = analysisDate.isAfter(earliestDate) || analysisDate.isEqual(earliestDate)
    val statusTags = if (isDue) listOf("due") else listOf("scheduled")
    
    return listOf(
        MissingItem(
            vaccineNameForPopup = displayName,
            description = buildString {
                append("Đã tiêm $validDosesCount/$dosesRequired mũi. Cần tiêm mũi tiếp theo.")
                if (dateByInterval != null) {
                    val minInterval = rule.minIntervalDays.getOrNull(nextDoseNumber - 1) ?: rule.intervalDays
                    append(" Mũi $nextDoseNumber cách mũi trước tối thiểu $minInterval ngày.")
                }
            },
            earliestNextDoseDate = earliestDate,
            statusTags = statusTags
        )
    )
}
