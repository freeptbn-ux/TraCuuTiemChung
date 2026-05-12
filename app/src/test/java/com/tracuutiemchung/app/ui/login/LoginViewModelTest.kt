package com.tracuutiemchung.app.ui.login

import com.tracuutiemchung.app.data.credentials.CredentialStore
import com.tracuutiemchung.app.data.credentials.SavedCredentials
import com.tracuutiemchung.app.data.portal.InMemorySessionStore
import com.tracuutiemchung.app.data.portal.LocalPortalAuthRepository
import com.tracuutiemchung.app.data.portal.fakeLogin
import com.tracuutiemchung.app.domain.usecase.LoginToPortalUseCase
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class LoginViewModelTest {
    private val dispatcher = StandardTestDispatcher()

    @Before
    fun setUp() {
        Dispatchers.setMain(dispatcher)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun blankInputShowsErrorWithoutLogin() = runTest {
        val viewModel = createViewModel()

        viewModel.login("", "pass")

        val state = viewModel.uiState.value
        assertTrue(state is LoginUiState.Error)
        assertEquals("Vui lòng nhập tài khoản và mật khẩu.", (state as LoginUiState.Error).message)
    }

    @Test
    fun validInputEmitsSuccess() = runTest {
        val viewModel = createViewModel()

        viewModel.login("dev-user", "dev-pass")
        advanceUntilIdle()

        assertTrue(viewModel.uiState.value is LoginUiState.Success)
    }

    @Test
    fun initLoadsSavedCredentialsIntoSavedCredentialUiState() = runTest {
        val credentialStore = FakeCredentialStore(SavedCredentials("saved-user", "saved-pass"))
        val viewModel = createViewModel(credentialStore)

        advanceUntilIdle()

        assertFalse(viewModel.savedCredentialUiState.value.isLoading)
        assertTrue(viewModel.savedCredentialUiState.value.hasSavedCredentials)
        assertEquals(
            SavedCredentials("saved-user", "saved-pass"),
            viewModel.savedCredentialUiState.value.credentials,
        )
    }

    @Test
    fun noSavedCredentialsLeavesSavedCredentialUiStateBlank() = runTest {
        val credentialStore = FakeCredentialStore()
        val viewModel = createViewModel(credentialStore)

        advanceUntilIdle()

        assertFalse(viewModel.savedCredentialUiState.value.isLoading)
        assertFalse(viewModel.savedCredentialUiState.value.hasSavedCredentials)
        assertEquals(null, viewModel.savedCredentialUiState.value.credentials)
    }

    @Test
    fun successfulLoginSavesOnlyWhenRememberLoginIsEnabled() = runTest {
        val credentialStore = FakeCredentialStore()
        val viewModel = createViewModel(credentialStore)

        advanceUntilIdle()
        viewModel.login("dev-user", "dev-pass", rememberLogin = true)
        advanceUntilIdle()

        assertEquals(SavedCredentials("dev-user", "dev-pass"), credentialStore.savedCredentials.value)
        assertEquals(SavedCredentials("dev-user", "dev-pass"), viewModel.savedCredentialUiState.value.credentials)
    }

    @Test
    fun successfulLoginClearsSavedCredentialWhenRememberLoginIsDisabled() = runTest {
        val credentialStore = FakeCredentialStore(SavedCredentials("old-user", "old-pass"))
        val viewModel = createViewModel(credentialStore)

        advanceUntilIdle()
        viewModel.login("dev-user", "dev-pass", rememberLogin = false)
        advanceUntilIdle()

        assertEquals(null, credentialStore.savedCredentials.value)
        assertFalse(viewModel.savedCredentialUiState.value.hasSavedCredentials)
    }

    @Test
    fun clearSavedCredentialsRemovesStoredCredentialsImmediately() = runTest {
        val credentialStore = FakeCredentialStore(SavedCredentials("old-user", "old-pass"))
        val viewModel = createViewModel(credentialStore)

        advanceUntilIdle()
        viewModel.clearSavedCredentials()
        advanceUntilIdle()

        assertEquals(null, credentialStore.savedCredentials.value)
        assertFalse(viewModel.savedCredentialUiState.value.hasSavedCredentials)
    }

    @Test
    fun credentialStoreSaveFailureDoesNotPreventSuccessfulLogin() = runTest {
        val credentialStore = FakeCredentialStore(failSave = true)
        val viewModel = createViewModel(credentialStore)

        advanceUntilIdle()
        viewModel.login("dev-user", "dev-pass", rememberLogin = true)
        advanceUntilIdle()

        assertTrue(viewModel.uiState.value is LoginUiState.Success)
        assertEquals(
            "Đăng nhập thành công nhưng không thể lưu thông tin đăng nhập.",
            viewModel.savedCredentialUiState.value.warningMessage,
        )
        assertEquals(null, credentialStore.savedCredentials.value)
    }

    @Test
    fun credentialStoreClearFailureDoesNotPreventSuccessfulLogin() = runTest {
        val credentialStore = FakeCredentialStore(
            initial = SavedCredentials("old-user", "old-pass"),
            failClear = true,
        )
        val viewModel = createViewModel(credentialStore)

        advanceUntilIdle()
        viewModel.login("dev-user", "dev-pass", rememberLogin = false)
        advanceUntilIdle()

        assertTrue(viewModel.uiState.value is LoginUiState.Success)
        assertEquals(
            "Đăng nhập thành công nhưng không thể xóa thông tin đăng nhập đã lưu.",
            viewModel.savedCredentialUiState.value.warningMessage,
        )
        assertEquals(SavedCredentials("old-user", "old-pass"), credentialStore.savedCredentials.value)
    }

    private fun createViewModel(credentialStore: CredentialStore? = null): LoginViewModel {
        val repository = LocalPortalAuthRepository(InMemorySessionStore(), ::fakeLogin)
        return LoginViewModel(LoginToPortalUseCase(repository), credentialStore)
    }

    private class FakeCredentialStore(
        initial: SavedCredentials? = null,
        private val failSave: Boolean = false,
        private val failClear: Boolean = false,
    ) : CredentialStore {
        override val savedCredentials = MutableStateFlow(initial)

        override suspend fun save(username: String, password: String) {
            if (failSave) error("store failure")
            savedCredentials.value = SavedCredentials(username, password)
        }

        override suspend fun clear() {
            if (failClear) error("store failure")
            savedCredentials.value = null
        }
    }
}
