package com.tracuutiemchung.app.ui

import com.tracuutiemchung.app.data.model.AnalysisResult

sealed interface AppUiState {
    data object NotLoggedIn : AppUiState
    data object LoggingIn : AppUiState
    data object ReadyForLookup : AppUiState
    data object LookingUp : AppUiState
    data class ShowingResult(val result: AnalysisResult) : AppUiState
    data class Error(val message: String) : AppUiState
}
