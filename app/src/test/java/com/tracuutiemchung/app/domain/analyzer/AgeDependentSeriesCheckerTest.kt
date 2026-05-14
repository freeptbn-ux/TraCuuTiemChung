package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.rules.AgeBasedRegimen
import com.tracuutiemchung.app.data.rules.DoseSpecificRule
import com.tracuutiemchung.app.data.rules.VaccineRule
import org.junit.Assert.*
import org.junit.Test
import java.time.LocalDate

class AgeDependentSeriesCheckerTest {

    private val dob = LocalDate.of(2023, 1, 1)
    private val analysisDate = LocalDate.of(2025, 1, 1)

    private val prevenar13Rules = listOf(
        AgeBasedRegimen(regimenName = "Phác đồ 4 mũi", maxAgeAtFirstDoseMonths = 6, dosesRequired = 4, minIntervalDays = listOf(null, 30, 30, 240)),
        AgeBasedRegimen(regimenName = "Phác đồ 3 mũi", minAgeAtFirstDoseMonths = 7, maxAgeAtFirstDoseMonths = 11, dosesRequired = 3, minIntervalDays = listOf(null, 30, 180)),
        AgeBasedRegimen(regimenName = "Phác đồ 2 mũi", minAgeAtFirstDoseMonths = 12, maxAgeAtFirstDoseMonths = 23, dosesRequired = 2, minIntervalDays = listOf(null, 60)),
        AgeBasedRegimen(regimenName = "Phác đồ 1 mũi", minAgeAtFirstDoseMonths = 24, dosesRequired = 1, minIntervalDays = listOf(null))
    )

    private fun createPrevenar13Rule() = VaccineRule(
        vaccineKey = "prevenar13",
        displayName = "Prevenar 13",
        type = "age_dependent_series",
        minAgeMonthsOverall = 2,
        rulesByAge = prevenar13Rules
    )

