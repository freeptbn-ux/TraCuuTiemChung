package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.rules.AgeBasedRegimen
import com.tracuutiemchung.app.data.rules.VaccineRule
import org.junit.Assert.*
import org.junit.Test
import java.time.LocalDate

class MmrEquivalentGroupCheckerTest {

    private val dob = LocalDate.of(2023, 1, 1)
    private val analysisDate = LocalDate.of(2025, 1, 1)

    private val mmrRegimens = listOf(
        AgeBasedRegimen(regimenName = "9-11m", minAgeAtFirstDoseMonths = 9, maxAgeAtFirstDoseMonths = 11, dosesRequired = 3, minIntervalDays = listOf(null, 90, 1095)),
        AgeBasedRegimen(regimenName = "12m-7y", minAgeAtFirstDoseMonths = 12, maxAgeAtFirstDoseMonths = 83, dosesRequired = 2, minIntervalDays = listOf(null, 90)),
        AgeBasedRegimen(regimenName = ">=7y", minAgeAtFirstDoseMonths = 84, dosesRequired = 2, minIntervalDays = listOf(null, 28))
    )

    private fun createMmrRule() = VaccineRule(
        vaccineKey = "MMR_Group",
        displayName = "MMR",
        type = "mmr_equivalent_group",
        regimens = mmrRegimens
    )

    @Test
    fun mvvacInteraction_0MmrDoses_needsMmrAfter84Days() {
        val rule = createMmrRule()
        val mvvacRecord = createRecord("MVVAC", "2023-10-01") // 9 months old
        val allAdministered = mapOf("MVVAC" to listOf(mvvacRecord))
        
        val result = checkMmrEquivalentGroup("MMR_Group", rule, emptyList(), dob, analysisDate, emptyMap(), allAdministered)
        
        assertEquals(1, result.size)
        assertTrue(result[0].description.contains("Đã tiêm Sởi (MVVAC)"))
        // Earliest date = max(2023-10-01 + 84d, dob + 12m)
        // 2023-10-01 + 84d = 2023-12-24
        // dob + 12m = 2024-01-01
        // Max = 2024-01-01
        assertEquals(LocalDate.of(2024, 1, 1), result[0].earliestNextDoseDate)
    }

    @Test
    fun mvvacInteraction_1MmrDose_needs2ndMmrAfter3Years() {
        val rule = createMmrRule()
        val mvvacRecord = createRecord("MVVAC", "2023-10-01")
        val mmrRecord = createRecord("MMR", "2024-01-01") // MMR 1 after MVVAC
        val allAdministered = mapOf(
            "MVVAC" to listOf(mvvacRecord),
            "MMR_Group" to listOf(mmrRecord)
        )
        
        val result = checkMmrEquivalentGroup("MMR_Group", rule, listOf(mmrRecord), dob, analysisDate, emptyMap(), allAdministered)
        
        println("DEBUG: Result size = ${result.size}")
        if (result.isNotEmpty()) {
            println("DEBUG: Description = ${result[0].description}")
            println("DEBUG: StatusTags = ${result[0].statusTags}")
        }
        
        assertTrue("Description should contain expected text", result[0].description.contains("MVVAC) và 1 mũi MMR"))
        // Earliest date = 2024-01-01 + 1095 days = 2026-12-31
        assertEquals(LocalDate.of(2026, 12, 31), result[0].earliestNextDoseDate)
    }

    @Test
    fun mvvacInteraction_2MmrDoses_completed() {
        val rule = createMmrRule()
        val mvvacRecord = createRecord("MVVAC", "2023-10-01")
        val mmr1 = createRecord("MMR", "2024-01-01")
        val mmr2 = createRecord("MMR", "2024-04-01")
        val records = listOf(mmr1, mmr2)
        val allAdministered = mapOf(
            "MVVAC" to listOf(mvvacRecord),
            "MMR_Group" to records
        )
        
        val result = checkMmrEquivalentGroup("MMR_Group", rule, records, dob, analysisDate, emptyMap(), allAdministered)
        
        assertTrue(result.isEmpty())
    }

    @Test
    fun noMvvac_standardMmr_regimen12mTo7y() {
        val rule = createMmrRule()
        // First dose at 12 months
        val mmr1 = createRecord("MMR", "2024-01-01")
        val records = listOf(mmr1)
        
        val result = checkMmrEquivalentGroup("MMR_Group", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].vaccineNameForPopup.contains("12m-7y"))
        assertEquals(LocalDate.of(2024, 1, 1).plusDays(90), result[0].earliestNextDoseDate)
    }

    @Test
    fun mvvacCoverageInSingleSeriesChecker() {
        val mvvacRule = VaccineRule(
            vaccineKey = "MVVAC",
            displayName = "MVVAC",
            type = "single_dose_min_age",
            requiredDoses = 1,
            providesMeaslesProtection = true
        )
        val mmrRule = VaccineRule(
            vaccineKey = "MMR_Group",
            displayName = "MMR",
            type = "mmr_equivalent_group",
            providesMeaslesProtection = true
        )
        
        val allRules = mapOf("MVVAC" to mvvacRule, "MMR_Group" to mmrRule)
        
        // Patient has MMR but no MVVAC
        val mmrRecord = createRecord("MMR", "2024-01-01")
        val allAdministered = mapOf("MMR_Group" to listOf(mmrRecord))
        
        // Check MVVAC rule
        val result = checkSingleVaccineSeries("MVVAC", mvvacRule, emptyList(), dob, analysisDate, allRules, allAdministered)
        
        // Should be empty because covered by MMR
        assertTrue(result.isEmpty())
    }

    private fun createRecord(name: String, date: String) = ParsedRecord(
        source = VaccinationRecord(vaccineName = name, vaccinationDate = date),
        date = LocalDate.parse(date)
    )
}
