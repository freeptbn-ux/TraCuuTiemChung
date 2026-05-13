package com.tracuutiemchung.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.BackHandler
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.toRoute
import com.tracuutiemchung.app.data.remote.VercelPortalRepository
import com.tracuutiemchung.app.domain.usecase.LookupVaccinationByPhoneUseCase
import com.tracuutiemchung.app.ui.lookup.LookupUiState
import com.tracuutiemchung.app.ui.lookup.PhoneLookupScreen
import com.tracuutiemchung.app.ui.lookup.PhoneLookupViewModel
import com.tracuutiemchung.app.ui.navigation.AppRoute
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
    val navController = rememberNavController()
    
    // Shared ViewModel for both screens
    val lookupViewModel: PhoneLookupViewModel = viewModel {
        PhoneLookupViewModel(
            LookupVaccinationByPhoneUseCase(VercelPortalRepository()),
        )
    }

    NavHost(
        navController = navController,
        startDestination = AppRoute.Lookup
    ) {
        composable<AppRoute.Lookup> {
            PhoneLookupScreen(
                viewModel = lookupViewModel,
                onNavigateToSelection = { phone ->
                    navController.navigate(AppRoute.Selection(phone))
                },
                onNavigateToResult = { phone ->
                    navController.navigate(AppRoute.Result(phone))
                },
                onSessionExpired = { /* Not applicable for Vercel */ },
            )
        }
        composable<AppRoute.Selection> { backStackEntry ->
            val route = backStackEntry.toRoute<AppRoute.Selection>()
            val state by lookupViewModel.uiState.collectAsState()
            val patients by lookupViewModel.patients.collectAsState()
            
            // Auto-navigate to Result if only one patient is found, 
            // but ONLY if we are in the PatientSelection state (not Success/LoadingDetail)
            // to avoid infinite loops when coming back from Result screen.
            LaunchedEffect(state) {
                if (state is LookupUiState.PatientSelection && (state as LookupUiState.PatientSelection).patients.size == 1) {
                    val patient = (state as LookupUiState.PatientSelection).patients.first()
                    lookupViewModel.selectPatient(patient)
                    navController.navigate(AppRoute.Result(route.phoneNumber))
                }
            }

            BackHandler {
                lookupViewModel.resetToSearch()
                navController.popBackStack()
            }

            if (patients.isNotEmpty()) {
                com.tracuutiemchung.app.ui.lookup.PatientSelectionScreen(
                    phone = route.phoneNumber,
                    patients = patients,
                    onSelectPatient = { patient ->
                        lookupViewModel.selectPatient(patient)
                        navController.navigate(AppRoute.Result(route.phoneNumber))
                    },
                    onEditPhone = {
                        lookupViewModel.resetToSearch()
                        navController.popBackStack()
                    }
                )
            } else {
                LaunchedEffect(Unit) {
                    if (state !is LookupUiState.Loading && state !is LookupUiState.PatientSelection) {
                        navController.popBackStack()
                    }
                }
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator()
                }
            }
        }
        composable<AppRoute.Result> {
            val result by lookupViewModel.lastResult.collectAsState()
            val state by lookupViewModel.uiState.collectAsState()
            
            result?.let {
                ResultScreen(
                    result = it,
                    onBackClick = { navController.popBackStack() }
                )
            } ?: run {
                if (state is LookupUiState.LoadingDetail) {
                    Box(
                        modifier = Modifier.fillMaxSize(),
                        contentAlignment = Alignment.Center
                    ) {
                        CircularProgressIndicator()
                    }
                } else {
                    // Avoid popping back too early if we just started
                    LaunchedEffect(state) {
                        if (state !is LookupUiState.LoadingDetail && state !is LookupUiState.Success) {
                            navController.popBackStack()
                        }
                    }
                }
            }
        }
    }
}

