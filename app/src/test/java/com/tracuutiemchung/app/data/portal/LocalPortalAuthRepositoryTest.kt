package com.tracuutiemchung.app.data.portal

import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertNotNull
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Test

class LocalPortalAuthRepositoryTest {
    @Test
    fun loginWithBlankInputFailsAndDoesNotStoreSession() = runTest {
        val store = InMemorySessionStore()
        val repository = LocalPortalAuthRepository(store, ::fakeLogin)

        val result = repository.login("", "secret")

        assertTrue(result.isFailure)
        assertNull(store.get())
    }

    @Test
    fun loginWithCredentialsStoresSessionForReuse() = runTest {
        val store = InMemorySessionStore()
        val repository = LocalPortalAuthRepository(store, ::fakeLogin)

        val result = repository.login("dev-user", "dev-pass")

        assertTrue(result.isSuccess)
        val session = result.getOrThrow()
        assertEquals(SessionSource.HTTP_CLIENT, session.source)
        assertEquals("dev_session=local_login", session.cookieHeader)
        assertEquals(session, repository.currentSession())
        assertNotNull(store.get())
    }

    @Test
    fun logoutClearsSession() = runTest {
        val store = InMemorySessionStore()
        val repository = LocalPortalAuthRepository(store, ::fakeLogin)
        repository.login("dev-user", "dev-pass")

        repository.logout()

        assertNull(repository.currentSession())
    }
}
