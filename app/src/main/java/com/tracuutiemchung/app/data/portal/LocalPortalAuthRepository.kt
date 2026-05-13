package com.tracuutiemchung.app.data.portal

class LoginInputException(message: String) : IllegalArgumentException(message)

class LocalPortalAuthRepository(
    private val sessionStore: SessionStore,
    private val loginProvider: suspend (username: String, password: String) -> PortalSession = { username, password ->
        VncdcPortalClient().login(username, password)
    },
) : PortalAuthRepository {
    override suspend fun login(username: String, password: String): Result<PortalSession> {
        if (username.isBlank() || password.isBlank()) {
            return Result.failure(LoginInputException("Vui lòng nhập tài khoản và mật khẩu."))
        }

        return runCatching { loginProvider(username.trim(), password) }
            .onSuccess(sessionStore::set)
            .recoverCatching { error ->
                if (error is LoginInputException) throw error
                throw LoginInputException(error.message ?: "Không thể đăng nhập VNCDC.")
            }
    }

    override suspend fun logout() {
        sessionStore.clear()
    }

    override suspend fun currentSession(): PortalSession? = sessionStore.get()

    override suspend fun refreshIfNeeded(): Result<PortalSession?> = Result.success(sessionStore.get())
}
