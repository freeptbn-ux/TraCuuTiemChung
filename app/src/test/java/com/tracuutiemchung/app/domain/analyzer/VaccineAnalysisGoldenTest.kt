package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.RecommendationStatus
import com.tracuutiemchung.app.data.portal.DefaultPortalParser
import com.tracuutiemchung.app.data.portal.RawSourceType
import com.tracuutiemchung.app.data.rules.VaccineRuleRepository
import org.junit.Assert.assertEquals
import org.junit.Assert.assertNotNull
import org.junit.Assert.fail
import org.junit.Test
import java.io.File
import java.time.LocalDate

class VaccineAnalysisGoldenTest {
    private val engine = VaccineAnalysisEngine()
    private val parser = DefaultPortalParser()

    private fun loadRealRules(): List<com.tracuutiemchung.app.data.rules.VaccineRule> {
        var assetsDir = File("app/src/main/assets")
        if (!assetsDir.exists()) {
            assetsDir = File("src/main/assets")
        }
        if (!assetsDir.exists()) {
            throw IllegalStateException("Assets directory not found at any known location")
        }
        val repository = VaccineRuleRepository(
            assetReader = { fileName -> File(assetsDir, fileName).readText() }
        )
        return repository.loadRules()
    }

    @Test
    fun `golden parity - MinhKhoi real data`() {
        // 1. Load data from Minhkhoi.html
        var htmlFile = File("test/Minhkhoi.html")
        if (!htmlFile.exists()) {
            htmlFile = File("../test/Minhkhoi.html")
        }
        if (!htmlFile.exists()) {
            println("Skipping test: Minhkhoi.html not found at ${htmlFile.absolutePath}")
            return
        }
        val html = htmlFile.readText()
        
        val parseResult = parser.parse(html, RawSourceType.HTML).getOrThrow()
        val patientInfo = parseResult.patient ?: throw IllegalStateException("Patient info not found")
        val records = parseResult.vaccinations

        // 2. Load real rules
        val rules = loadRealRules()

        // 3. Run analysis with fixed system date: 12/05/2026
        val analysisDate = LocalDate.of(2026, 5, 12)
        val result = engine.analyze(
            AnalysisInput(
                patient = patientInfo,
                vaccinationRecords = records,
                rules = rules,
                analysisDate = analysisDate
            )
        )

        // 4. Assert parity based on Phase 10 expectations
        
        // BCG: 1 mũi -> COMPLETED
        assertRuleStatus(result, "Lao", RecommendationStatus.COMPLETED)
        
        // 6in1 (Infanrix Hexa): 4 mũi -> COMPLETED
        assertRuleStatus(result, "6 trong 1", RecommendationStatus.COMPLETED)
        
        // Prevenar 13: 4 mũi -> COMPLETED
        assertRuleStatus(result, "Prevenar", RecommendationStatus.COMPLETED)
        
        // Rota: 3 mũi -> COMPLETED
        assertRuleStatus(result, "Rota", RecommendationStatus.COMPLETED)
        
        // MMR Group: Completed (2 MMR-II + 1 MVVAC)
        assertRuleStatus(result, "MMR", RecommendationStatus.COMPLETED)
        
        // Varivax: 2 mũi -> COMPLETED
        assertRuleStatus(result, "Varivax", RecommendationStatus.COMPLETED)
        
        // VA-MENGOC-BC: 2 mũi -> COMPLETED
        assertRuleStatus(result, "MENGOC", RecommendationStatus.COMPLETED)

        // Flu: Should be DUE_LATER or DUE_NOW depending on the season, 
        // but since last dose was 16/11/2025, and analysis is 12/05/2026, 
        // it might be COMPLETED for the current season or DUE_LATER for next.
        // Rule says interval 365 days. 16/11/2025 + 365 = 16/11/2026.
        // So at 12/05/2026, it should be DUE_LATER (scheduled for 2026-11-16).
        assertRuleStatus(result, "Cúm", RecommendationStatus.DUE_LATER)
        val fluRec = result.find { it.vaccineName.contains("Cúm", ignoreCase = true) }
        assertEquals("2026-11-16", fluRec?.suggestedDate)

        // HepA: Never vaccinated -> DUE_NOW or DUE_LATER (too old? no, min age 12-24 months)
        // Minh Khoi born 2020-07-27, at 2026-05-12 he is ~5.8 years old.
        // HepA rules: min age 12m or 24m. He is eligible.
        assertRuleStatus(result, "Viêm gan A", RecommendationStatus.DUE_NOW)

        // Meningococcal ACYW: 1 dose at 34 months (21/05/2023)
        // Menactra rule for >=24m: 1 dose required.
        // But check if it needs booster.
        // In rule: MENACTRA >= 24m -> 1 dose. No booster defined in JSON for >=24m.
        // So it should be COMPLETED.
        assertRuleStatus(result, "ACYW", RecommendationStatus.COMPLETED)
    }

