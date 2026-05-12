package com.tracuutiemchung.app.data.credentials

import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import java.security.KeyStore
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec
import kotlin.io.encoding.Base64
import kotlin.io.encoding.ExperimentalEncodingApi

interface CredentialCipher {
    fun encrypt(plainText: String): EncryptedPayload
    fun decrypt(payload: EncryptedPayload): String
}

class CredentialCrypto(
    private val keyAlias: String = KEY_ALIAS,
) : CredentialCipher {
    fun getOrCreateSecretKey(): SecretKey {
        val keyStore = KeyStore.getInstance(ANDROID_KEYSTORE).apply { load(null) }
        (keyStore.getEntry(keyAlias, null) as? KeyStore.SecretKeyEntry)?.let { return it.secretKey }

        val keyGenerator = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, ANDROID_KEYSTORE)
        val keySpec = KeyGenParameterSpec.Builder(
            keyAlias,
            KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT,
        )
            .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
            .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
            .setKeySize(AES_KEY_SIZE_BITS)
            .build()
        keyGenerator.init(keySpec)
        return keyGenerator.generateKey()
    }

    override fun encrypt(plainText: String): EncryptedPayload {
        val cipher = Cipher.getInstance(TRANSFORMATION)
        cipher.init(Cipher.ENCRYPT_MODE, getOrCreateSecretKey())
        val cipherText = cipher.doFinal(plainText.encodeToByteArray())
        return EncryptedPayload(
            ivBase64 = cipher.iv.base64Encode(),
            cipherTextBase64 = cipherText.base64Encode(),
        )
    }

    override fun decrypt(payload: EncryptedPayload): String {
        val cipher = Cipher.getInstance(TRANSFORMATION)
        cipher.init(
            Cipher.DECRYPT_MODE,
            getOrCreateSecretKey(),
            GCMParameterSpec(GCM_TAG_LENGTH_BITS, payload.ivBase64.base64Decode()),
        )
        return cipher.doFinal(payload.cipherTextBase64.base64Decode()).decodeToString()
    }

    companion object {
        const val KEY_ALIAS = "tracuutiemchung_vncdc_credentials_aes_gcm"
        private const val ANDROID_KEYSTORE = "AndroidKeyStore"
        private const val TRANSFORMATION = "AES/GCM/NoPadding"
        private const val AES_KEY_SIZE_BITS = 256
        private const val GCM_TAG_LENGTH_BITS = 128
    }
}

@OptIn(ExperimentalEncodingApi::class)
private fun ByteArray.base64Encode(): String = Base64.encode(this)

@OptIn(ExperimentalEncodingApi::class)
private fun String.base64Decode(): ByteArray = Base64.decode(this)
