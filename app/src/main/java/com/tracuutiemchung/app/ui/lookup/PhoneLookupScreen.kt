package com.tracuutiemchung.app.ui.lookup

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tracuutiemchung.app.data.model.AnalysisResult
import com.tracuutiemchung.app.data.portal.PortalPatientSummary

@Composable
fun PhoneLookupScreen(
    viewModel: PhoneLookupViewModel,
    onLookupSuccess: (AnalysisResult) -> Unit,
    onSessionExpired: () -> Unit,
) {
    var phoneNumber by remember { mutableStateOf("") }
    val state by viewModel.uiState.collectAsState()

    LaunchedEffect(state) {
        if (state is LookupUiState.Success) {
            onLookupSuccess((state as LookupUiState.Success).result)
        }
    }

    Surface(
        modifier = Modifier.fillMaxSize(),
        color = MaterialTheme.colorScheme.background,
    ) {
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(horizontal = 20.dp, vertical = 28.dp),
            contentAlignment = Alignment.Center,
        ) {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(28.dp),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
                elevation = CardDefaults.cardElevation(defaultElevation = 2.dp),
            ) {
                Column(
                    modifier = Modifier
                        .padding(22.dp)
                        .verticalScroll(rememberScrollState()),
                    verticalArrangement = Arrangement.spacedBy(14.dp),
                ) {
                    Text(
                        text = "Tra cứu lịch sử tiêm",
                        style = MaterialTheme.typography.headlineSmall,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.primary,
                    )
                    Text(
                        text = "Nhập đúng một số điện thoại để xem lịch sử và khuyến nghị.",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                    when (val current = state) {
                        LookupUiState.Idle,
                        is LookupUiState.Error,
                        LookupUiState.Loading,
                        -> SearchForm(
                            phoneNumber = phoneNumber,
                            onPhoneChange = { phoneNumber = it },
                            isLoading = current == LookupUiState.Loading,
                            errorMessage = (current as? LookupUiState.Error)?.message,
                            onSearch = { viewModel.search(phoneNumber) },
                        )
                        is LookupUiState.PatientSelection -> PatientSelectionList(
                            phone = current.phone,
                            patients = current.patients,
                            onSelectPatient = viewModel::selectPatient,
                            onEditPhone = {
                                phoneNumber = current.phone
                                viewModel.resetToSearch()
                            },
                        )
                        is LookupUiState.LoadingDetail -> LoadingDetailContent(current.patient)
                        LookupUiState.SessionExpired -> SessionExpiredContent(onSessionExpired)
                        is LookupUiState.Success -> Unit
                    }
                }
            }
        }
    }
}

@Composable
private fun SearchForm(
    phoneNumber: String,
    onPhoneChange: (String) -> Unit,
    isLoading: Boolean,
    errorMessage: String?,
    onSearch: () -> Unit,
) {
    OutlinedTextField(
        modifier = Modifier.fillMaxWidth(),
        value = phoneNumber,
        onValueChange = onPhoneChange,
        label = { Text("Số điện thoại") },
        singleLine = true,
        enabled = !isLoading,
    )
    Button(
        modifier = Modifier.fillMaxWidth(),
        enabled = !isLoading,
        onClick = onSearch,
    ) {
        if (isLoading) {
            Row(
                horizontalArrangement = Arrangement.spacedBy(10.dp),
                verticalAlignment = Alignment.CenterVertically,
            ) {
                CircularProgressIndicator(
                    modifier = Modifier.semantics { contentDescription = "Đang tra cứu" },
                    color = MaterialTheme.colorScheme.onPrimary,
                    strokeWidth = 2.dp,
                )
                Text("Đang tra cứu...")
            }
        } else {
            Text("Tra cứu")
        }
    }
    errorMessage?.let {
        Text(
            text = it,
            color = MaterialTheme.colorScheme.error,
            style = MaterialTheme.typography.bodyMedium,
        )
    }
}

@Composable
fun PatientSelectionList(
    phone: String,
    patients: List<PortalPatientSummary>,
    onSelectPatient: (PortalPatientSummary) -> Unit,
    onEditPhone: () -> Unit,
) {
    Column(verticalArrangement = Arrangement.spacedBy(10.dp)) {
        Text(
            text = "Chọn đúng người cần tra cứu tiêm chủng",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.SemiBold,
        )
        Text(
            text = "Số điện thoại: $phone • ${patients.size} hồ sơ phù hợp",
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
        )
        patients.forEach { patient ->
            PatientSummaryCard(
                patient = patient,
                onClick = { onSelectPatient(patient) },
            )
        }
        OutlinedButton(
            modifier = Modifier.fillMaxWidth(),
            onClick = onEditPhone,
        ) {
            Text("Nhập số điện thoại khác")
        }
    }
}

@Composable
fun PatientSummaryCard(
    patient: PortalPatientSummary,
    onClick: () -> Unit,
) {
    val secondaryLine = listOfNotNull(
        patient.birthDateOrYear,
        patient.gender,
        patient.phone,
    ).joinToString(" • ")
    val identifier = patient.receivedDate ?: patient.patientId
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .semantics {
                contentDescription = "Hồ sơ ${patient.fullName}, $secondaryLine, ${patient.address.orEmpty()}, mã $identifier"
            },
        onClick = onClick,
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.secondaryContainer),
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(6.dp),
        ) {
            Text(patient.fullName, fontWeight = FontWeight.SemiBold)
            if (secondaryLine.isNotBlank()) {
                Text(secondaryLine, style = MaterialTheme.typography.bodyMedium)
            }
            patient.address?.let { Text(it, style = MaterialTheme.typography.bodySmall) }
            Text(
                text = if (patient.receivedDate != null) {
                    "Ngày tiếp nhận: ${patient.receivedDate}"
                } else {
                    "Mã hồ sơ: ${patient.patientId}"
                },
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSecondaryContainer,
            )
        }
    }
}

@Composable
private fun LoadingDetailContent(patient: PortalPatientSummary) {
    Row(
        modifier = Modifier.fillMaxWidth(),
        horizontalArrangement = Arrangement.spacedBy(12.dp),
        verticalAlignment = Alignment.CenterVertically,
    ) {
        CircularProgressIndicator(
            modifier = Modifier.semantics { contentDescription = "Đang tải chi tiết tiêm chủng" },
        )
        Column {
            Text("Đang tải chi tiết tiêm chủng...")
            Text(
                text = patient.fullName,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurfaceVariant,
            )
        }
    }
}

@Composable
private fun SessionExpiredContent(onSessionExpired: () -> Unit) {
    Text(
        text = "Phiên đăng nhập hết hạn, vui lòng đăng nhập lại.",
        color = MaterialTheme.colorScheme.error,
        style = MaterialTheme.typography.bodyMedium,
    )
    OutlinedButton(
        modifier = Modifier.fillMaxWidth(),
        onClick = onSessionExpired,
    ) {
        Text("Đăng nhập lại")
    }
}
