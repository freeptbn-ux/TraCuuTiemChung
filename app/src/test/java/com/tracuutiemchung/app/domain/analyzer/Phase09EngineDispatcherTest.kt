package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.PatientInfo
import com.tracuutiemchung.app.data.model.RecommendationStatus
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.rules.VaccineRule
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import java.time.LocalDate

class Phase09EngineDispatcherTest {
    private val engine = VaccineAnalysisEngine()

    @Test
    fun testPneumococcalPreProcessing_SkipsUnstartedRules() {
        val rules = listOf(
            VaccineRule(vaccineKey = "PNEU_10", displayName = "Synflorix", type = "age_dependent_series"),
            VaccineRule(vaccineKey = "PNEU_13", displayName = "Prevenar 13", type = "age_dependent_series")
        )
        val input = AnalysisInput(
            patient = PatientInfo("Bé A", "090", dateOfBirth = "2024-01-01"),
            vaccinationRecords = listOf(
                VaccinationRecord("Synflorix", vaccinationDate = "2024-03-01")
            ),
            rules = rules,
            analysisDate = LocalDate.parse("2024-05-01")
        )
        val result = engine.analyze(input)
        
        // Should only have Synflorix, Prevenar 13 should be skipped
        assertEquals(1, result.size)
        assertEquals("Synflorix", result[0].vaccineName)
    }

    @Test
    fun testPneumococcalPreProcessing_FavorsPneu13IfNoneStarted() {
        val rules = listOf(
            VaccineRule(vaccineKey = "PNEU_10", displayName = "Synflorix", type = "age_dependent_series"),
            VaccineRule(vaccineKey = "PNEU_13", displayName = "Prevenar 13", type = "age_dependent_series")
        )
        val input = AnalysisInput(
            patient = PatientInfo("Bé A", "090", dateOfBirth = "2024-01-01"),
            vaccinationRecords = emptyList(),
            rules = rules,
            analysisDate = LocalDate.parse("2024-05-01")
        )
        val result = engine.analyze(input)
        
        // Should favor PNEU_13
        assertEquals(1, result.size)
        assertEquals("Prevenar 13", result[0].vaccineName)
    }

    @Test
    fun testUnknownRuleType_ReturnsWarning() {
        val rules = listOf(
            VaccineRule(vaccineKey = "UNKNOWN", displayName = "Unknown Vaccine", type = "unsupported_type")
        )
        val input = AnalysisInput(
            patient = PatientInfo("Bé A", "090", dateOfBirth = "2024-01-01"),
            vaccinationRecords = emptyList(),
            rules = rules,
            analysisDate = LocalDate.parse("2024-05-01")
        )
        val result = engine.analyze(input)
        
        assertEquals(1, result.size)
        assertEquals(RecommendationStatus.NEEDS_REVIEW, result[0].status)
        assertTrue(result[0].reason.contains("chưa được hỗ trợ"))
    }

    @Test
    fun testVaMengocBcReverseInteraction() {
        val rules = listOf(
            VaccineRule(vaccineKey = "VA-MENGOC-BC", displayName = "VA-MENGOC-BC", type = "single_series", requiredDoses = 2),
            VaccineRule(vaccineKey = "MEN_ACYW", displayName = "Men ACYW", type = "single_series")
        )
        val input = AnalysisInput(
            patient = PatientInfo("Bé A", "090", dateOfBirth = "2020-01-01"),
            vaccinationRecords = listOf(
                VaccinationRecord("Men ACYW", vaccinationDate = "2021-01-01")
            ),
            rules = rules,
            analysisDate = LocalDate.parse("2024-01-01")
        )
        val result = engine.analyze(input)
        
        val mengocBc = result.find { it.vaccineName == "VA-MENGOC-BC" }
        assertTrue(mengocBc?.warnings?.any { it.contains("không cần tiêm VA-MENGOC-BC") } == true)
    }
}
