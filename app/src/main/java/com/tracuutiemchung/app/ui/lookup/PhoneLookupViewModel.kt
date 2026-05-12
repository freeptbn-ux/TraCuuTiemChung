package com.tracuutiemchung.app.ui.lookup

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.tracuutiemchung.app.data.model.AnalysisResult
import com.tracuutiemchung.app.data.portal.PortalLookupException
import com.tracuutiemchung.app.data.portal.PortalPatientSummary
import com.tracuutiemchung.app.domain.usecase.LookupVaccinationByPhoneUseCase
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

sealed interface LookupUiState {
    data object Idle : LookupUiState
    data object Loading : LookupUiState
    data class Success(val result: AnalysisResult) : LookupUiState
    data class PatientSelection(
        val phone: String,
        val patients: List<PortalPatientSummary>,
    ) : LookupUiState
    data class LoadingDetail(val patient: PortalPatientSummary) : LookupUiState
    data class Error(val message: String) : LookupUiState
    data object SessionExpired : LookupUiState
}

class PhoneLookupViewModel(
    private val lookupVaccinationByPhone: LookupVaccinationByPhoneUseCase,
    private val dispatcher: CoroutineDispatcher = Dispatchers.Main,
) : ViewModel() {
    private val _uiState = MutableStateFlow<LookupUiState>(LookupUiState.Idle)
    val uiState: StateFlow<LookupUiState> = _uiState.asStateFlow()

    fun validatePhone(phone: String): Boolean = Regex("^0\\d{9}$").matches(phone.trim())

    fun lookup(phone: String) = search(phone)

    fun search(phone: String) {
        val normalizedPhone = phone.trim()
        if (!validatePhone(normalizedPhone)) {
            _uiState.value = LookupUiState.Error("Số điện thoại phải gồm 10 chữ số và bắt đầu bằng 0.")
            return
        }

        viewModelScope.launch(dispatcher) {
            _uiState.value = LookupUiState.Loading
            val result = lookupVaccinationByPhone.searchPatients(normalizedPhone)
            _uiState.value = result.fold(
                onSuccess = { patients ->
                    when (patients.size) {
                        0 -> LookupUiState.Error("Không tìm thấy dữ liệu tiêm chủng.")
                        else -> LookupUiState.PatientSelection(normalizedPhone, patients)
                    }
                },
                onFailure = ::failureState,
            )
        }
    }

    fun selectPatient(patient: PortalPatientSummary) {
        viewModelScope.launch(dispatcher) {
            _uiState.value = LookupUiState.LoadingDetail(patient)
            _uiState.value = lookupVaccinationByPhone.lookupPatient(patient).fold(
                onSuccess = { LookupUiState.Success(it) },
                onFailure = ::failureState,
            )
        }
    }

    fun resetToSearch() {
        _uiState.value = LookupUiState.Idle
    }

    private fun failureState(error: Throwable): LookupUiState = when (error) {
        PortalLookupException.SessionExpired,
        PortalLookupException.MissingSession -> LookupUiState.SessionExpired
        else -> LookupUiState.Error(error.message ?: "Không thể tra cứu dữ liệu VNCDC.")
    }
}