    @Test
    fun `golden test - 4 months old no vaccinations`() {
        val dob = LocalDate.of(2026, 1, 12) // Born 4 months ago relative to 2026-05-12
        val analysisDate = LocalDate.of(2026, 5, 12)
        val patientInfo = com.tracuutiemchung.app.data.model.PatientInfo(fullName = "Bé Bốn Tháng", phoneNumber = "0123456789", dateOfBirth = dob.toString())
        val records = emptyList<com.tracuutiemchung.app.data.model.VaccinationRecord>()
        val rules = loadRealRules()

        val result = engine.analyze(
            AnalysisInput(patient = patientInfo, vaccinationRecords = records, rules = rules, analysisDate = analysisDate)
        )

        // 4 months old should be DUE_NOW for 6-in-1 (mũi 1-2-3), Prevenar (mũi 1-2), Rota (mũi 1-2)
        assertRuleStatus(result, "6 trong 1", RecommendationStatus.DUE_NOW)
        assertRuleStatus(result, "Prevenar", RecommendationStatus.DUE_NOW)
        assertRuleStatus(result, "Rota", RecommendationStatus.DUE_NOW)
        
        // Lao (BCG) should be DUE_NOW (usually given at birth)
        assertRuleStatus(result, "Lao", RecommendationStatus.DUE_NOW)

        // MMR should be DUE_LATER (min age 9m or 12m)
        assertRuleStatus(result, "MMR", RecommendationStatus.DUE_LATER)
    }

    @Test
    fun `golden test - MMR interaction MVVAC then MMR`() {
        val dob = LocalDate.of(2024, 1, 1)
        val analysisDate = LocalDate.of(2026, 5, 12)
        val patientInfo = com.tracuutiemchung.app.data.model.PatientInfo(fullName = "Bé MMR", phoneNumber = "0123456789", dateOfBirth = dob.toString())
        
        // Tiêm MVVAC lúc 9 tháng (2024-10-01)
        // Tiêm MMR-II lúc 12 tháng (2025-01-01)
        val records = listOf(
            com.tracuutiemchung.app.data.model.VaccinationRecord(vaccineName = "Sởi (MVVAC)", vaccinationDate = "2024-10-01"),
            com.tracuutiemchung.app.data.model.VaccinationRecord(vaccineName = "MMR-II", vaccinationDate = "2025-01-01")
        )
        val rules = loadRealRules()

        val result = engine.analyze(
            AnalysisInput(patient = patientInfo, vaccinationRecords = records, rules = rules, analysisDate = analysisDate)
        )

        // Sau 1 MVVAC + 1 MMR, cần thêm 1 MMR nữa để COMPLETED (tổng 3 mũi sởi)
        // Interval MMR sau MVVAC là 84 ngày. 
        // Sau đó mũi tiếp theo cách 1095 ngày (3 năm) hoặc theo phác đồ.
        // Trong logic checkMmrEquivalentGroup: mmrAfterMvvac.size == 1 -> cần thêm 1 mũi nữa, interval 1095 ngày.
        // 2025-01-01 + 1095 days = 2028-01-01.
        assertRuleStatus(result, "MMR", RecommendationStatus.DUE_LATER)
        val mmrRec = result.find { it.vaccineName.contains("MMR", ignoreCase = true) }
        assertEquals("2028-01-01", mmrRec?.suggestedDate)
    }

    @Test
    fun `golden test - Flu age-dependent logic`() {
        // Scenario 1: First dose at 3 years old -> 2 doses required
        val dob1 = LocalDate.of(2023, 1, 1)
        val analysisDate = LocalDate.of(2026, 5, 12)
        val patientInfo1 = com.tracuutiemchung.app.data.model.PatientInfo(fullName = "Bé Cúm Nhỏ", phoneNumber = "0123456789", dateOfBirth = dob1.toString())
        val records1 = emptyList<com.tracuutiemchung.app.data.model.VaccinationRecord>()
        val rules = loadRealRules()

        val result1 = engine.analyze(
            AnalysisInput(patient = patientInfo1, vaccinationRecords = records1, rules = rules, analysisDate = analysisDate)
        )
        val flu1 = result1.find { it.vaccineName.contains("Cúm", ignoreCase = true) }
        assertNotNull("Flu recommendation not found for Scenario 1", flu1)
        assert(flu1?.reason?.contains("Chưa tiêm") == true) { "Expected reason to contain 'Chưa tiêm' for no records, but got: ${flu1?.reason}" }

        // Scenario 2: First dose at 10 years old -> 1 dose required
        val dob2 = LocalDate.of(2016, 1, 1)
        val patientInfo2 = com.tracuutiemchung.app.data.model.PatientInfo(fullName = "Bé Cúm Lớn", phoneNumber = "0123456789", dateOfBirth = dob2.toString())
        val result2 = engine.analyze(
            AnalysisInput(patient = patientInfo2, vaccinationRecords = records1, rules = rules, analysisDate = analysisDate)
        )
        val flu2 = result2.find { it.vaccineName.contains("Cúm", ignoreCase = true) }
        assertNotNull("Flu recommendation not found for Scenario 2", flu2)
        // For >=9y, checkFluGroup says "Chưa tiêm vắc-xin Cúm. ..." or similar if no records
        // Wait, Scenario 2 has empty records too.
        assert(flu2?.reason?.contains("Chưa tiêm") == true) { "Expected reason to contain 'Chưa tiêm' for >=9y flu, but got: ${flu2?.reason}" }
    }

    private fun assertRuleStatus(results: List<com.tracuutiemchung.app.data.model.VaccineRecommendation>, ruleNamePart: String, expectedStatus: RecommendationStatus) {
        val rec = results.find { it.vaccineName.contains(ruleNamePart, ignoreCase = true) }
        if (rec == null) {
            val available = results.joinToString(", ") { it.vaccineName }
            fail("Rule matching '$ruleNamePart' not found in results. Available: [$available]")
            return
        }
        if (expectedStatus != rec.status) {
            println("FAILURE: ${rec.vaccineName} expected $expectedStatus but got ${rec.status}. Reason: ${rec.reason}. Tags: ${rec.statusTags}")
        }
        assertEquals("Status mismatch for ${rec.vaccineName}. Reason: ${rec.reason}", expectedStatus, rec.status)
    }
}
