package com.tracuutiemchung.app.data.portal

import kotlinx.serialization.Serializable

@Serializable
data class PortalSession(
    val cookieHeader: String,
    val cookies: Map<String, String>,
    val csrfToken: String? = null,
    val expiresAtMillis: Long? = null,
    val source: SessionSource,
)

@Serializable
enum class SessionSource {
    HTTP_CLIENT,
    WEBVIEW,
}
