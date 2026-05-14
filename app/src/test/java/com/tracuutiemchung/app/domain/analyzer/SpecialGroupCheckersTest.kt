package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.rules.AgeBasedRegimen
import com.tracuutiemchung.app.data.rules.BoosterConfig
import com.tracuutiemchung.app.data.rules.InteractionConfig
import com.tracuutiemchung.app.data.rules.MemberConfig
import com.tracuutiemchung.app.data.rules.VaccineRule
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test
import java.time.LocalDate

class SpecialGroupCheckersTest {

    private val dob = LocalDate.of(2023, 1, 1)
    private val analysisDate = LocalDate.of(2024, 5, 1)

    @Test
    fun testFluGroup_NoDoses() {
        val rule = VaccineRule(
            vaccineKey = "Flu",
            displayName = "Cúm",
            type = "flu_group",
            minAgeMonthsAtFirstDose = 6
        )
        val results = checkFluGroup("Flu", rule, emptyList(), dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("Chưa tiêm"))
        assertEquals(analysisDate, results[0].earliestNextDoseDate)
    }

    @Test
    fun testFluGroup_Under9Years_OneDose() {
        val rule = VaccineRule(
            vaccineKey = "Flu",
            displayName = "Cúm",
            type = "flu_group"
        )
        val records = listOf(
            createParsedRecord("Cúm", dob.plusMonths(7))
        )
        val results = checkFluGroup("Flu", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("cần 2 mũi"))
        assertEquals(records[0].date.plusDays(30), results[0].earliestNextDoseDate)
    }

    @Test
    fun testFluGroup_AnnualBooster() {
        val rule = VaccineRule(
            vaccineKey = "Flu",
            displayName = "Cúm",
            type = "flu_group"
        )
        // >= 9 years old at first dose (simulated by setting DOB far back)
        val adultDob = LocalDate.of(2000, 1, 1)
        val records = listOf(
            createParsedRecord("Cúm", LocalDate.of(2023, 1, 1))
        )
        val results = checkFluGroup("Flu", rule, records, adultDob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("nhắc lại hàng năm"))
        assertEquals(LocalDate.of(2024, 1, 1), results[0].earliestNextDoseDate)
        assertTrue(results[0].statusTags.contains("booster_due"))
    }

    @Test
    fun testCumulativeGroup() {
        val rule = VaccineRule(
            vaccineKey = "Group_ABC",
            displayName = "Nhóm ABC",
            type = "group_cumulative_unique_doses",
            requiredDoses = 3,
            minIntervalDays = listOf(null, 30, 30)
        )
        val records = listOf(
            createParsedRecord("A", dob.plusMonths(2)),
            createParsedRecord("B", dob.plusMonths(3))
        )
        val results = checkCumulativeGroupDoses("Group_ABC", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("2/3 liều"))
        assertEquals(records.last().date.plusDays(30), results[0].earliestNextDoseDate)
    }

    @Test
    fun testMeningococcalAcyw_MenQuadfi_Infant() {
        val rule = createMeningococcalRule()
        val records = listOf(
            createParsedRecord("MenQuadfi", dob.plusMonths(2))
        )
        val results = checkMeningococcalAcywGroup("MeningococcalACYW_Group", rule, records, dob, analysisDate, emptyMap(), emptyMap())
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("Đã tiêm 1/3 mũi MenQuadfi"))
    }

    @Test
    fun testMeningococcalAcyw_Menactra_Older() {
        val rule = createMeningococcalRule()
        val records = listOf(
            createParsedRecord("MENACTRA", dob.plusMonths(25))
        )
        // Analysis date needs to be later
        val futureAnalysis = dob.plusMonths(26)
        val results = checkMeningococcalAcywGroup("MeningococcalACYW_Group", rule, records, dob, futureAnalysis, emptyMap(), emptyMap())
        // Should be complete as first dose was at 25 months
        assertEquals(0, results.size)
    }

    @Test
    fun testMeningococcalAcyw_Interaction_SixInOne() {
        val rule = createMeningococcalRule()
        val sixInOneRecords = listOf(
            createParsedRecord("6in1", analysisDate.minusDays(10))
        )
        val allAdministered = mapOf("Six_In_One_Combined" to sixInOneRecords)
        
        val results = checkMeningococcalAcywGroup("MeningococcalACYW_Group", rule, emptyList(), dob, analysisDate, emptyMap(), allAdministered)
        assertEquals(1, results.size)
        assertTrue(results[0].description.contains("MenQuadfi nên cách ít nhất 1 tháng"))
    }

    private fun createParsedRecord(name: String, date: LocalDate): ParsedRecord {
        return ParsedRecord(
            source = VaccinationRecord(vaccineName = name, vaccinationDate = date.toString()),
            date = date
        )
    }

    private fun createMeningococcalRule(): VaccineRule {
        return VaccineRule(
            vaccineKey = "MeningococcalACYW_Group",
            displayName = "Não mô cầu ACYW",
            type = "meningococcal_acyw_group",
            members = mapOf(
                "MENACTRA" to MemberConfig(
                    rawNames = listOf("MENACTRA"),
                    display = "Menactra",
                    rulesByAge = listOf(
                        AgeBasedRegimen(maxAgeAtFirstDoseMonths = 23, dosesRequired = 2, minIntervalDays = listOf(null, 90)),
                        AgeBasedRegimen(minAgeAtFirstDoseMonths = 24, dosesRequired = 1)
                    )
                ),
                "MENQUADFI" to MemberConfig(
                    rawNames = listOf("MenQuadfi"),
                    display = "MenQuadfi",
                    rulesByAge = listOf(
                        AgeBasedRegimen(minAgeWeeksAtFirstDose = 6, maxAgeAtFirstDoseMonths = 5, dosesRequired = 3, minIntervalDays = listOf(null, 60, 60),
                            booster = BoosterConfig(minAgeMonths = 12, minIntervalDaysFromLast = 60, description = "Mũi nhắc MenQuadfi")),
                        AgeBasedRegimen(minAgeAtFirstDoseMonths = 6, maxAgeAtFirstDoseMonths = 11, dosesRequired = 1,
                            booster = BoosterConfig(minAgeMonths = 12, minIntervalDaysFromLast = 60, description = "Mũi nhắc MenQuadfi")),
                        AgeBasedRegimen(minAgeAtFirstDoseMonths = 12, dosesRequired = 1)
                    )
                )
            ),
            interactions = mapOf(
                "Six_In_One_Combined" to InteractionConfig(minIntervalDays = 30, message = "MenQuadfi nên cách ít nhất 1 tháng.")
            )
        )
    }
}
