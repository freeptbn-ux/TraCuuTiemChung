package com.tracuutiemchung.app.ui.navigation

import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import org.junit.Assert.assertEquals
import org.junit.Test

class NavigationTest {

    @Test
    fun `test AppRoute serialization`() {
        val lookup = AppRoute.Lookup
        val selection = AppRoute.Selection("0987654321")
        val result = AppRoute.Result("0123456789")

        val lookupJson = Json.encodeToString<AppRoute>(lookup)
        val selectionJson = Json.encodeToString<AppRoute>(selection)
        val resultJson = Json.encodeToString<AppRoute>(result)

        assertEquals("{\"type\":\"com.tracuutiemchung.app.ui.navigation.AppRoute.Lookup\"}", lookupJson)
        assertEquals("{\"type\":\"com.tracuutiemchung.app.ui.navigation.AppRoute.Selection\",\"phoneNumber\":\"0987654321\"}", selectionJson)
        assertEquals("{\"type\":\"com.tracuutiemchung.app.ui.navigation.AppRoute.Result\",\"phoneNumber\":\"0123456789\"}", resultJson)
    }

    @Test
    fun `test AppRoute deserialization`() {
        val selectionJson = "{\"type\":\"com.tracuutiemchung.app.ui.navigation.AppRoute.Selection\",\"phoneNumber\":\"0987654321\"}"
        val selection = Json.decodeFromString<AppRoute>(selectionJson)

        assert(selection is AppRoute.Selection)
        assertEquals("0987654321", (selection as AppRoute.Selection).phoneNumber)
    }
}
