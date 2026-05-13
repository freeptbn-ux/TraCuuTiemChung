package com.tracuutiemchung.app.ui.theme

import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.sp
import org.junit.Assert.assertEquals
import org.junit.Test

class ThemeConfigTest {

    @Test
    fun testLightPrimaryColor() {
        // Verify that the light primary color is correctly defined
        assertEquals(Color(0xFF006A60), md_theme_light_primary)
    }

    @Test
    fun testDarkPrimaryColor() {
        // Verify that the dark primary color is correctly defined
        assertEquals(Color(0xFF53DBC9), md_theme_dark_primary)
    }

    @Test
    fun testTypographyStyles() {
        // Verify that headlineLarge typography has the correct font size
        // Note: We can't easily check Typography object in a pure unit test without compose-ui
        // but we can check if it's initialized
        assert(Typography.headlineLarge.fontSize == 32.sp)
    }
}
