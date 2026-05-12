package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.PatientInfo
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.model.VaccineRecommendation
import com.tracuutiemchung.app.data.model.RecommendationStatus
import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.VaccineNameNormalizer
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Data class representing the input for vaccination analysis.
 */
data class AnalysisInput(
    val patient: PatientInfo,
    val vaccinationRecords: List<VaccinationRecord>,
    val rules: List<VaccineRule>,
    val analysisDate: LocalDate,
)

/**
 * Core engine for analyzing vaccination records and providing recommendations.
 * This class achieves functional parity with the legacy Python analysis system.
 */
class VaccineAnalysisEngine {

    /**
     * Analyzes the patient's vaccination records against a set of rules and provides recommendations.
     *
     * @param input The input data for analysis including patient info, records, and rules.
     * @return A list of [VaccineRecommendation] for each applicable rule.
     */
    fun analyze(input: AnalysisInput): List<VaccineRecommendation> {
        if (input.rules.isEmpty()) {
            return listOf(
                VaccineRecommendation(
                    vaccineName = "Quy tắc vắc-xin",
                    status = RecommendationStatus.NOT_ENOUGH_DATA,
                    reason = "Chưa có dữ liệu quy tắc để phân tích.",
                    warnings = listOf("Thiếu danh sách quy tắc vắc-xin."),
                )
            )
        }

        val dob = VaccineDateUtils.parseFlexibleDate(input.patient.dateOfBirth)
        val globalWarnings = buildList {
            if (input.patient.dateOfBirth.isNullOrBlank()) {
                add("Thiếu ngày sinh, không thể kiểm tra điều kiện tuổi.")
            } else if (dob == null) {
                add("Không parse được ngày sinh: ${input.patient.dateOfBirth}.")
            }
        }
        val normalizer = VaccineNameNormalizer(input.rules)
        val recordsByRule = groupValidRecords(input.vaccinationRecords, normalizer)
        
        val recordWarnings = collectRecordWarnings(input.vaccinationRecords)

        val allRulesMap = input.rules.associateBy { it.vaccineKey }

        // Phase 09.2: Pneumococcal pre-processing
        val skippedRuleKeys = preProcessPneumococcal(input.rules, recordsByRule)

        val recommendations = input.rules
            .filter { it.vaccineKey !in skippedRuleKeys }
            .map { rule ->
                analyzeRule(
                    rule = rule,
                    administered = recordsByRule[rule.vaccineKey].orEmpty(),
                    dob = dob,
                    analysisDate = input.analysisDate,
                    warnings = globalWarnings + recordWarnings,
                    allRules = allRulesMap,
                    allAdministered = recordsByRule
                )
            }

        return recommendations
    }

