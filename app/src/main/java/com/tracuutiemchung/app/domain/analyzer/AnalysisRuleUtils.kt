package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.RecommendationStatus
import com.tracuutiemchung.app.data.model.VaccineRecommendation
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate
import java.time.format.DateTimeFormatter

data class AgeStatus(
    val message: String,
    val earliestDate: LocalDate?,
    val statusTags: List<String>,
)

data class MissingItem(
    val vaccineNameForPopup: String,
    val description: String,
    val earliestNextDoseDate: LocalDate?,
    val statusTags: List<String>,
)

fun MissingItem.toRecommendation(analysisDate: LocalDate): VaccineRecommendation {
    val status = when {
        statusTags.any { it.startsWith("error_") && it != "error_dob" } -> RecommendationStatus.NEEDS_REVIEW
        statusTags.contains("error_dob") -> RecommendationStatus.NOT_ENOUGH_DATA
        statusTags.any { it.contains("due") && !it.contains("later") } || 
        statusTags.contains("due") || 
        statusTags.contains("booster_due") ||
        statusTags.contains("eligible") ||
        statusTags.contains("overdue") -> RecommendationStatus.DUE_NOW
        statusTags.any { it.contains("later") } ||
        statusTags.contains("too_young") || 
        statusTags.any { it.startsWith("too_old_") } ||
        statusTags.contains("booster_upcoming") ||
        statusTags.contains("info") ||
        statusTags.contains("scheduled") ||
        statusTags.contains("upcoming") -> RecommendationStatus.DUE_LATER
        else -> RecommendationStatus.NEEDS_REVIEW
    }

    val dateStr = earliestNextDoseDate?.format(DateTimeFormatter.ISO_LOCAL_DATE)

    val warnings = mutableListOf<String>()
    if (status == RecommendationStatus.NEEDS_REVIEW || statusTags.contains("warning")) {
        warnings.add(description)
    }
    if (statusTags.contains("too_early")) {
        warnings.add("Cảnh báo: Tiêm quá sớm so với quy định.")
    }

    return VaccineRecommendation(
        vaccineName = vaccineNameForPopup,
        status = status,
        suggestedDate = dateStr,
        reason = description,
        warnings = warnings,
        statusTags = statusTags
    )
}

object AnalysisRuleUtils {
    private const val GRACE_PERIOD_DAYS = 0L

    fun getAgeStatusAndEarliestDate(
        dob: LocalDate?,
        analysisDate: LocalDate,
        rule: VaccineRule,
        displayName: String,
    ): AgeStatus {
        val prefix = if (displayName.isNotEmpty()) "$displayName - " else ""

        if (dob == null) {
            return AgeStatus("${prefix}Không có ngày sinh để kiểm tra tuổi", null, listOf("error_dob"))
        }

        val age = VaccineDateUtils.getAgeAtDate(dob, analysisDate)
        if (age == null) {
            return AgeStatus("${prefix}Ngày phân tích trước ngày sinh", null, listOf("error_date"))
        }

        val minAgeMonths = rule.minAgeMonthsAtFirstDose ?: rule.minAgeMonthsOverall
        val minAgeWeeks = rule.minAgeWeeksAtFirstDose
        val minAgeYears = rule.minAgeYearsAtFirstDose
        val minAgeDaysVal = rule.minAgeDaysAtFirstDose
        val minAgeMonthsGroup = rule.minAgeMonthsOverallGroup

        // Python's _get_age_status_and_earliest_date priority: months > years > weeks > days
        var earliestAcceptableDate: LocalDate? = null
        var statusMessage = ""
        var statusTags = listOf("eligible")

        val effectiveMinAgeMonths = minAgeMonths ?: minAgeMonthsGroup
        if (effectiveMinAgeMonths != null) {
            val target = VaccineDateUtils.addMonths(dob, effectiveMinAgeMonths)
            earliestAcceptableDate = target.minusDays(GRACE_PERIOD_DAYS)
            if (analysisDate.isBefore(earliestAcceptableDate)) {
                statusMessage = "cần $effectiveMinAgeMonths tháng tuổi"
                statusTags = listOf("too_young")
            }
        } else if (minAgeYears != null) {
            val target = VaccineDateUtils.addYears(dob, minAgeYears)
            earliestAcceptableDate = target.minusDays(GRACE_PERIOD_DAYS)
            if (analysisDate.isBefore(earliestAcceptableDate)) {
                statusMessage = "cần $minAgeYears tuổi"
                statusTags = listOf("too_young")
            }
        } else if (minAgeWeeks != null) {
            val effectiveDays = (minAgeWeeks * 7).toLong() - GRACE_PERIOD_DAYS
            if (age.totalDays < effectiveDays) {
                earliestAcceptableDate = dob.plusDays(effectiveDays)
                statusMessage = "cần $minAgeWeeks tuần tuổi"
                statusTags = listOf("too_young")
            }
        } else if (minAgeDaysVal != null) {
            val effectiveDays = minAgeDaysVal.toLong() - GRACE_PERIOD_DAYS
            if (age.totalDays < effectiveDays) {
                earliestAcceptableDate = dob.plusDays(effectiveDays)
                val displayReq = if (minAgeDaysVal >= 60) "${minAgeDaysVal / 30} tháng" else "$minAgeDaysVal ngày"
                statusMessage = "cần >$displayReq tuổi"
                statusTags = listOf("too_young")
            }
        }

        if (statusTags.contains("too_young")) {
            return AgeStatus(statusMessage, earliestAcceptableDate, statusTags)
        }

        return AgeStatus("đủ điều kiện tuổi", analysisDate, statusTags)
    }

