package com.tracuutiemchung.app.data.remote

import com.tracuutiemchung.app.data.model.*
import com.tracuutiemchung.app.data.portal.*
import com.tracuutiemchung.app.data.remote.model.AnalysisRequestDto
import com.tracuutiemchung.app.data.remote.model.RecommendationDto

class VercelPortalRepository(
    private val apiService: VercelApiService = RetrofitClient.apiService
) : PortalLookupRepository {

    override suspend fun searchPatientsByPhone(phone: String): Result<List<PortalPatientSummary>> {
        return try {
            val response = apiService.lookupPatients(
                request = mapOf("phone" to phone)
            )
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.status == "success") {
                    val patients = body.data ?: emptyList()
                    Result.success(patients.map { dto ->
                        PortalPatientSummary(
                            patientId = dto.id,
                            fullName = dto.name,
                            birthDateOrYear = dto.dob,
                            gender = dto.gender,
                            phone = phone
                        )
                    })
                } else {
                    Result.failure(PortalLookupException.NetworkFailed(body?.detail ?: "Unknown error"))
                }
            } else {
                Result.failure(PortalLookupException.NetworkFailed("Lỗi máy chủ: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(PortalLookupException.NetworkFailed("Lỗi kết nối: ${e.localizedMessage}"))
        }
    }

    override suspend fun lookupVaccinations(patient: PortalPatientSummary): Result<PortalLookupResult> {
        // We bypass the old flow and go straight to analysis
        return Result.failure(UnsupportedOperationException("Use analyze() for Vercel backend"))
    }

    suspend fun analyze(patientId: String, phone: String): Result<AnalysisResult> {
        return try {
            val response = apiService.analyzeVaccinations(
                request = AnalysisRequestDto(patientId, phone)
            )
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.status == "success") {
                    val data = body.data
                    val patientInfo = PatientInfo(
                        fullName = data.patientName,
                        phoneNumber = phone,
                        dateOfBirth = data.dob
                    )
                    val recommendations = data.missingVaccines.map { it.toDomain() }
                    val records = data.administeredVaccines.map {
                        VaccinationRecord(
                            vaccineName = it.vaccineName,
                            vaccinationDate = it.date,
                            doseNumber = it.dose?.toIntOrNull(),
                            provider = it.provider ?: "VNCDC"
                        )
                    }
                    
                    Result.success(AnalysisResult(
                        patientInfo = patientInfo,
                        records = records,
                        recommendations = recommendations,
                        warnings = emptyList()
                    ))
                } else {
                    Result.failure(PortalLookupException.NetworkFailed(body?.detail ?: "Không thể phân tích dữ liệu"))
                }
            } else {
                Result.failure(PortalLookupException.NetworkFailed("Lỗi máy chủ: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(PortalLookupException.NetworkFailed("Lỗi kết nối: ${e.localizedMessage}"))
        }
    }

    private fun RecommendationDto.toDomain(): VaccineRecommendation {
        return VaccineRecommendation(
            vaccineName = this.vaccineName,
            status = when (this.status) {
                "COMPLETED" -> RecommendationStatus.COMPLETED
                "DUE_NOW" -> RecommendationStatus.DUE_NOW
                "DUE_LATER" -> RecommendationStatus.DUE_LATER
                "OVERDUE" -> RecommendationStatus.OVERDUE
                "NEEDS_REVIEW" -> RecommendationStatus.NEEDS_REVIEW
                else -> RecommendationStatus.NEEDS_REVIEW
            },
            suggestedDate = this.nextDose,
            reason = this.message,
            statusTags = this.statusTags
        )
    }
}
