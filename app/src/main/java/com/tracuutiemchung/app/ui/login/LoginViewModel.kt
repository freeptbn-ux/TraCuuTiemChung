package com.tracuutiemchung.app.ui.login

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.tracuutiemchung.app.data.credentials.CredentialStore
import com.tracuutiemchung.app.data.credentials.SavedCredentials
import com.tracuutiemchung.app.data.portal.PortalSession
import com.tracuutiemchung.app.domain.usecase.LoginToPortalUseCase
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch

sealed interface LoginUiState {
    data object Idle : LoginUiState
    data object Loading : LoginUiState
    data class Warning(val message: String) : LoginUiState
    data class Success(val session: PortalSession) : LoginUiState
    data class Error(val message: String) : LoginUiState
}

data class SavedCredentialUiState(
    val isLoading: Boolean = false,
    val credentials: SavedCredentials? = null,
    val warningMessage: String? = null,
) {
    val hasSavedCredentials: Boolean = credentials != null
}

class LoginViewModel(
    private val loginToPortal: LoginToPortalUseCase,
    private val credentialStore: CredentialStore? = null,
) : ViewModel() {
    private val _uiState = MutableStateFlow<LoginUiState>(LoginUiState.Idle)
    val uiState: StateFlow<LoginUiState> = _uiState.asStateFlow()

    private val _savedCredentialUiState = MutableStateFlow(SavedCredentialUiState(isLoading = credentialStore != null))
    val savedCredentialUiState: StateFlow<SavedCredentialUiState> = _savedCredentialUiState.asStateFlow()

    init {
        loadSavedCredentials()
    }

    fun loadSavedCredentials() {
        val store = credentialStore ?: return
        viewModelScope.launch {
            _savedCredentialUiState.value = _savedCredentialUiState.value.copy(isLoading = true, warningMessage = null)
            runCatching { store.savedCredentials.first() }
                .onSuccess { credentials ->
                    _savedCredentialUiState.value = SavedCredentialUiState(
                        isLoading = false,
                        credentials = credentials,
                    )
                }
                .onFailure {
                    _savedCredentialUiState.value = SavedCredentialUiState(
                        isLoading = false,
                        warningMessage = "Không thể tải thông tin đăng nhập đã lưu.",
                    )
                }
        }
    }

    fun login(username: String, password: String, rememberLogin: Boolean = false) {
        if (username.isBlank() || password.isBlank()) {
            _uiState.value = LoginUiState.Error("Vui lòng nhập tài khoản và mật khẩu.")
            return
        }

        viewModelScope.launch {
            _uiState.value = LoginUiState.Loading
            val result = loginToPortal(username, password)
            result.onSuccess {
                updateSavedCredentialsAfterLogin(username, password, rememberLogin)
                _uiState.value = LoginUiState.Success(it)
                return@launch
            }
            _uiState.value = LoginUiState.Error(
                result.exceptionOrNull()?.message ?: "Không thể đăng nhập. Vui lòng thử lại.",
            )
        }
    }

    fun clearSavedCredentials() {
        val store = credentialStore ?: return
        viewModelScope.launch {
            runCatching { store.clear() }
                .onSuccess {
                    _savedCredentialUiState.value = SavedCredentialUiState()
                    _uiState.value = LoginUiState.Idle
                }
                .onFailure {
                    _savedCredentialUiState.value = _savedCredentialUiState.value.copy(
                        warningMessage = "Không thể xóa thông tin đăng nhập đã lưu.",
                    )
                    _uiState.value = LoginUiState.Warning("Không thể xóa thông tin đăng nhập đã lưu.")
                }
        }
    }

    fun resetError() {
        if (_uiState.value is LoginUiState.Error || _uiState.value is LoginUiState.Warning) {
            _uiState.value = LoginUiState.Idle
        }
    }

    private suspend fun updateSavedCredentialsAfterLogin(
        username: String,
        password: String,
        rememberLogin: Boolean,
    ): String? {
        val store = credentialStore ?: return null
        return runCatching {
            if (rememberLogin) {
                store.save(username, password)
                _savedCredentialUiState.value = SavedCredentialUiState(credentials = SavedCredentials(username, password))
            } else {
                store.clear()
                _savedCredentialUiState.value = SavedCredentialUiState()
            }
        }.fold(
            onSuccess = { null },
            onFailure = {
                val message = if (rememberLogin) {
                    "Đăng nhập thành công nhưng không thể lưu thông tin đăng nhập."
                } else {
                    "Đăng nhập thành công nhưng không thể xóa thông tin đăng nhập đã lưu."
                }
                _savedCredentialUiState.value = _savedCredentialUiState.value.copy(warningMessage = message)
                message
            },
        )
    }
}
