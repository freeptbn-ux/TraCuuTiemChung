package com.tracuutiemchung.app.data.credentials

import kotlinx.serialization.Serializable

/** Credentials are kept in this model only after a successful decrypt. */
data class SavedCredentials(
    val username: String,
    val password: String,
)

@Serializable
data class EncryptedPayload(
    val ivBase64: String,
    val cipherTextBase64: String,
)
