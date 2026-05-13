package com.tracuutiemchung.app.ui.lookup

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import com.tracuutiemchung.app.data.portal.PortalPatientSummary
import org.junit.Rule
import org.junit.Test

class PatientSelectionScreenTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun testPatientSelectionScreen_displaysPatients() {
        val patients = listOf(
            PortalPatientSummary(
                patientId = "1",
                fullName = "Nguyen Van A",
                birthDateOrYear = "1990",
                gender = "Nam"
            ),
            PortalPatientSummary(
                patientId = "2",
                fullName = "Tran Thi B",
                birthDateOrYear = "1995",
                gender = "Nu"
            )
        )

        var selectedPatient: PortalPatientSummary? = null

        composeTestRule.setContent {
            PatientSelectionScreen(
                phone = "0123456789",
                patients = patients,
                onSelectPatient = { selectedPatient = it },
                onEditPhone = {}
            )
        }

        // Check if title is displayed
        composeTestRule.onNodeWithText("Chọn Người Cần Tra Cứu").assertIsDisplayed()
        
        // Check if patients are displayed
        composeTestRule.onNodeWithText("Nguyen Van A").assertIsDisplayed()
        composeTestRule.onNodeWithText("Tran Thi B").assertIsDisplayed()

        // Test click
        composeTestRule.onNodeWithText("Nguyen Van A").performClick()
        assert(selectedPatient?.fullName == "Nguyen Van A")
    }

    @Test
    fun testPatientSelectionScreen_editPhoneClick() {
        var editPhoneClicked = false

        composeTestRule.setContent {
            PatientSelectionScreen(
                phone = "0123456789",
                patients = emptyList(),
                onSelectPatient = {},
                onEditPhone = { editPhoneClicked = true }
            )
        }

        composeTestRule.onNodeWithText("Nhập số điện thoại khác").performClick()
        assert(editPhoneClicked)
    }
}