    fun checkFirstDoseAgeValidity(
        dob: LocalDate?,
        firstDoseDate: LocalDate,
        rule: VaccineRule,
        displayName: String,
    ): Pair<Boolean, MissingItem?> {
        if (dob == null) return true to null

        val age = VaccineDateUtils.getAgeAtDate(dob, firstDoseDate)
        if (age == null) {
            return false to MissingItem(
                vaccineNameForPopup = displayName,
                description = "$displayName - Lỗi tính tuổi cho mũi đầu (ngày tiêm có thể trước ngày sinh).",
                earliestNextDoseDate = null,
                statusTags = listOf("error_age_calculation")
            )
        }

        val minAgeMonths = rule.minAgeMonthsAtFirstDose ?: rule.minAgeMonthsOverall
        val minAgeWeeks = rule.minAgeWeeksAtFirstDose
        val minAgeYears = rule.minAgeYearsAtFirstDose
        val minAgeDaysVal = rule.minAgeDaysAtFirstDose

        var errorDetail: String? = null

        // Priority from Phase 03 doc (Days > Weeks > Months > Years)
        if (minAgeDaysVal != null) {
            val effectiveMinDays = minAgeDaysVal.toLong() - GRACE_PERIOD_DAYS
            if (age.totalDays < effectiveMinDays) {
                val displayAge = if (minAgeDaysVal >= 60) "${minAgeDaysVal / 30} tháng" else "$minAgeDaysVal ngày"
                errorDetail = "Mũi 1 tiêm quá sớm (cần >$displayAge, thực tế ${age.totalDays} ngày tuổi)."
            }
        } else if (minAgeWeeks != null) {
            val effectiveMinDays = (minAgeWeeks * 7).toLong() - GRACE_PERIOD_DAYS
            if (age.totalDays < effectiveMinDays) {
                errorDetail = "Mũi 1 tiêm quá sớm (cần $minAgeWeeks tuần, thực tế ${age.totalDays} ngày tuổi)."
            }
        } else if (minAgeMonths != null) {
            val earliestAllowed = VaccineDateUtils.addMonths(dob, minAgeMonths).minusDays(GRACE_PERIOD_DAYS)
            if (firstDoseDate.isBefore(earliestAllowed)) {
                errorDetail = "Mũi 1 tiêm quá sớm (cần $minAgeMonths tháng tuổi)."
            }
        } else if (minAgeYears != null) {
            val earliestAllowed = VaccineDateUtils.addYears(dob, minAgeYears).minusDays(GRACE_PERIOD_DAYS)
            if (firstDoseDate.isBefore(earliestAllowed)) {
                errorDetail = "Mũi 1 tiêm quá sớm (cần $minAgeYears tuổi)."
            }
        }

        return if (errorDetail != null) {
            false to MissingItem(
                vaccineNameForPopup = displayName,
                description = "$displayName - $errorDetail",
                earliestNextDoseDate = null,
                statusTags = listOf("error_age_first_dose", "too_early")
            )
        } else {
            true to null
        }
    }
}
