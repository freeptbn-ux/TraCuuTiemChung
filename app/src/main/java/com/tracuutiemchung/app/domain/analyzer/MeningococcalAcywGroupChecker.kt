package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.VaccineRule
import com.tracuutiemchung.app.data.rules.MemberConfig
import com.tracuutiemchung.app.data.rules.AgeBasedRegimen
import java.time.LocalDate
import java.time.temporal.ChronoUnit

/**
 * Port of Python's check_meningococcal_acyw_group()
 * Handles complex logic for Menactra and MenQuadfi, including interactions.
 */
fun checkMeningococcalAcywGroup(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>,
): List<MissingItem> {
    val displayName = rule.displayName

    // 1. No doses - Suggest MenQuadfi by default as it can start earlier (6 weeks)
    if (administeredRecords.isEmpty()) {
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        
        // Interaction Check for suggesting first dose
        val interactionWarning = checkInteractions(null, allAdministered, dob, analysisDate, rule)
        val description = if (interactionWarning != null) {
            "${ageStatus.message}. $interactionWarning"
        } else {
            ageStatus.message
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

    // 2. Identify active member and regimen
    val firstRecord = administeredRecords.first()
    val memberKey = identifyMember(firstRecord, rule.members)
    val memberConfig = rule.members[memberKey] ?: return emptyList()
    
    val ageAtFirst = dob?.let { VaccineDateUtils.getAgeAtDate(it, firstRecord.date) }
    val regimen = findApplicableRegimen(ageAtFirst, memberConfig.rulesByAge) ?: return emptyList()

    val dosesRequired = regimen.dosesRequired
    val dosesTaken = administeredRecords.size
    val lastDoseDate = administeredRecords.last().date

    // 3. Series Completion & Booster Check
    if (dosesTaken >= dosesRequired) {
        val booster = regimen.booster
        if (booster != null && dob != null) {
            val dateByAge = VaccineDateUtils.addMonths(dob, booster.minAgeMonths ?: 0)
            val dateByInterval = booster.minIntervalDaysFromLast?.let { lastDoseDate.plusDays(it.toLong()) } ?: lastDoseDate
            
            val earliestBoosterDate = if (dateByInterval.isAfter(dateByAge)) dateByInterval else dateByAge
            val isDue = analysisDate.isAfter(earliestBoosterDate) || analysisDate.isEqual(earliestBoosterDate)
            
            // Check if booster already taken? 
            // Usually, if dosesTaken > dosesRequired, the extra doses are considered boosters.
            if (dosesTaken > dosesRequired) return emptyList()

            return listOf(
                MissingItem(
                    vaccineNameForPopup = displayName,
                    description = booster.description ?: "Cần tiêm mũi nhắc (booster).",
                    earliestNextDoseDate = earliestBoosterDate,
                    statusTags = if (isDue) listOf("due", "booster_due") else listOf("info", "booster_upcoming")
                )
            )
        }
        return emptyList()
    }

    // 4. Next dose logic
    val nextDoseIndex = dosesTaken
    val interval = regimen.minIntervalDays.getOrNull(nextDoseIndex) ?: 30 // Fallback to 30 days
    val earliestByInterval = lastDoseDate.plusDays(interval.toLong())
    
    var earliestDate = earliestByInterval
    val interactionWarning = checkInteractions(lastDoseDate, allAdministered, dob, analysisDate, rule)
    
    val isDue = analysisDate.isAfter(earliestDate) || analysisDate.isEqual(earliestDate)
    val description = buildString {
        append("Đã tiêm $dosesTaken/$dosesRequired mũi ${memberConfig.display}. Cần tiêm mũi tiếp theo.")
        if (interactionWarning != null) {
            append(" Lưu ý: $interactionWarning")
        }
    }

    return listOf(
        MissingItem(
            vaccineNameForPopup = displayName,
            description = description,
            earliestNextDoseDate = earliestDate,
            statusTags = if (isDue) listOf("due") else listOf("scheduled")
        )
    )
}

private fun identifyMember(parsedRecord: ParsedRecord, members: Map<String, MemberConfig>): String? {
    val name = parsedRecord.source.vaccineName.lowercase()
    for ((key, config) in members) {
        if (key.lowercase() == name || config.rawNames.any { it.lowercase() == name }) {
            return key
        }
        // Partial match
        if (name.contains(key.lowercase()) || config.rawNames.any { name.contains(it.lowercase()) }) {
            return key
        }
    }
    return null
}

private fun findApplicableRegimen(ageAtFirst: AgeAtDate?, regimens: List<AgeBasedRegimen>): AgeBasedRegimen? {
    if (ageAtFirst == null) return regimens.lastOrNull()
    
    return regimens.firstOrNull { regimen ->
        val minMatch = regimen.minAgeAtFirstDoseMonths?.let { ageAtFirst.totalMonths >= it } ?: true
        val maxMatch = regimen.maxAgeAtFirstDoseMonths?.let { ageAtFirst.totalMonths <= it } ?: true
        val minWeekMatch = regimen.minAgeWeeksAtFirstDose?.let { ageAtFirst.totalDays >= it * 7 } ?: true
        minMatch && maxMatch && minWeekMatch
    } ?: regimens.lastOrNull()
}

private fun checkInteractions(
    lastDoseDate: LocalDate?,
    allAdministered: Map<String, List<ParsedRecord>>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    rule: VaccineRule
): String? {
    val warnings = mutableListOf<String>()
    
    // 1. VA-MENGOC-BC interaction
    val mengocBcRecords = allAdministered["VA-MENGOC-BC"] ?: emptyList()
    if (mengocBcRecords.isNotEmpty()) {
        val lastBc = mengocBcRecords.last().date
        val ageNow = dob?.let { VaccineDateUtils.getAgeAtDate(it, analysisDate) }
        
        val interaction = rule.interactions["VA-MENGOC-BC"]
        if (interaction != null) {
            val minAgeGte = interaction.appliesWhenAgeMonthsGte ?: 24
            if (ageNow != null && ageNow.totalMonths >= minAgeGte) {
                val earliestSafeDate = lastBc.plusDays(interaction.minIntervalDays.toLong())
                if (analysisDate.isBefore(earliestSafeDate)) {
                    warnings.add(interaction.message)
                }
            }
        }
    }
    
    // 2. Six_In_One_Combined interaction
    val sixInOneRecords = allAdministered["Six_In_One_Combined"] ?: emptyList()
    if (sixInOneRecords.isNotEmpty()) {
        val lastSixInOne = sixInOneRecords.last().date
        val interaction = rule.interactions["Six_In_One_Combined"]
        if (interaction != null) {
            val earliestSafeDate = lastSixInOne.plusDays(interaction.minIntervalDays.toLong())
            if (analysisDate.isBefore(earliestSafeDate)) {
                warnings.add(interaction.message)
            }
        }
    }
    
    return if (warnings.isNotEmpty()) warnings.joinToString(" ") else null
}
