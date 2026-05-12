package com.tracuutiemchung.app.data.credentials

import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec

class JvmCredentialCrypto(
    private val secretKey: SecretKey = KeyGenerator.getInstance("AES").apply { init(256) }.generateKey(),
) : CredentialCipher {
    override fun encrypt(plainText: String): EncryptedPayload {
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        cipher.init(Cipher.ENCRYPT_MODE, secretKey)
        return EncryptedPayload(
            ivBase64 = java.util.Base64.getEncoder().encodeToString(cipher.iv),
            cipherTextBase64 = java.util.Base64.getEncoder().encodeToString(cipher.doFinal(plainText.encodeToByteArray())),
        )
    }

    override fun decrypt(payload: EncryptedPayload): String {
        val cipher = Cipher.getInstance("AES/GCM/NoPadding")
        cipher.init(
            Cipher.DECRYPT_MODE,
            secretKey,
            GCMParameterSpec(128, java.util.Base64.getDecoder().decode(payload.ivBase64)),
        )
        return cipher.doFinal(java.util.Base64.getDecoder().decode(payload.cipherTextBase64)).decodeToString()
    }
}
