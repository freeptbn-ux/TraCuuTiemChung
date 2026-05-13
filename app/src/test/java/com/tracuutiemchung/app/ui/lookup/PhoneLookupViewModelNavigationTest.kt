package com.tracuutiemchung.app.ui.lookup

import com.tracuutiemchung.app.data.model.AnalysisResult
import com.tracuutiemchung.app.data.model.PatientInfo
import com.tracuutiemchung.app.data.portal.PortalPatientSummary
import com.tracuutiemchung.app.domain.usecase.LookupVaccinationByPhoneUseCase
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Assert.assertTrue
import org.junit.Before
import org.junit.Test
import org.mockito.kotlin.any
import org.mockito.kotlin.doReturn
import org.mockito.kotlin.mock
import org.mockito.kotlin.whenever

@OptIn(ExperimentalCoroutinesApi::class)
class PhoneLookupViewModelNavigationTest {

    private val testDispatcher = StandardTestDispatcher()
    private val useCase: LookupVaccinationByPhoneUseCase = mock()
    private lateinit var viewModel: PhoneLookupViewModel

    @Before
    fun setup() {
        Dispatchers.setMain(testDispatcher)
        viewModel = PhoneLookupViewModel(useCase, testDispatcher)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `selectPatient should update uiState to Success and update lastResult`() = runTest(testDispatcher) {
        // Arrange
        val patient = PortalPatientSummary(
            patientId = "123",
            fullName = "Test User",
            birthDateOrYear = "2000",
            gender = "Nam",
            address = "Hanoi"
        )
        val mockResult = AnalysisResult(
            patientInfo = PatientInfo("Test User", "0123456789"),
            records = emptyList(),
            recommendations = emptyList(),
            warnings = emptyList()
        )
        
        doReturn(Result.success(mockResult)).whenever(useCase).lookupPatient(any())

        // Act
        viewModel.selectPatient(patient)
        testDispatcher.scheduler.advanceUntilIdle()

        // Assert
        assertTrue(viewModel.uiState.value is LookupUiState.Success)
        assertEquals(mockResult, (viewModel.uiState.value as LookupUiState.Success).result)
        assertEquals(mockResult, viewModel.lastResult.value)
    }

    @Test
    fun `resetToSearch should set uiState to Idle but keep lastResult`() = runTest(testDispatcher) {
        // Arrange
        val mockResult = AnalysisResult(
            patientInfo = PatientInfo("Test User", "0123456789"),
            records = emptyList(),
            recommendations = emptyList(),
            warnings = emptyList()
        )
        // Manual set lastResult since we don't have a direct setter
        val patient = PortalPatientSummary(patientId = "123", fullName = "Test User")
        doReturn(Result.success(mockResult)).whenever(useCase).lookupPatient(any())
        viewModel.selectPatient(patient)
        testDispatcher.scheduler.advanceUntilIdle()
        
        // Pre-condition
        assertTrue(viewModel.uiState.value is LookupUiState.Success)
        assertEquals(mockResult, viewModel.lastResult.value)

        // Act
        viewModel.resetToSearch()

        // Assert
        assertEquals(LookupUiState.Idle, viewModel.uiState.value)
        assertEquals(mockResult, viewModel.lastResult.value) // Should persist
    }
}
