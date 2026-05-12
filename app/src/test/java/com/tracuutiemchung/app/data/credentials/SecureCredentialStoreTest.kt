package com.tracuutiemchung.app.data.credentials

import androidx.datastore.preferences.core.PreferenceDataStoreFactory
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.Job
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.test.TestScope
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.runTest
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNull
import org.junit.Before
import org.junit.Test
import java.io.File

@OptIn(ExperimentalCoroutinesApi::class)
class SecureCredentialStoreTest {
    private lateinit var tempDir: File

    @Before
    fun setUp() {
        tempDir = createTempDir(prefix = "credential-store-test")
    }

    @After
    fun tearDown() {
        tempDir.deleteRecursively()
    }

    @Test
    fun saveLoadClearCredentialsWithEncryptedFilePayload() = runTest {
        val dataStoreFile = File(tempDir, "vncdc_credentials.preferences_pb")
        val dataStore = testDataStore(dataStoreFile)
        val store = SecureCredentialStore(dataStore, JvmCredentialCrypto())
        val username = "vncdc-user"
        val password = "plain-password-456"

        store.save(username, password)

        val saved = store.savedCredentials.first()
        assertEquals(SavedCredentials(username, password), saved)
        assertFalse(dataStoreFile.readText(Charsets.ISO_8859_1).contains(password))

        store.clear()

        assertNull(store.savedCredentials.first())
    }

    @Test
    fun corruptPayloadClearsCredentialsAndReturnsNull() = runTest {
        val dataStore = testDataStore(File(tempDir, "corrupt.preferences_pb"))
        val store = SecureCredentialStore(dataStore, JvmCredentialCrypto())

        dataStore.edit { preferences ->
            preferences[SecureCredentialStore.USERNAME_KEY] = "vncdc-user"
            preferences[SecureCredentialStore.PASSWORD_PAYLOAD_KEY] = "not-json"
        }

        assertNull(store.savedCredentials.first())
        assertNull(dataStore.data.first()[SecureCredentialStore.USERNAME_KEY])
        assertNull(dataStore.data.first()[SecureCredentialStore.PASSWORD_PAYLOAD_KEY])
    }

    private fun testDataStore(file: File): androidx.datastore.core.DataStore<Preferences> =
        PreferenceDataStoreFactory.create(
            scope = TestScope(UnconfinedTestDispatcher() + Job()),
            produceFile = { file },
        )
}
