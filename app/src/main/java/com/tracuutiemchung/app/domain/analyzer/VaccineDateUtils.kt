package com.tracuutiemchung.app.domain.analyzer

import java.time.LocalDate
import java.time.format.DateTimeFormatter
import java.time.format.DateTimeParseException
import java.time.temporal.ChronoUnit

data class AgeAtDate(
    val totalMonths: Int,
    val totalDays: Int,
    val totalYears: Int,
)

object VaccineDateUtils {
    private val supportedDateFormats = listOf(
        DateTimeFormatter.ISO_LOCAL_DATE,
        DateTimeFormatter.ofPattern("d/M/uuuu"),
        DateTimeFormatter.ofPattern("dd/MM/uuuu"),
    )

    fun getAgeAtDate(dob: LocalDate, atDate: LocalDate): AgeAtDate? {
        if (atDate.isBefore(dob)) return null

        val totalDays = ChronoUnit.DAYS.between(dob, atDate).toInt()

        var years = atDate.year - dob.year
        if (atDate.monthValue < dob.monthValue || (atDate.monthValue == dob.monthValue && atDate.dayOfMonth < dob.dayOfMonth)) {
            years--
        }
        val totalYears = maxOf(0, years)

        var monthsTotal = (atDate.year - dob.year) * 12 + (atDate.monthValue - dob.monthValue)
        if (atDate.dayOfMonth < dob.dayOfMonth) {
            monthsTotal--
        }
        val totalMonths = maxOf(0, monthsTotal)

        return AgeAtDate(
            totalMonths = totalMonths,
            totalDays = totalDays,
            totalYears = totalYears
        )
    }

    fun addMonths(source: LocalDate, months: Int): LocalDate {
        return source.plusMonths(months.toLong())
    }

    fun addYears(source: LocalDate, years: Int): LocalDate {
        return source.plusYears(years.toLong())
    }

    fun parseFlexibleDate(dateStr: String?): LocalDate? {
        if (dateStr.isNullOrBlank()) return null
        return supportedDateFormats.firstNotNullOfOrNull { formatter ->
            try {
                LocalDate.parse(dateStr.trim(), formatter)
            } catch (_: DateTimeParseException) {
                null
            }
        }
    }
}