    @Test
    fun prevenar13_firstDoseAt4Months_needs4Doses() {
        val rule = createPrevenar13Rule()
        // First dose at 4 months (2023-05-01)
        val records = listOf(createRecord("Prevenar 13", "2023-05-01"))
        val result = checkAgeDependentSeries("prevenar13", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].description.contains("Đã tiêm 1/4 mũi"))
        assertTrue(result[0].vaccineNameForPopup.contains("Phác đồ 4 mũi"))
    }

    @Test
    fun prevenar13_firstDoseAt9Months_needs3Doses() {
        val rule = createPrevenar13Rule()
        // First dose at 9 months (2023-10-01)
        val records = listOf(createRecord("Prevenar 13", "2023-10-01"))
        val result = checkAgeDependentSeries("prevenar13", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].description.contains("Đã tiêm 1/3 mũi"))
        assertTrue(result[0].vaccineNameForPopup.contains("Phác đồ 3 mũi"))
    }

    @Test
    fun prevenar13_firstDoseAt14Months_needs2Doses() {
        val rule = createPrevenar13Rule()
        // First dose at 14 months (2024-03-01)
        val records = listOf(createRecord("Prevenar 13", "2024-03-01"))
        val result = checkAgeDependentSeries("prevenar13", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].description.contains("Đã tiêm 1/2 mũi"))
        assertTrue(result[0].vaccineNameForPopup.contains("Phác đồ 2 mũi"))
    }

    @Test
    fun prevenar13_firstDoseAt24Months_needs1Dose() {
        val rule = createPrevenar13Rule()
        // First dose at 24 months (2025-01-01)
        val records = listOf(createRecord("Prevenar 13", "2025-01-01"))
        val result = checkAgeDependentSeries("prevenar13", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        // Completed because required doses = 1
        assertTrue(result.isEmpty())
    }

    @Test
    fun prevenar13_3of4DosesCompleted_nextDoseInterval240Days() {
        val rule = createPrevenar13Rule()
        // First dose at 4 months (2023-05-01) -> Phác đồ 4 mũi
        // Dose 2: 2023-06-01
        // Dose 3: 2023-07-01
        // Dose 4: 2023-07-01 + 240 days = 2024-02-26
        val records = listOf(
            createRecord("Prevenar 13", "2023-05-01"),
            createRecord("Prevenar 13", "2023-06-01"),
            createRecord("Prevenar 13", "2023-07-01")
        )
        val result = checkAgeDependentSeries("prevenar13", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertEquals(LocalDate.of(2023, 7, 1).plusDays(240), result[0].earliestNextDoseDate)
    }

    @Test
    fun synflorix_7to11Months_dose2RequiresMinAbsAge12Months() {
        val synflorixRules = listOf(
            AgeBasedRegimen(
                regimenName = "7-11 tháng",
                minAgeAtFirstDoseMonths = 7,
                maxAgeAtFirstDoseMonths = 11,
                dosesRequired = 3,
                minIntervalDays = listOf(null, 30, 180),
                doseSpecificRules = mapOf("2" to DoseSpecificRule(minAbsoluteAgeMonths = 12))
            )
        )
        val rule = VaccineRule(
            vaccineKey = "synflorix",
            displayName = "Synflorix",
            type = "age_dependent_series",
            rulesByAge = synflorixRules
        )
        
        // First dose at 8 months (2023-09-01)
        val records = listOf(createRecord("Synflorix", "2023-09-01"))
        
        // Next dose (2) needs:
        // - Interval: 2023-09-01 + 30 days = 2023-10-01
        // - Abs age: 2023-01-01 + 12 months = 2024-01-01
        // -> Max(2023-10-01, 2024-01-01) = 2024-01-01
        
        val result = checkAgeDependentSeries("synflorix", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(LocalDate.of(2024, 1, 1), result[0].earliestNextDoseDate)
    }

    @Test
    fun noDoses_returnsGuidanceBasedOnCurrentAge() {
        val rule = createPrevenar13Rule()
        // Analysis at 3 months
        val result = checkAgeDependentSeries("prevenar13", rule, emptyList(), dob, LocalDate.of(2023, 4, 1), emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].description.contains("đủ điều kiện tuổi"))
        assertTrue(result[0].description.contains("Phác đồ dự kiến 4 mũi"))
    }

    @Test
    fun noDob_returnsErrorDob() {
        val rule = createPrevenar13Rule()
        val result = checkAgeDependentSeries("prevenar13", rule, emptyList(), null, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].statusTags.contains("error_dob"))
    }

    @Test
    fun firstDose_tooEarlyOverall_returnsError() {
        val rule = createPrevenar13Rule() // minAgeMonthsOverall = 2
        // First dose at 1 month (2023-02-01)
        val records = listOf(createRecord("Prevenar 13", "2023-02-01"))
        val result = checkAgeDependentSeries("prevenar13", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].statusTags.contains("error_age_first_dose"))
        assertTrue(result[0].description.contains("quá sớm"))
    }

    @Test
    fun firstDose_outsideAllRanges_returnsNoMatchingRule() {
        // Create rule with limited ranges
        val rule = VaccineRule(
            vaccineKey = "test",
            displayName = "Test",
            type = "age_dependent_series",
            rulesByAge = listOf(
                AgeBasedRegimen(minAgeAtFirstDoseMonths = 0, maxAgeAtFirstDoseMonths = 5, dosesRequired = 1)
            )
        )
        // First dose at 6 months
        val records = listOf(createRecord("Test", "2023-07-01"))
        val result = checkAgeDependentSeries("test", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertEquals(1, result.size)
        assertTrue(result[0].statusTags.contains("error_no_matching_rule"))
    }

    @Test
    fun allDosesCompleted_returnsEmpty() {
        val rule = createPrevenar13Rule()
        // First dose at 24 months (2025-01-01) -> 1 dose required
        val records = listOf(createRecord("Prevenar 13", "2025-01-01"))
        val result = checkAgeDependentSeries("prevenar13", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        
        assertTrue(result.isEmpty())
    }

    private fun createRecord(name: String, date: String) = ParsedRecord(
        source = VaccinationRecord(vaccineName = name, vaccinationDate = date),
        date = LocalDate.parse(date)
    )
}
