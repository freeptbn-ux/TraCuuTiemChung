package com.tracuutiemchung.app.data.portal

interface PortalAuthRepository {
    suspend fun login(username: String, password: String): Result<PortalSession>
    suspend fun logout()
    suspend fun currentSession(): PortalSession?
    suspend fun refreshIfNeeded(): Result<PortalSession?>
}
