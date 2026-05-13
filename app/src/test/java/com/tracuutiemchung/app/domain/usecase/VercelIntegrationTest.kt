package com.tracuutiemchung.app.domain.usecase

import com.jakewharton.retrofit2.converter.kotlinx.serialization.asConverterFactory
import com.tracuutiemchung.app.data.model.RecommendationStatus
import com.tracuutiemchung.app.data.portal.PortalPatientSummary
import com.tracuutiemchung.app.data.remote.VercelApiService
import com.tracuutiemchung.app.data.remote.VercelPortalRepository
import kotlinx.coroutines.runBlocking
import kotlinx.serialization.json.Json
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.mockwebserver.MockResponse
import okhttp3.mockwebserver.MockWebServer
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import retrofit2.Retrofit

class VercelIntegrationTest {

    private lateinit var mockWebServer: MockWebServer
    private lateinit var apiService: VercelApiService
    private lateinit var repository: VercelPortalRepository
    private lateinit var useCase: LookupVaccinationByPhoneUseCase

    private val json = Json { ignoreUnknownKeys = true }

    @Before
    fun setup() {
        mockWebServer = MockWebServer()
        mockWebServer.start()

        val contentType = "application/json".toMediaType()
        apiService = Retrofit.Builder()
            .baseUrl(mockWebServer.url("/"))
            .addConverterFactory(json.asConverterFactory(contentType))
            .build()
            .create(VercelApiService::class.java)

        repository = VercelPortalRepository(apiService)
        useCase = LookupVaccinationByPhoneUseCase(repository)
    }

    @After
    fun teardown() {
        mockWebServer.shutdown()
    }

    @Test
    fun `searchPatients returns list of patients from Vercel`() = runBlocking {
        val mockResponseBody = """
            {
                "status": "success",
                "data": [
                    {
                        "id": "P123",
                        "name": "Nguyen Van A",
                        "dob": "01/01/2020",
                        "gender": "Nam"
                    }
                ]
            }
        """.trimIndent()

        mockWebServer.enqueue(MockResponse().setBody(mockResponseBody).setResponseCode(200))

        val result = useCase.searchPatients("0123456789")

        assertTrue(result.isSuccess)
        val patients = result.getOrThrow()
        assertEquals(1, patients.size)
        assertEquals("Nguyen Van A", patients[0].fullName)
        assertEquals("P123", patients[0].patientId)
    }

    @Test
    fun `lookupPatient returns analysis results from Vercel`() = runBlocking {
        val mockResponseBody = """
            {
                "status": "success",
                "data": {
                    "patient_name": "Nguyen Van A",
                    "dob": "01/01/2020",
                    "analysis_date": "12/05/2026",
                    "missing_vaccines": [
                        {
                            "vaccine_name": "Lao (BCG)",
                            "rule_type": "BCG",
                            "status": "DUE_NOW",
                            "next_dose": "12/05/2026",
                            "message": "Cần tiêm ngay",
                            "status_tags": ["TIÊM_NGAY"]
                        }
                    ]
                }
            }
        """.trimIndent()

        mockWebServer.enqueue(MockResponse().setBody(mockResponseBody).setResponseCode(200))

        val patient = PortalPatientSummary(patientId = "P123", fullName = "Nguyen Van A", phone = "0123456789")
        val result = useCase.lookupPatient(patient)

        assertTrue(result.isSuccess)
        val analysis = result.getOrThrow()
        assertEquals("Nguyen Van A", analysis.patientInfo.fullName)
        assertEquals(1, analysis.recommendations.size)
        assertEquals("Lao (BCG)", analysis.recommendations[0].vaccineName)
        assertEquals(RecommendationStatus.DUE_NOW, analysis.recommendations[0].status)
    }
}
