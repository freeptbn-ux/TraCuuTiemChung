package com.tracuutiemchung.app.ui.theme

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp

@Composable
fun ThemeVerificationScreen() {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(MaterialTheme.colorScheme.background)
            .verticalScroll(rememberScrollState())
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        Text(
            text = "Design System Verification",
            style = MaterialTheme.typography.headlineLarge,
            color = MaterialTheme.colorScheme.primary
        )

        Divider()

        Text(text = "Typography", style = MaterialTheme.typography.titleLarge)
        
        Text(text = "Display Large", style = MaterialTheme.typography.displayLarge)
        Text(text = "Headline Medium", style = MaterialTheme.typography.headlineMedium)
        Text(text = "Title Large", style = MaterialTheme.typography.titleLarge)
        Text(text = "Body Large", style = MaterialTheme.typography.bodyLarge)
        Text(text = "Label Small", style = MaterialTheme.typography.labelSmall)

        Divider()

        Text(text = "Colors", style = MaterialTheme.typography.titleLarge)
        
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            ColorBox("Primary", MaterialTheme.colorScheme.primary, MaterialTheme.colorScheme.onPrimary)
            ColorBox("Secondary", MaterialTheme.colorScheme.secondary, MaterialTheme.colorScheme.onSecondary)
            ColorBox("Tertiary", MaterialTheme.colorScheme.tertiary, MaterialTheme.colorScheme.onTertiary)
        }
        
        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
            ColorBox("Surface", MaterialTheme.colorScheme.surface, MaterialTheme.colorScheme.onSurface)
            ColorBox("Error", MaterialTheme.colorScheme.error, MaterialTheme.colorScheme.onError)
        }

        Divider()

        Text(text = "Shapes & Elevation", style = MaterialTheme.typography.titleLarge)
        
        Button(
            onClick = {},
            shape = MaterialTheme.shapes.medium,
            modifier = Modifier.fillMaxWidth()
        ) {
            Text("Medium Shape Button")
        }

        Card(
            shape = MaterialTheme.shapes.large,
            modifier = Modifier.fillMaxWidth().height(100.dp),
            colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
        ) {
            Box(contentAlignment = Alignment.Center, modifier = Modifier.fillMaxSize()) {
                Text("Large Shape Card", style = MaterialTheme.typography.bodyLarge)
            }
        }
    }
}

@Composable
fun ColorBox(name: String, color: Color, onColor: Color) {
    Box(
        modifier = Modifier
            .size(100.dp)
            .background(color, shape = MaterialTheme.shapes.small)
            .padding(8.dp),
        contentAlignment = Alignment.Center
    ) {
        Text(text = name, color = onColor, style = MaterialTheme.typography.labelSmall)
    }
}

@Preview(showBackground = true, name = "Light Theme")
@Composable
fun LightThemePreview() {
    TraCuuTiemChungTheme(darkTheme = false) {
        ThemeVerificationScreen()
    }
}

@Preview(showBackground = true, name = "Dark Theme")
@Composable
fun DarkThemePreview() {
    TraCuuTiemChungTheme(darkTheme = true) {
        ThemeVerificationScreen()
    }
}
