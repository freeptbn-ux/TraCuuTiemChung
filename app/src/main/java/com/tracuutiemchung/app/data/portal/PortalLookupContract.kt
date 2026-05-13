package com.tracuutiemchung.app.data.portal

import com.tracuutiemchung.app.data.model.PatientInfo
import com.tracuutiemchung.app.data.model.VaccinationRecord

data class PortalLookupResult(
    val patient: PatientInfo?,
    val vaccinations: List<VaccinationRecord>,
    val rawSourceType: RawSourceType,
    val warnings: List<String> = emptyList(),
    val systemDate: String? = null,
)

data class PortalPatientSummary(
    val patientId: String,
    val detailUrl: String? = null,
    val detailPayload: String? = null,
    val fullName: String,
    val birthDateOrYear: String? = null,
    val gender: String? = null,
    val phone: String? = null,
    val address: String? = null,
    val receivedDate: String? = null,
)

enum class RawSourceType {
    HTML,
    JSON,
    WEBVIEW_DOM,
    MOCK,
}

sealed class PortalLookupException(message: String) : Exception(message) {
    data object MissingSession : PortalLookupException("Chưa có phiên đăng nhập.")
    data object SessionExpired : PortalLookupException("Phiên đăng nhập đã hết hạn.")
    data object InvalidPhone : PortalLookupException("Số điện thoại không hợp lệ.")
    data object NotFound : PortalLookupException("Không tìm thấy dữ liệu tiêm chủng.")
    data object CaptchaRequired : PortalLookupException("VNCDC yêu cầu captcha hoặc OTP.")
    data object ParseFailed : PortalLookupException("Không đọc được dữ liệu từ VNCDC.")
    data class NetworkFailed(val detail: String) : PortalLookupException(detail)
}

interface PortalLookupRepository {
    suspend fun searchPatientsByPhone(phone: String): Result<List<PortalPatientSummary>>
    suspend fun lookupVaccinations(patient: PortalPatientSummary): Result<PortalLookupResult>

    suspend fun lookupByPhone(phone: String): Result<PortalLookupResult> =
        searchPatientsByPhone(phone).mapCatching { patients ->
            val patient = patients.firstOrNull() ?: throw PortalLookupException.NotFound
            lookupVaccinations(patient).getOrThrow()
        }
}

interface PortalParser {
    fun parse(raw: String, sourceType: RawSourceType): Result<PortalLookupResult>
    fun parsePatientSummaries(raw: String): Result<List<PortalPatientSummary>>
}
