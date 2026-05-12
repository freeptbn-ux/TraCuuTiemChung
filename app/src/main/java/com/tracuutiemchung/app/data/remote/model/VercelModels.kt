package com.tracuutiemchung.app.data.remote.model

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable

@Serializable
data class VercelResponse<T>(
    val status: String,
    val data: T,
    val detail: String? = null
)

@Serializable
data class PatientLookupDto(
    val id: String,
    val name: String,
    val dob: String,
    val gender: String
)

@Serializable
data class AnalysisRequestDto(
    @SerialName("patient_id") val patientId: String,
    val phone: String
)

@Serializable
data class RecommendationDto(
    @SerialName("vaccine_name") val vaccineName: String,
    @SerialName("rule_type") val ruleType: String,
    val status: String,
    @SerialName("next_dose") val nextDose: String? = null,
    val message: String,
    @SerialName("status_tags") val statusTags: List<String> = emptyList()
)

@Serializable
data class RecordDto(
    @SerialName("vaccine_name") val vaccineName: String,
    val date: String,
    val dose: String? = null,
    val provider: String? = null
)

@Serializable
data class AnalysisResponseDto(
    @SerialName("patient_name") val patientName: String,
    val dob: String,
    @SerialName("analysis_date") val analysisDate: String,
    @SerialName("missing_vaccines") val missingVaccines: List<RecommendationDto>,
    @SerialName("administered_vaccines") val administeredVaccines: List<RecordDto> = emptyList()
)
