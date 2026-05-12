package com.tracuutiemchung.app.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Typography
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val AppColorScheme = lightColorScheme(
    primary = Color(0xFF006A60),
    onPrimary = Color.White,
    secondary = Color(0xFF4A635F),
    tertiary = Color(0xFF456179),
    background = Color(0xFFF5FBF8),
    surface = Color(0xFFFFFFFF),
    surfaceVariant = Color(0xFFDBE5E1),
    error = Color(0xFFBA1A1A),
)

private val AppTypography = Typography()

@Composable
fun TraCuuTiemChungTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        colorScheme = AppColorScheme,
        typography = AppTypography,
        content = content,
    )
}
