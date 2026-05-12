package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.rules.DoseSpecificRule
import com.tracuutiemchung.app.data.rules.VaccineRule
import org.junit.Assert.*
import org.junit.Test
import java.time.LocalDate

class SingleSeriesCheckerTest {

    private val dob = LocalDate.of(2023, 1, 1)
    private val analysisDate = LocalDate.of(2024, 1, 1)

    @Test
    fun completedSeries_noBooster_returnsEmpty() {
        val rule = createRule(requiredDoses = 2)
        val records = listOf(
            createRecord("V1", "2023-02-01"),
            createRecord("V1", "2023-03-01")
        )
        val result = checkSingleVaccineSeries("V1", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertTrue(result.isEmpty())
    }

    @Test
    fun completedSeries_withBooster_due_returnsDueItem() {
        val rule = createRule(requiredDoses = 2, boosterIntervalYears = 1)
        val records = listOf(
            createRecord("V1", "2023-02-01"),
            createRecord("V1", "2023-03-01")
        )
        // Booster due on 2024-03-01.
        val result = checkSingleVaccineSeries("V1", rule, records, dob, LocalDate.of(2024, 3, 1), emptyMap(), emptyMap())
        assertEquals(1, result.size)
        assertTrue(result[0].statusTags.contains("booster_due"))
        assertEquals(LocalDate.of(2024, 3, 1), result[0].earliestNextDoseDate)
    }

    @Test
    fun completedSeries_withBooster_upcoming_returnsInfoItem() {
        val rule = createRule(requiredDoses = 2, boosterIntervalYears = 1)
        val records = listOf(
            createRecord("V1", "2023-02-01"),
            createRecord("V1", "2023-03-01")
        )
        // Booster due on 2024-03-01. Analysis on 2024-01-01 is upcoming.
        val result = checkSingleVaccineSeries("V1", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, result.size)
        assertTrue(result[0].statusTags.contains("booster_upcoming"))
    }

    @Test
    fun completedSeries_withBooster_maxAgeReached_returnsEmpty() {
        val rule = createRule(requiredDoses = 1, boosterIntervalYears = 5, boosterMaxAgeYears = 4)
        val records = listOf(createRecord("V1", "2023-02-01"))
        // Booster due in 5 years (2028). But max age is 4 years.
        val result = checkSingleVaccineSeries("V1", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertTrue(result.isEmpty())
    }

    @Test
    fun noDoses_tooYoung_returnsWithEarliestDate() {
        val rule = createRule(minAgeMonthsOverall = 12)
        // dob: 2023-01-01. 12 months = 2024-01-01.
        // Analysis on 2023-12-01 is too young.
        val result = checkSingleVaccineSeries("V1", rule, emptyList(), dob, LocalDate.of(2023, 12, 1), emptyMap(), emptyMap())
        assertEquals(1, result.size)
        assertTrue(result[0].statusTags.contains("too_young"))
        assertEquals(LocalDate.of(2024, 1, 1), result[0].earliestNextDoseDate)
    }

    @Test
    fun oneDoseOf3_usesCorrectIntervalForDose2() {
        // min_interval_days=[null, 30, 30, 360]
        val rule = createRule(requiredDoses = 3, minIntervalDays = listOf(null, 30, 30, 360))
        val records = listOf(createRecord("V1", "2023-02-01"))
        // Next dose (2) needs 30 days interval -> 2023-03-03
        val result = checkSingleVaccineSeries("V1", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, result.size)
        assertEquals(LocalDate.of(2023, 2, 1).plusDays(30), result[0].earliestNextDoseDate)
    }

    @Test
    fun doseSpecific_alternativeAgeRange_usesEarlierDate() {
        val rule = createRule(
            requiredDoses = 2,
            minIntervalDays = listOf(null, 180),
            doseSpecificRules = mapOf("2" to DoseSpecificRule(alternativeMinAgeYears = 1))
        )
        // dob: 2023-01-01
        // dose 1: 2023-10-01
        // interval: 2023-10-01 + 180 = 2024-03-29
        // alt_age (1yr): 2024-01-01
        // min(interval, alt_age) = 2024-01-01
        val records = listOf(createRecord("V1", "2023-10-01"))
        val result = checkSingleVaccineSeries("V1", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(LocalDate.of(2024, 1, 1), result[0].earliestNextDoseDate)
    }

    @Test
    fun doseSpecific_minAbsoluteAgeMonths_constrainsDate() {
        val rule = createRule(
            requiredDoses = 2,
            minIntervalDays = listOf(null, 30),
            doseSpecificRules = mapOf("2" to DoseSpecificRule(minAbsoluteAgeMonths = 6))
        )
        val records = listOf(createRecord("V1", "2023-02-01"))
        // interval: 2023-03-03
        // abs_min_age: 2023-01-01 + 6 months = 2023-07-01
        // max(interval, abs_min_age) = 2023-07-01
        val result = checkSingleVaccineSeries("V1", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(LocalDate.of(2023, 7, 1), result[0].earliestNextDoseDate)
    }

    private fun createRule(
        requiredDoses: Int? = 1,
        minAgeMonthsOverall: Int? = null,
        minIntervalDays: List<Int?> = emptyList(),
        boosterIntervalYears: Int? = null,
        boosterMaxAgeYears: Int? = null,
        doseSpecificRules: Map<String, DoseSpecificRule> = emptyMap()
    ) = VaccineRule(
        vaccineKey = "TEST",
        displayName = "Test Vaccine",
        type = "single_series",
        requiredDoses = requiredDoses,
        minAgeMonthsOverall = minAgeMonthsOverall,
        minIntervalDays = minIntervalDays,
        boosterIntervalYears = boosterIntervalYears,
        boosterMaxAgeYears = boosterMaxAgeYears,
        doseSpecificRules = doseSpecificRules
    )

    private fun createRecord(name: String, date: String) = ParsedRecord(
        source = VaccinationRecord(vaccineName = name, vaccinationDate = date),
        date = LocalDate.parse(date)
    )
}
