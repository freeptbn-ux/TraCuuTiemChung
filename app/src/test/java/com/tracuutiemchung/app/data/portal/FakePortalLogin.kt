package com.tracuutiemchung.app.data.portal

fun fakeLogin(username: String, password: String): PortalSession {
    require(username.isNotBlank() && password.isNotBlank())
    return PortalSession(
        cookieHeader = "dev_session=local_login",
        cookies = mapOf("dev_session" to "local_login"),
        source = SessionSource.HTTP_CLIENT,
    )
}
