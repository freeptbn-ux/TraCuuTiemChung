package com.tracuutiemchung.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import com.tracuutiemchung.app.data.model.AnalysisResult
import com.tracuutiemchung.app.data.remote.VercelPortalRepository
import com.tracuutiemchung.app.domain.usecase.LookupVaccinationByPhoneUseCase
import com.tracuutiemchung.app.ui.lookup.PhoneLookupScreen
import com.tracuutiemchung.app.ui.lookup.PhoneLookupViewModel
import com.tracuutiemchung.app.ui.result.ResultScreen
import com.tracuutiemchung.app.ui.theme.TraCuuTiemChungTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            TraCuuTiemChungTheme {
                TraCuuTiemChungApp()
            }
        }
    }
}

@Composable
fun TraCuuTiemChungApp() {
    var screen by remember { mutableStateOf(AppScreen.Lookup) }
    var result by remember { mutableStateOf<AnalysisResult?>(null) }
    
    val lookupViewModel = remember {
        PhoneLookupViewModel(
            LookupVaccinationByPhoneUseCase(VercelPortalRepository()),
        )
    }

    when (screen) {
        AppScreen.Lookup -> PhoneLookupScreen(
            viewModel = lookupViewModel,
            onLookupSuccess = { analysisResult ->
                result = analysisResult
                screen = AppScreen.Result
            },
            onSessionExpired = { /* Not applicable for Vercel */ },
        )
        AppScreen.Result -> ResultScreen(requireNotNull(result))
    }
}

enum class AppScreen {
    Lookup,
    Result,
}
