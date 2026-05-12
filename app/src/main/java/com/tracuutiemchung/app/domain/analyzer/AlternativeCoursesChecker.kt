package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.ParsedRecord
import com.tracuutiemchung.app.data.rules.CourseConfig
import com.tracuutiemchung.app.data.rules.VaccineRule
import java.time.LocalDate

/**
 * Port of group_checkers_alternative.py
 */
fun checkAlternativeCoursesGroup(
    ruleKey: String,
    rule: VaccineRule,
    administeredRecords: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
    allRules: Map<String, VaccineRule>,
    allAdministered: Map<String, List<ParsedRecord>>,
): List<MissingItem> {
    val displayName = rule.displayName
    val courses = rule.courses

    if (courses.isEmpty()) return emptyList()

    // 1. Check if any course is already completed
    for (course in courses) {
        val courseRecords = filterRecordsByCourse(administeredRecords, course)
        if (courseRecords.size >= course.dosesRequired) {
            return emptyList() // Completed
        }
    }

    // Special JE_Group mixing logic (8B)
    if (ruleKey == "JE_Group") {
        return checkJeGroupMixing(rule, administeredRecords, dob, analysisDate)
    }

    // Special Rota logic (8A)
    if (ruleKey == "Rota") {
        return checkRotaLogic(rule, administeredRecords, dob, analysisDate)
    }

    // Default Alternative Courses Logic (e.g. HepA)
    return checkDefaultAlternativeCourses(rule, administeredRecords, dob, analysisDate)
}

private fun filterRecordsByCourse(
    records: List<ParsedRecord>,
    course: CourseConfig
): List<ParsedRecord> {
    if (course.rawNames.isEmpty()) return records
    val normalizedCourseNames = course.rawNames.map { it.lowercase() }
    return records.filter { record ->
        val recordName = record.source.vaccineName.lowercase()
        normalizedCourseNames.any { it == recordName || recordName.contains(it) }
    }
}

private fun checkRotaLogic(
    rule: VaccineRule,
    administered: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
): List<MissingItem> {
    val displayName = rule.displayName
    if (dob == null) {
        return listOf(
            MissingItem(
                displayName,
                "Thiếu ngày sinh để kiểm tra Rota",
                null,
                listOf("error_dob")
            )
        )
    }

    val ageAtAnalysis = VaccineDateUtils.getAgeAtDate(dob, analysisDate) ?: return emptyList()

    // Check max age to start (6 months)
    val maxAgeToStartMonths = rule.maxAgeMonthsToStart ?: 6
    val maxDateToStart = VaccineDateUtils.addMonths(dob, maxAgeToStartMonths)
    
    if (administered.isEmpty()) {
        if (analysisDate.isAfter(maxDateToStart)) {
            return listOf(
                MissingItem(
                    displayName,
                    "Quá tuổi bắt đầu uống Rota (tối đa $maxAgeToStartMonths tháng).",
                    null,
                    listOf("too_old_to_start")
                )
            )
        }
        
        // Suggest first dose
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        return listOf(
            MissingItem(
                displayName,
                ageStatus.message,
                ageStatus.earliestDate,
                ageStatus.statusTags
            )
        )
    }

    // Check max age for completion (8 months)
    val maxAgeForCompletionMonths = rule.maxAgeMonthsForCompletion ?: 8
    val maxDateForCompletion = VaccineDateUtils.addMonths(dob, maxAgeForCompletionMonths)
    
    if (analysisDate.isAfter(maxDateForCompletion)) {
        return listOf(
            MissingItem(
                displayName,
                "Quá tuổi hoàn thành Rota (tối đa $maxAgeForCompletionMonths tháng).",
                null,
                listOf("too_old_to_complete")
            )
        )
    }

    // Select course based on first dose or administered records
    // For Rota, usually we follow the first dose's course
    val firstDoseName = administered.first().source.vaccineName.lowercase()
    val matchedCourse = rule.courses.find { course ->
        course.rawNames.any { it.lowercase().let { name -> firstDoseName == name || firstDoseName.contains(name) } }
    } ?: rule.courses.first() // Fallback to first course

    val courseRecords = filterRecordsByCourse(administered, matchedCourse)
    val dosesRequired = matchedCourse.dosesRequired
    val validDosesCount = courseRecords.size
    
    if (validDosesCount >= dosesRequired) return emptyList()

    val nextDoseNumber = validDosesCount + 1
    val lastDoseDate = courseRecords.lastOrNull()?.date ?: administered.last().date
    val interval = matchedCourse.minIntervalDays.getOrNull(validDosesCount) ?: 28 // Default 4 weeks
    val earliestNextDate = lastDoseDate.plusDays(interval.toLong())

    val isDue = analysisDate.isAfter(earliestNextDate) || analysisDate.isEqual(earliestNextDate)
    
    return listOf(
        MissingItem(
            displayName,
            "Đã uống $validDosesCount/$dosesRequired liều ${matchedCourse.display}. Cần uống liều tiếp theo.",
            earliestNextDate,
            if (isDue) listOf("due") else listOf("scheduled")
        )
    )
}

