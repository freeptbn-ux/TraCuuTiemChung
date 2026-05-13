package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Port of Pneumococcal special logic from rule_processor.py
 */
fun checkPneumococcalSpecial(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>,
): List<MissingItem> {
    val displayName = rule.displayName
    
    // 1. Identify all pneumococcal types
    val pneuTypes = listOf("PNEU_13", "PNEU_10", "PNEU_15", "PNEU_23")
    val allPneuRecords = allAdministered.filterKeys { it in pneuTypes }
    
    // 2. Check for interleaving (mixing different types in primary series)
    val distinctTypesPresent = allPneuRecords.filter { it.value.isNotEmpty() }.keys
    if (distinctTypesPresent.size > 1 && !distinctTypesPresent.contains("PNEU_23")) {
        // Warning about mixing different conjugate vaccines
        return listOf(
            MissingItem(
                displayName,
                "Cảnh báo: Phát hiện tiêm xen kẽ các loại vắc-xin phế cầu khác nhau (${distinctTypesPresent.joinToString()}). Nên trung thành với một loại.",
                analysisDate,
                listOf("warning", "interleaved_pneu")
            )
        )
    }

    // 3. Special logic for Pneumovax 23 (PNEU_23)
    if (ruleKey == "PNEU_23") {
        if (dob != null) {
            val ageAtAnalysis = VaccineDateUtils.getAgeAtDate(dob, analysisDate)
            if (ageAtAnalysis != null && ageAtAnalysis.totalYears < 2) {
                return listOf(
                    MissingItem(
                        displayName,
                        "Pneumovax 23 chỉ dành cho trẻ từ 2 tuổi trở lên.",
                        VaccineDateUtils.addYears(dob, 2),
                        listOf("too_young")
                    )
                )
            }
        }
        
        if (administeredRecords.isNotEmpty()) return emptyList() // Typically one dose

        // Suggest PNEU_23 if older than 2 years and completed primary conjugate series
        val primaryCompleted = pneuTypes.filter { it != "PNEU_23" }.any { type ->
            val records = allAdministered[type].orEmpty()
            val required = allRules[type]?.requiredDoses ?: 4
            records.size >= required
        }

        if (primaryCompleted) {
            return listOf(
                MissingItem(
                    displayName,
                    "Đã hoàn thành phác đồ phế cầu cộng hợp. Khuyến cáo tiêm thêm 1 mũi Pneumovax 23 để mở rộng phạm vi bảo vệ.",
                    analysisDate,
                    listOf("due", "booster_pneu23")
                )
            )
        }
        
        return emptyList() // Don't suggest if primary not done
    }

    // 4. Delegate to SingleSeries or AgeDependent based on the specific rule
    // Note: Synflorix (PNEU_10) and Prevenar 13 (PNEU_13) are usually age_dependent_series
    return when (rule.type) {
        "age_dependent_series" -> checkAgeDependentSeries(ruleKey, rule, administeredRecords, dob, analysisDate, allRules, allAdministered)
        else -> checkSingleVaccineSeries(ruleKey, rule, administeredRecords, dob, analysisDate, allRules, allAdministered)
    }
}
