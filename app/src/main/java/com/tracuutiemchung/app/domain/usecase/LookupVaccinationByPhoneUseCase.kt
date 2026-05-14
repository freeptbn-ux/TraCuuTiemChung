package com.tracuutiemchung.app.domain.usecase

import com.tracuutiemchung.app.data.model.AnalysisResult
import com.tracuutiemchung.app.data.portal.PortalPatientSummary
import com.tracuutiemchung.app.data.remote.VercelPortalRepository

class LookupVaccinationByPhoneUseCase(
    private val repository: VercelPortalRepository,
) {
    suspend fun searchPatients(phone: String): Result<List<PortalPatientSummary>> =
        repository.searchPatientsByPhone(phone)

    suspend fun lookupPatient(patient: PortalPatientSummary): Result<AnalysisResult> =
        repository.analyze(patient.patientId, patient.phone.orEmpty())

    // Backward compatibility or legacy support if needed
    suspend operator fun invoke(phone: String): Result<AnalysisResult> = 
        searchPatients(phone).mapCatching { patients ->
            val patient = patients.firstOrNull() ?: throw Exception("Không tìm thấy bệnh nhân")
            lookupPatient(patient).getOrThrow()
        }
}
