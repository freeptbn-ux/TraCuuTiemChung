package com.tracuutiemchung.app.data.remote

import com.tracuutiemchung.app.data.remote.model.*
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.Header
import retrofit2.http.POST

interface VercelApiService {

    @POST("api/lookup")
    suspend fun lookupPatients(
        @Body request: Map<String, String>
    ): Response<VercelResponse<List<PatientLookupDto>>>

    @POST("api/analyze")
    suspend fun analyzeVaccinations(
        @Body request: AnalysisRequestDto
    ): Response<VercelResponse<AnalysisResponseDto>>
}
