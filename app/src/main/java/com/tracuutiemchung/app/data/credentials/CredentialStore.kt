package com.tracuutiemchung.app.data.credentials

import android.content.Context
import androidx.datastore.core.CorruptionException
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.emptyPreferences
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.launch
import kotlinx.serialization.SerializationException
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import java.io.IOException

interface CredentialStore {
    val savedCredentials: Flow<SavedCredentials?>
    suspend fun save(username: String, password: String)
    suspend fun clear()
}

class SecureCredentialStore(
    private val dataStore: DataStore<Preferences>,
    private val crypto: CredentialCipher,
    private val scope: CoroutineScope? = null,
) : CredentialStore {
    override val savedCredentials: Flow<SavedCredentials?> = dataStore.data
        .catch { exception ->
            if (exception is IOException || exception is CorruptionException) {
                emit(emptyPreferences())
            } else {
                throw exception
            }
        }
        .map { preferences ->
            val username = preferences[USERNAME_KEY]
            val payloadJson = preferences[PASSWORD_PAYLOAD_KEY]
            if (username.isNullOrBlank() || payloadJson.isNullOrBlank()) {
                return@map null
            }

            runCatching {
                val payload = json.decodeFromString<EncryptedPayload>(payloadJson)
                SavedCredentials(username = username, password = crypto.decrypt(payload))
            }.getOrElse {
                clearAfterCorruption()
                null
            }
        }

    override suspend fun save(username: String, password: String) {
        val encryptedPassword = crypto.encrypt(password)
        dataStore.edit { preferences ->
            preferences[USERNAME_KEY] = username
            preferences[PASSWORD_PAYLOAD_KEY] = json.encodeToString(encryptedPassword)
        }
    }

    override suspend fun clear() {
        dataStore.edit { preferences -> preferences.clear() }
    }

    private suspend fun clearAfterCorruption() {
        if (scope == null) {
            clear()
        } else {
            scope.launch { clear() }
        }
    }

    companion object {
        internal val USERNAME_KEY = stringPreferencesKey("vncdc_username")
        internal val PASSWORD_PAYLOAD_KEY = stringPreferencesKey("vncdc_password_payload")
        private val json = Json { ignoreUnknownKeys = true }
    }
}

val Context.vncdcCredentialDataStore by preferencesDataStore(name = "vncdc_credentials")

suspend fun CredentialStore.loadOnce(): SavedCredentials? = savedCredentials.first()
