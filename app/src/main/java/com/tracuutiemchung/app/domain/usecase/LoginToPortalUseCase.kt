package com.tracuutiemchung.app.domain.usecase

import com.tracuutiemchung.app.data.portal.PortalAuthRepository
import com.tracuutiemchung.app.data.portal.PortalSession

class LoginToPortalUseCase(
    private val repository: PortalAuthRepository,
) {
    suspend operator fun invoke(username: String, password: String): Result<PortalSession> {
        return repository.login(username.trim(), password)
    }
}
