package com.tracuutiemchung.app.data.credentials

import com.tracuutiemchung.app.data.security.SecureStorage
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first

interface CredentialStore {
    val savedCredentials: Flow<SavedCredentials?>
    suspend fun save(username: String, password: String)
    suspend fun clear()
}

class SecureCredentialStore(
    private val secureStorage: SecureStorage
) : CredentialStore {
    private val _credentials = MutableStateFlow<SavedCredentials?>(loadFromStorage())
    override val savedCredentials: Flow<SavedCredentials?> = _credentials.asStateFlow()

    override suspend fun save(username: String, password: String) {
        secureStorage.saveString(KEY_USERNAME, username)
        secureStorage.saveString(KEY_PASSWORD, password)
        _credentials.value = SavedCredentials(username, password)
    }

    override suspend fun clear() {
        secureStorage.clear()
        _credentials.value = null
    }

    private fun loadFromStorage(): SavedCredentials? {
        val username = secureStorage.getString(KEY_USERNAME)
        val password = secureStorage.getString(KEY_PASSWORD)
        return if (!username.isNullOrBlank() && !password.isNullOrBlank()) {
            SavedCredentials(username, password)
        } else {
            null
        }
    }

    companion object {
        private const val KEY_USERNAME = "vncdc_username"
        private const val KEY_PASSWORD = "vncdc_password"
    }
}

suspend fun CredentialStore.loadOnce(): SavedCredentials? = savedCredentials.first()