private fun checkJeGroupMixing(
    rule: VaccineRule,
    administered: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
): List<MissingItem> {
    val displayName = rule.displayName
    val jevaxRecords = administered.filter { it.source.vaccineName.lowercase().contains("jevax") }
    val imojevRecords = administered.filter { it.source.vaccineName.lowercase().contains("imojev") }
    val jeevRecords = administered.filter { it.source.vaccineName.lowercase().contains("jeev") }

    // Logic 21: Imojev trước Jevax: error interchange
    if (imojevRecords.isNotEmpty() && jevaxRecords.isNotEmpty()) {
        val firstImojev = imojevRecords.minByOrNull { it.date }!!
        val firstJevax = jevaxRecords.minByOrNull { it.date }!!
        if (firstImojev.date.isBefore(firstJevax.date)) {
            return listOf(
                MissingItem(
                    displayName,
                    "Cảnh báo: Đã tiêm Imojev trước Jevax. Không khuyến cáo chuyển đổi ngược lại.",
                    null,
                    listOf("error_interchange")
                )
            )
        }
    }

    // Logic 22: Jevax >= 3 mũi + Imojev >= 1: hoàn thành
    if (jevaxRecords.size >= 3 && (imojevRecords.isNotEmpty() || jeevRecords.isNotEmpty())) {
        return emptyList()
    }

    // Logic 23 & 24 & 25
    if (jevaxRecords.isNotEmpty() && imojevRecords.isEmpty() && jeevRecords.isEmpty()) {
        if (jevaxRecords.size >= 3) {
            // Gợi ý booster Jevax hoặc 1 mũi Imojev
            val lastJevax = jevaxRecords.last().date
            val nextBooster = lastJevax.plusYears(3)
            val isDue = analysisDate.isAfter(nextBooster) || analysisDate.isEqual(nextBooster)
            return listOf(
                MissingItem(
                    displayName,
                    "Đã tiêm phác đồ cơ bản Jevax. Cần tiêm nhắc sau 3 năm hoặc chuyển sang 1 mũi Imojev.",
                    nextBooster,
                    if (isDue) listOf("due", "booster_due") else listOf("info", "booster_upcoming")
                )
            )
        } else {
            // Jevax 1-2 mũi: gợi ý chuyển sang Imojev
            return listOf(
                MissingItem(
                    displayName,
                    "Đã tiêm ${jevaxRecords.size} mũi Jevax. Khuyến cáo chuyển sang tiêm Imojev để đạt hiệu quả cao hơn.",
                    analysisDate,
                    listOf("due", "switch_imojev")
                )
            )
        }
    }

    // Imojev / JEEV completion check
    for (course in rule.courses) {
        val courseRecords = filterRecordsByCourse(administered, course)
        if (courseRecords.size >= course.dosesRequired) return emptyList()
    }

    // Default suggest next dose for Imojev/JEEV/Jevax
    val mostRecent = administered.maxByOrNull { it.date }
    if (mostRecent == null) {
        val ageStatus = AnalysisRuleUtils.getAgeStatusAndEarliestDate(dob, analysisDate, rule, displayName)
        return listOf(MissingItem(displayName, ageStatus.message, ageStatus.earliestDate, ageStatus.statusTags))
    }

    // Simple fallback: suggest Imojev if not finished
    return listOf(
        MissingItem(
            displayName,
            "Cần tiêm phòng Nhật Bản B (Khuyến cáo dùng Imojev).",
            analysisDate,
            listOf("due")
        )
    )
}

private fun checkDefaultAlternativeCourses(
    rule: VaccineRule,
    administered: List<ParsedRecord>,
    dob: LocalDate?,
    analysisDate: LocalDate,
): List<MissingItem> {
    val displayName = rule.displayName
    
    // Select course by age or by existing doses
    val ageAtAnalysis = if (dob != null) VaccineDateUtils.getAgeAtDate(dob, analysisDate) else null
    
    val selectedCourse = if (administered.isNotEmpty()) {
        val firstDoseName = administered.first().source.vaccineName.lowercase()
        rule.courses.find { course ->
            course.rawNames.any { name -> firstDoseName.contains(name.lowercase()) }
        } ?: rule.courses.first()
    } else if (ageAtAnalysis != null) {
        rule.courses.find { course ->
            val minAge = course.minAgeMonthsAtFirstDose ?: 0
            val maxAge = course.maxAgeYearsAtFirstDose?.let { it * 12 } ?: Int.MAX_VALUE
            ageAtAnalysis.totalMonths >= minAge && ageAtAnalysis.totalMonths <= maxAge
        } ?: rule.courses.first()
    } else {
        rule.courses.first()
    }

    val courseRecords = filterRecordsByCourse(administered, selectedCourse)
    if (courseRecords.size >= selectedCourse.dosesRequired) return emptyList()

    val nextDoseNumber = courseRecords.size + 1
    val lastDate = courseRecords.lastOrNull()?.date
    val interval = selectedCourse.minIntervalDays.getOrNull(courseRecords.size) ?: 180 // Default 6 months for HepA
    
    val earliestDate = lastDate?.plusDays(interval.toLong()) ?: analysisDate
    val isDue = analysisDate.isAfter(earliestDate) || analysisDate.isEqual(earliestDate)
    
    return listOf(
        MissingItem(
            displayName,
            "Cần tiêm ${selectedCourse.display} (Mũi $nextDoseNumber).",
            earliestDate,
            if (isDue) listOf("due") else listOf("scheduled")
        )
    )
}
