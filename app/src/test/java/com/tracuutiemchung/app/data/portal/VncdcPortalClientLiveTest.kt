package com.tracuutiemchung.app.data.portal

import kotlinx.coroutines.test.runTest
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Assume.assumeTrue
import org.junit.Test

class VncdcPortalClientLiveTest {
    @Test
    fun loginWithProvidedCredentialsReturnsAuthenticatedSession() = runTest {
        assumeTrue("Live test disabled. Set RUN_VNCDC_LIVE_TEST=true to run.", System.getenv("RUN_VNCDC_LIVE_TEST") == "true")
        val username = System.getenv("VNCDC_USERNAME").orEmpty()
        val password = System.getenv("VNCDC_PASSWORD").orEmpty()
        assertFalse("VNCDC_USERNAME is required", username.isBlank())
        assertFalse("VNCDC_PASSWORD is required", password.isBlank())

        val session = VncdcPortalClient().login(username, password)

        assertEquals(SessionSource.HTTP_CLIENT, session.source)
        assertTrue("Expected non-empty cookie header", session.cookieHeader.isNotBlank())
        assertTrue(
            "Expected authenticated VNCDC cookie or a reusable cookie set",
            session.cookies.containsKey(".ASPXAUTH") || session.cookies.isNotEmpty(),
        )
    }
}
