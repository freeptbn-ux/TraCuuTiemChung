package com.tracuutiemchung.app.data.model

import kotlinx.serialization.Serializable
import java.time.LocalDate

@Serializable
data class PatientInfo(
    val fullName: String,
    val phoneNumber: String,
    val dateOfBirth: String? = null,
    val address: String? = null,
)

@Serializable
data class VaccinationRecord(
    val vaccineName: String,
    val doseNumber: Int? = null,
    val vaccinationDate: String,
    val provider: String? = null,
)

@Serializable
enum class RecommendationStatus {
    COMPLETED,
    DUE_NOW,
    DUE_LATER,
    OVERDUE,
    NEEDS_REVIEW,
    NOT_ENOUGH_DATA,
}

@Serializable
data class VaccineRecommendation(
    val vaccineName: String,
    val status: RecommendationStatus = RecommendationStatus.NEEDS_REVIEW,
    val suggestedDate: String? = null,
    val reason: String,
    val warnings: List<String> = emptyList(),
    val recommendedDate: String? = suggestedDate,
    val warning: String? = warnings.firstOrNull(),
    val statusTags: List<String> = emptyList(),
)

@Serializable
data class AnalysisResult(
    val patientInfo: PatientInfo,
    val records: List<VaccinationRecord>,
    val recommendations: List<VaccineRecommendation>,
    val warnings: List<String> = emptyList(),
)

data class ParsedRecord(
    val source: VaccinationRecord,
    val date: LocalDate,
)
