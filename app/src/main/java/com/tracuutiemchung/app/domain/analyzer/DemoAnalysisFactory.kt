package com.tracuutiemchung.app.domain.analyzer

import com.tracuutiemchung.app.data.model.AnalysisResult
import com.tracuutiemchung.app.data.model.PatientInfo
import com.tracuutiemchung.app.data.model.VaccinationRecord
import com.tracuutiemchung.app.data.model.VaccineRecommendation

object DemoAnalysisFactory {
    fun create(phoneNumber: String): AnalysisResult {
        val normalizedPhone = phoneNumber.trim().ifBlank { "Chưa nhập" }
        return AnalysisResult(
            patientInfo = PatientInfo(
                fullName = "Nguyễn Văn A",
                phoneNumber = normalizedPhone,
                dateOfBirth = "01/01/2020",
                address = "Hồ Chí Minh",
            ),
            records = listOf(
                VaccinationRecord(
                    vaccineName = "Sởi - Quai bị - Rubella",
                    doseNumber = 1,
                    vaccinationDate = "15/03/2021",
                    provider = "VNCDC",
                ),
                VaccinationRecord(
                    vaccineName = "Viêm gan B",
                    doseNumber = 3,
                    vaccinationDate = "20/08/2021",
                    provider = "VNCDC",
                ),
            ),
            recommendations = listOf(
                VaccineRecommendation(
                    vaccineName = "Sởi - Quai bị - Rubella",
                    reason = "Demo mũi nhắc theo lịch tiêm chủng",
                    recommendedDate = "15/03/2026",
                ),
            ),
            warnings = listOf("Dữ liệu demo, chưa kết nối cổng thật ở phase này."),
        )
    }
}