    private fun analyzeRule(
        rule: VaccineRule,
        administered: List<ParsedRecord>,
        dob: LocalDate?,
        analysisDate: LocalDate,
        warnings: List<String>,
        allRules: Map<String, VaccineRule>,
        allAdministered: Map<String, List<ParsedRecord>>
    ): VaccineRecommendation {
        // Phase 09.4: VA-MENGOC-BC reverse interaction
        val enhancedWarnings = if (rule.vaccineKey == "VA-MENGOC-BC" && dob != null) {
            val ageMonths = VaccineDateUtils.getAgeAtDate(dob, analysisDate)?.totalMonths ?: 0
            val hasAcyw = allAdministered.keys.any { it.contains("ACYW", ignoreCase = true) } ||
                    allAdministered.values.flatten().any { it.source.vaccineName.contains("ACYW", ignoreCase = true) }
            if (ageMonths >= 24 && hasAcyw) {
                warnings + "Trẻ từ 24 tháng tuổi và đã tiêm ACYW thì không cần tiêm VA-MENGOC-BC (trừ trường hợp đặc biệt)."
            } else warnings
        } else warnings

        // Delegate to specific checkers based on rule type
        val missingItems = when (rule.type) {
            "single_series", "single_dose_min_age", "single_series_min_age", "series" -> {
                checkSingleVaccineSeries(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            "age_dependent_series" -> {
                checkAgeDependentSeries(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            "mmr_equivalent_group" -> {
                checkMmrEquivalentGroup(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            "flu_group" -> {
                checkFluGroup(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            "group_cumulative_unique_doses", "group_cumulative_unique_doses_min_age" -> {
                checkCumulativeGroupDoses(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            "meningococcal_acyw_group" -> {
                checkMeningococcalAcywGroup(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            "group_alternative_courses", "group_alternative_courses_min_age", "group_alternative_courses_age_range" -> {
                checkAlternativeCoursesGroup(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            "pneumococcal_special" -> {
                checkPneumococcalSpecial(
                    ruleKey = rule.vaccineKey,
                    rule = rule,
                    administeredRecords = administered,
                    dob = dob,
                    analysisDate = analysisDate,
                    allRules = allRules,
                    allAdministered = allAdministered
                )
            }
            else -> {
                listOf(
                    MissingItem(
                        vaccineNameForPopup = rule.displayName,
                        description = "Cảnh báo: Loại quy tắc '${rule.type}' chưa được hỗ trợ.",
                        earliestNextDoseDate = null,
                        statusTags = listOf("unsupported_rule_type", "warning")
                    )
                )
            }
        }

        // Map MissingItem to Recommendation
        if (missingItems.isEmpty()) {
            val count = administered.size
            val required = rule.requiredDoses ?: count
            val reason = if (count > 0) "Đã hoàn thành phác đồ ($count/$required mũi)." else "Đã hoàn thành phác đồ."
            return VaccineRecommendation(
                vaccineName = rule.displayName,
                status = RecommendationStatus.COMPLETED,
                reason = reason,
                warnings = enhancedWarnings
            )
        }

        val recommendation = missingItems.first().toRecommendation(analysisDate)
        return recommendation.copy(warnings = enhancedWarnings + recommendation.warnings)
    }

    private fun preProcessPneumococcal(
        rules: List<VaccineRule>,
        recordsByRule: Map<String, List<ParsedRecord>>
    ): Set<String> {
        val skipped = mutableSetOf<String>()
        val pneuRuleKeys = listOf("PNEU_10", "PNEU_13", "PNEU_15")
        val pneuRulesInInput = rules.filter { it.vaccineKey in pneuRuleKeys }
        if (pneuRulesInInput.size <= 1) return emptySet()

        // logic: favor the one already started
        val rulesWithDoses = pneuRulesInInput.filter { recordsByRule[it.vaccineKey]?.isNotEmpty() == true }
        if (rulesWithDoses.isNotEmpty()) {
            val bestRule = rulesWithDoses.maxByOrNull { recordsByRule[it.vaccineKey]?.size ?: 0 }!!
            pneuRulesInInput.forEach {
                if (it.vaccineKey != bestRule.vaccineKey) skipped.add(it.vaccineKey)
            }
        } else {
            val preferred = pneuRulesInInput.find { it.vaccineKey == "PNEU_13" } ?: pneuRulesInInput.first()
            pneuRulesInInput.forEach {
                if (it.vaccineKey != preferred.vaccineKey) skipped.add(it.vaccineKey)
            }
        }
        return skipped
    }

    private fun groupValidRecords(
        records: List<VaccinationRecord>,
        normalizer: VaccineNameNormalizer,
    ): Map<String, List<ParsedRecord>> = records.mapNotNull { record ->
        val date = VaccineDateUtils.parseFlexibleDate(record.vaccinationDate) ?: return@mapNotNull null
        val rule = normalizer.normalize(record.vaccineName) ?: return@mapNotNull null
        rule.vaccineKey to ParsedRecord(record, date)
    }.groupBy({ it.first }, { it.second })
        .mapValues { (_, value) -> value.sortedBy { it.date } }

    private fun collectRecordWarnings(records: List<VaccinationRecord>): List<String> = records.flatMap { record ->
        buildList {
            if (VaccineDateUtils.parseFlexibleDate(record.vaccinationDate) == null) {
                add("Không parse được ngày tiêm của ${record.vaccineName}: ${record.vaccinationDate}.")
            }
        }
    }.distinct()
}
