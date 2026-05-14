package com.tracuutiemchung.app.domain.analyzer

import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Test
import java.time.LocalDate

class VaccineDateUtilsTest {

    @Test
    fun testParseFlexibleDate() {
        // ISO format
        assertEquals(
            LocalDate.of(2023, 5, 20),
            VaccineDateUtils.parseFlexibleDate("2023-05-20")
        )
        
        // d/M/uuuu format
        assertEquals(
            LocalDate.of(2023, 5, 20),
            VaccineDateUtils.parseFlexibleDate("20/5/2023")
        )
        
        // dd/MM/uuuu format
        assertEquals(
            LocalDate.of(2023, 5, 20),
            VaccineDateUtils.parseFlexibleDate("20/05/2023")
        )
        
        // Trim whitespace
        assertEquals(
            LocalDate.of(2023, 5, 20),
            VaccineDateUtils.parseFlexibleDate("  2023-05-20  ")
        )
        
        // Invalid formats
        assertNull(VaccineDateUtils.parseFlexibleDate("2023.05.20"))
        assertNull(VaccineDateUtils.parseFlexibleDate("invalid"))
        assertNull(VaccineDateUtils.parseFlexibleDate(null))
        assertNull(VaccineDateUtils.parseFlexibleDate(""))
    }

    @Test
    fun testGetAgeAtDate() {
        val dob = LocalDate.of(2020, 1, 1)
        
        // 1 year exactly
        val date1y = LocalDate.of(2021, 1, 1)
        val age1y = VaccineDateUtils.getAgeAtDate(dob, date1y)
        assertEquals(1, age1y?.totalYears)
        assertEquals(12, age1y?.totalMonths)
        assertEquals(366, age1y?.totalDays) // 2020 is leap year
        
        // Just before 1 year
        val dateBefore1y = LocalDate.of(2020, 12, 31)
        val ageBefore1y = VaccineDateUtils.getAgeAtDate(dob, dateBefore1y)
        assertEquals(0, ageBefore1y?.totalYears)
        assertEquals(11, ageBefore1y?.totalMonths)
        
        // After 1 year and 2 months
        val date1y2m = LocalDate.of(2021, 3, 1)
        val age1y2m = VaccineDateUtils.getAgeAtDate(dob, date1y2m)
        assertEquals(1, age1y2m?.totalYears)
        assertEquals(14, age1y2m?.totalMonths)
    }
}
