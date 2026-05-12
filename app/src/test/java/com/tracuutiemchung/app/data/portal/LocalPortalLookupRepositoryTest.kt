package com.tracuutiemchung.app.data.portal

import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertSame
import org.junit.Assert.assertTrue
import org.junit.Test

class LocalPortalLookupRepositoryTest {
    @Test
    fun invalidPhoneFailsBeforeSessionCheck() = runTest {
        val repository = LocalPortalLookupRepository(InMemorySessionStore())

        val result = repository.lookupByPhone("123")

        assertTrue(result.isFailure)
        assertTrue(result.exceptionOrNull() is PortalLookupException.InvalidPhone)
    }

    @Test
    fun missingSessionReturnsSessionError() = runTest {
        val repository = LocalPortalLookupRepository(InMemorySessionStore())

        val result = repository.lookupByPhone("0912345678")

        assertTrue(result.isFailure)
        assertTrue(result.exceptionOrNull() is PortalLookupException.MissingSession)
    }

    @Test
    fun expiredSessionReturnsSessionExpired() = runTest {
        val store = InMemorySessionStore().apply {
            set(testSession(expiresAtMillis = 1L))
        }
        val repository = LocalPortalLookupRepository(store)

        val result = repository.lookupByPhone("0912345678")

        assertTrue(result.isFailure)
        assertTrue(result.exceptionOrNull() is PortalLookupException.SessionExpired)
    }

    @Test
    fun validSessionParsesProvidedResponse() = runTest {
        val store = InMemorySessionStore().apply { set(testSession()) }
        lateinit var forwardedSession: PortalSession
        var forwardedPhone = ""
        val repository = LocalPortalLookupRepository(
            sessionStore = store,
            responseProvider = { phone, session ->
                forwardedPhone = phone
                forwardedSession = session
                """
                    {
                      "patient": {"name":"Phạm Bé D","phone":"$phone"},
                      "vaccinations": [
                        {"vaccineName":"BCG","doseText":"Mũi 1","dateText":"05/05/2022","facilityName":"VNCDC"}
                      ]
                    }
                """.trimIndent() to RawSourceType.JSON
            },
        )

        val result = repository.lookupByPhone("0912345678").getOrThrow()

        assertEquals("0912345678", forwardedPhone)
        assertSame(store.get(), forwardedSession)
        assertEquals("SESSION=abc", forwardedSession.cookieHeader)
        assertEquals("Phạm Bé D", result.patient?.fullName)
        assertEquals("0912345678", result.patient?.phoneNumber)
        assertEquals("BCG", result.vaccinations.first().vaccineName)
        assertTrue(result.warnings.contains("Thiếu ngày sinh, không tự suy đoán."))
    }

    @Test
    fun responseProviderFailureReturnsNetworkFailedMessage() = runTest {
        val store = InMemorySessionStore().apply { set(testSession()) }
        val repository = LocalPortalLookupRepository(
            sessionStore = store,
            responseProvider = { _, _ -> error("timeout") },
        )

        val result = repository.lookupByPhone("0912345678")

        assertTrue(result.isFailure)
        assertEquals("timeout", result.exceptionOrNull()?.message)
        assertTrue(result.exceptionOrNull() is PortalLookupException.NetworkFailed)
    }

    private fun testSession(expiresAtMillis: Long? = null): PortalSession = PortalSession(
        cookieHeader = "SESSION=abc",
        cookies = mapOf("SESSION" to "abc"),
        csrfToken = "csrf",
        expiresAtMillis = expiresAtMillis,
        source = SessionSource.HTTP_CLIENT,
    )
}
