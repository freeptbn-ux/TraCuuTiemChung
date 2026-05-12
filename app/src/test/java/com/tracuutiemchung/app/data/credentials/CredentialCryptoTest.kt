package com.tracuutiemchung.app.data.credentials

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertNotEquals
import org.junit.Test

class CredentialCryptoTest {
    @Test
    fun encryptDecryptRoundTripKeepsPlainPasswordOutOfPayload() {
        val crypto = JvmCredentialCrypto()
        val password = "secret-pass-123"

        val payload = crypto.encrypt(password)
        val decrypted = crypto.decrypt(payload)

        assertEquals(password, decrypted)
        assertFalse(payload.ivBase64.contains(password))
        assertFalse(payload.cipherTextBase64.contains(password))
    }

    @Test
    fun encryptUsesDifferentIvForSamePlainText() {
        val crypto = JvmCredentialCrypto()
        val password = "same-password"

        val first = crypto.encrypt(password)
        val second = crypto.encrypt(password)

        assertNotEquals(first.ivBase64, second.ivBase64)
        assertNotEquals(first.cipherTextBase64, second.cipherTextBase64)
    }
}
