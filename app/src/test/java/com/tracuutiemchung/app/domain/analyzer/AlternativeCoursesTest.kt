package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.rules.CourseConfig
import com.tracuutiemchung.app.data.rules.VaccineRule
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import java.time.LocalDate

class AlternativeCoursesTest {

    private val dob = LocalDate.of(2023, 1, 1)
    private val analysisDate = LocalDate.of(2023, 10, 1)

    @Test
    fun testRota_NoDoses_Under6Weeks() {
        val rule = createRotaRule()
        val earlyAnalysis = dob.plusWeeks(2)
        val results = checkAlternativeCoursesGroup("Rota", rule, emptyList(), dob, earlyAnalysis, emptyMap(), emptyMap())
        println("testRota_NoDoses_Under6Weeks results: ${results.getOrNull(0)?.description}")
        assertEquals(1, results.size)
        assertTrue("Description was: ${results[0].description}", results[0].description.contains("cần 6 tuần tuổi"))
    }

    @Test
    fun testRota_NoDoses_Over6Months() {
        val rule = createRotaRule()
        val results = checkAlternativeCoursesGroup("Rota", rule, emptyList(), dob, dob.plusMonths(7), emptyMap(), emptyMap())
        println("testRota_NoDoses_Over6Months results: ${results.getOrNull(0)?.description}")
        assertEquals(1, results.size)
        assertTrue(results[0].statusTags.contains("too_old_to_start"))
    }

    @Test
    fun testRota_Rotarix_Completed() {
        val rule = createRotaRule()
        val records = listOf(
            createParsedRecord("Rotarix", dob.plusWeeks(6)),
            createParsedRecord("Rotarix", dob.plusWeeks(10))
        )
        val results = checkAlternativeCoursesGroup("Rota", rule, records, dob, dob.plusMonths(4), emptyMap(), emptyMap())
        assertEquals(0, results.size) // Completed
    }

    @Test
    fun testRota_Rotavin_NextDose() {
        val rule = createRotaRule()
        val records = listOf(
            createParsedRecord("Rotavin", dob.plusWeeks(6))
        )
        // Use an analysis date within 8 months
        val results = checkAlternativeCoursesGroup("Rota", rule, records, dob, dob.plusMonths(4), emptyMap(), emptyMap())
        println("testRota_Rotavin_NextDose results: ${results.getOrNull(0)?.description}")
        assertEquals(1, results.size)
        assertTrue("Description was: ${results[0].description}", results[0].description.contains("1/2 liều Rotavin"))
    }

    @Test
    fun testJE_JevaxToImojevSwitch() {
        val rule = createJeRule()
        val records = listOf(
            createParsedRecord("Jevax", dob.plusMonths(12)),
            createParsedRecord("Jevax", dob.plusMonths(12).plusDays(14))
        )
        val results = checkAlternativeCoursesGroup("JE_Group", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("chuyển sang tiêm Imojev"))
        assertTrue(results[0].statusTags.contains("switch_imojev"))
    }

    @Test
    fun testJE_ImojevBeforeJevaxError() {
        val rule = createJeRule()
        val records = listOf(
            createParsedRecord("Imojev", dob.plusMonths(12)),
            createParsedRecord("Jevax", dob.plusMonths(13))
        )
        val results = checkAlternativeCoursesGroup("JE_Group", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].statusTags.contains("error_interchange"))
    }

    @Test
    fun testHepA_AgeBasedSelection() {
        val rule = createHepARule()
        // Child under 15 years should get Avaxim 80
        val results = checkAlternativeCoursesGroup("HepA", rule, emptyList(), dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("Avaxim 80"))
    }

    @Test
    fun testPneumococcal_MixingWarning() {
        val rule = VaccineRule(vaccineKey = "PNEU_13", displayName = "Phế cầu 13", type = "pneumococcal_special")
        val allAdministered = mapOf(
            "PNEU_13" to listOf(createParsedRecord("Prevenar 13", dob.plusMonths(2))),
            "PNEU_10" to listOf(createParsedRecord("Synflorix", dob.plusMonths(4)))
        )
        val results = checkPneumococcalSpecial("PNEU_13", rule, allAdministered["PNEU_13"]!!, dob, analysisDate, emptyMap(), allAdministered)
        assertEquals(1, results.size)
        assertTrue(results[0].statusTags.contains("interleaved_pneu"))
    }

    @Test
    fun testPneumovax23_Suggestion() {
        val pneu13Rule = VaccineRule(vaccineKey = "PNEU_13", displayName = "Phế cầu 13", type = "age_dependent_series", requiredDoses = 4)
        val pneu23Rule = VaccineRule(vaccineKey = "PNEU_23", displayName = "Phế cầu 23", type = "pneumococcal_special")
        
        val allRules = mapOf("PNEU_13" to pneu13Rule, "PNEU_23" to pneu23Rule)
        val allAdministered = mapOf(
            "PNEU_13" to listOf(
                createParsedRecord("P13", dob.plusMonths(2)),
                createParsedRecord("P13", dob.plusMonths(4)),
                createParsedRecord("P13", dob.plusMonths(6)),
                createParsedRecord("P13", dob.plusMonths(12))
            ),
            "PNEU_23" to emptyList<ParsedRecord>()
        )
        
        // Child is 3 years old
        val results = checkPneumococcalSpecial("PNEU_23", pneu23Rule, emptyList(), dob, dob.plusYears(3), allRules, allAdministered)
        assertEquals(1, results.size)
        assertTrue(results[0].statusTags.contains("booster_pneu23"))
    }

    private fun createParsedRecord(name: String, date: LocalDate): ParsedRecord {
        return ParsedRecord(
            source = VaccinationRecord(vaccineName = name, vaccinationDate = date.toString()),
            date = date
        )
    }

    private fun createRotaRule(): VaccineRule {
        return VaccineRule(
            vaccineKey = "Rota",
            displayName = "Rota",
            type = "group_alternative_courses",
            minAgeWeeksAtFirstDose = 6,
            maxAgeMonthsToStart = 6,
            maxAgeMonthsForCompletion = 8,
            courses = listOf(
                CourseConfig(rawNames = listOf("Rotarix"), dosesRequired = 2, display = "Rotarix"),
                CourseConfig(rawNames = listOf("Rotavin"), dosesRequired = 2, display = "Rotavin"),
                CourseConfig(rawNames = listOf("RotaTeq"), dosesRequired = 3, display = "RotaTeq")
            )
        )
    }

    private fun createJeRule(): VaccineRule {
        return VaccineRule(
            vaccineKey = "JE_Group",
            displayName = "Nhật Bản B",
            type = "group_alternative_courses_age_range",
            courses = listOf(
                CourseConfig(rawNames = listOf("Imojev"), dosesRequired = 2, display = "Imojev"),
                CourseConfig(rawNames = listOf("JEEV"), dosesRequired = 2, display = "JEEV"),
                CourseConfig(rawNames = listOf("Jevax"), dosesRequired = 3, display = "Jevax")
            )
        )
    }

    private fun createHepARule(): VaccineRule {
        return VaccineRule(
            vaccineKey = "HepA",
            displayName = "Viêm gan A",
            type = "group_alternative_courses_age_range",
            courses = listOf(
                CourseConfig(rawNames = listOf("Avaxim 80"), dosesRequired = 2, display = "Avaxim 80", maxAgeYearsAtFirstDose = 15),
                CourseConfig(rawNames = listOf("Avaxim 160"), dosesRequired = 2, display = "Avaxim 160", minAgeMonthsAtFirstDose = 180)
            )
        )
    }
}
