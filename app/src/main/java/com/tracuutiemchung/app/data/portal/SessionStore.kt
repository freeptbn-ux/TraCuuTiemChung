package com.tracuutiemchung.app.data.portal

interface SessionStore {
    fun get(): PortalSession?
    fun set(session: PortalSession)
    fun clear()
}

class InMemorySessionStore : SessionStore {
    private var session: PortalSession? = null

    override fun get(): PortalSession? = session

    override fun set(session: PortalSession) {
        this.session = session
    }

    override fun clear() {
        session = null
    }
}
