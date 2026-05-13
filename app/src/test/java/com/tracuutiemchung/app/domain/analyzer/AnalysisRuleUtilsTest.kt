package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.RecommendationStatus
import com.tracuutiemchung.app.data.rules.VaccineRule
import org.junit.Assert.*
import org.junit.Test
import java.time.LocalDate

class AnalysisRuleUtilsTest {

    private fun createMockRule(
        minAgeMonths: Int? = null,
        minAgeWeeks: Int? = null,
        minAgeDays: Int? = null,
        minAgeYears: Int? = null
    ): VaccineRule {
        return VaccineRule(
            vaccineKey = "test",
            displayName = "Test Vaccine",
            type = "single",
            minAgeMonthsAtFirstDose = minAgeMonths,
            minAgeWeeksAtFirstDose = minAgeWeeks,
            minAgeDaysAtFirstDose = minAgeDays,
            minAgeYearsAtFirstDose = minAgeYears
        )
    }

    @Test
    fun ageStatus_tooYoungByMonths() {
        val dob = LocalDate.of(2024, 1, 1)
        val analysisDate = LocalDate.of(2024, 6, 30) // 5 months 29 days
        val rule = createMockRule(minAgeMonths = 6)
        
        val result = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, "BCG")
        
        assertEquals("too_young", result.statusTags.first())
        assertEquals("cần 6 tháng tuổi", result.message)
        assertEquals(LocalDate.of(2024, 7, 1), result.earliestDate)
    }

    @Test
    fun ageStatus_eligibleByMonths() {
        val dob = LocalDate.of(2024, 1, 1)
        val analysisDate = LocalDate.of(2024, 7, 1)
        val rule = createMockRule(minAgeMonths = 6)
        
        val result = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, "BCG")
        
        assertEquals("eligible", result.statusTags.first())
        assertEquals("đủ điều kiện tuổi", result.message)
    }

    @Test
    fun ageStatus_noDob() {
        val result = AnalysisRuleUtils.getAgeStatusAndEarliestDate(null, LocalDate.now(), createMockRule(), "BCG")
        assertEquals("error_dob", result.statusTags.first())
        assertTrue(result.message.contains("Không có ngày sinh"))
    }

    @Test
    fun firstDoseValidity_tooEarlyByWeeks() {
        val dob = LocalDate.of(2024, 1, 1)
        val firstDoseDate = LocalDate.of(2024, 1, 28) // exactly 27 days, need 4 weeks (28 days)
        val rule = createMockRule(minAgeWeeks = 4)
        
        val (isValid, missingItem) = AnalysisRuleUtils.checkFirstDoseAgeValidity(dob, firstDoseDate, rule, "Test")
        
        assertFalse(isValid)
        assertNotNull(missingItem)
        assertTrue(missingItem?.description?.contains("cần 4 tuần") == true)
        assertEquals(listOf("error_age_first_dose", "too_early"), missingItem?.statusTags)
    }

    @Test
    fun firstDoseValidity_validAge() {
        val dob = LocalDate.of(2024, 1, 1)
        val firstDoseDate = LocalDate.of(2024, 1, 29)
        val rule = createMockRule(minAgeWeeks = 4)
        
        val (isValid, _) = AnalysisRuleUtils.checkFirstDoseAgeValidity(dob, firstDoseDate, rule, "Test")
        
        assertTrue(isValid)
    }

    @Test
    fun missingItemToRecommendation_dueStatus() {
        val item = MissingItem(
            vaccineNameForPopup = "6in1",
            description = "Đã đến lịch",
            earliestNextDoseDate = LocalDate.of(2024, 5, 1),
            statusTags = listOf("due")
        )
        
        val rec = item.toRecommendation(LocalDate.of(2024, 5, 1))
        
        assertEquals(RecommendationStatus.DUE_NOW, rec.status)
        assertEquals("2024-05-01", rec.suggestedDate)
    }

    @Test
    fun missingItemToRecommendation_errorStatus() {
        val item = MissingItem(
            vaccineNameForPopup = "Test",
            description = "Lỗi gì đó",
            earliestNextDoseDate = null,
            statusTags = listOf("error_unknown")
        )
        
        val rec = item.toRecommendation(LocalDate.of(2024, 5, 1))
        
        assertEquals(RecommendationStatus.NEEDS_REVIEW, rec.status)
        assertTrue(rec.warnings.contains("Lỗi gì đó"))
    }
}
