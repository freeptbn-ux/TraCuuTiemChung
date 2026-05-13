package com.tracuutiemchung.app.ui.lookup

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.ui.text.input.KeyboardType
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

import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Phone
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.CenterAlignedTopAppBar
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.Scaffold
import androidx.compose.material3.TopAppBarDefaults
import androidx.compose.ui.text.style.TextOverflow

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PhoneLookupScreen(
    viewModel: PhoneLookupViewModel,
    onNavigateToSelection: (String) -> Unit,
    onNavigateToResult: (String) -> Unit,
    onSessionExpired: () -> Unit,
) {
    var phoneNumber by remember { mutableStateOf("") }
    val state by viewModel.uiState.collectAsState()

    LaunchedEffect(state) {
        when (state) {
            is LookupUiState.PatientSelection -> {
                onNavigateToSelection((state as LookupUiState.PatientSelection).phone)
            }
            is LookupUiState.Success -> {
                onNavigateToResult((state as LookupUiState.Success).result.patientInfo.phoneNumber)
            }
            else -> Unit
        }
    }

    Scaffold(
        topBar = {
            CenterAlignedTopAppBar(
                title = {
                    Text(
                        "Tra Cứu Tiêm Chủng",
                        style = MaterialTheme.typography.titleLarge,
                        fontWeight = FontWeight.Bold
                    )
                },
                colors = TopAppBarDefaults.centerAlignedTopAppBarColors(
                    containerColor = MaterialTheme.colorScheme.surface,
                    titleContentColor = MaterialTheme.colorScheme.primary,
                )
            )
        }
    ) { paddingValues ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(16.dp),
            contentAlignment = Alignment.TopCenter,
        ) {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .verticalScroll(rememberScrollState()),
                horizontalAlignment = Alignment.CenterHorizontally,
                verticalArrangement = Arrangement.spacedBy(24.dp)
            ) {
                // Header section
                Column(
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    Text(
                        text = "Chào mừng bạn",
                        style = MaterialTheme.typography.headlineSmall,
                        fontWeight = FontWeight.ExtraBold,
                        color = MaterialTheme.colorScheme.onSurface
                    )
                    Text(
                        text = "Tra cứu lịch sử và nhận khuyến nghị tiêm chủng",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                        textAlign = androidx.compose.ui.text.style.TextAlign.Center
                    )
                }

                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = MaterialTheme.shapes.extraLarge,
                    colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
                    elevation = CardDefaults.cardElevation(defaultElevation = 0.dp),
                    border = androidx.compose.foundation.BorderStroke(
                        1.dp,
                        MaterialTheme.colorScheme.outlineVariant
                    )
                ) {
                    Column(
                        modifier = Modifier.padding(24.dp),
                        verticalArrangement = Arrangement.spacedBy(20.dp)
                    ) {
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
                            is LookupUiState.PatientSelection,
                            is LookupUiState.LoadingDetail -> Unit
                            LookupUiState.SessionExpired -> SessionExpiredContent(onSessionExpired)
                            is LookupUiState.Success -> Unit
                        }
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
    Column(verticalArrangement = Arrangement.spacedBy(16.dp)) {
        OutlinedTextField(
            modifier = Modifier.fillMaxWidth(),
            value = phoneNumber,
            onValueChange = onPhoneChange,
            label = { Text("Số điện thoại") },
            placeholder = { Text("0123 456 789") },
            leadingIcon = {
                Icon(
                    imageVector = Icons.Default.Phone,
                    contentDescription = null,
                    tint = MaterialTheme.colorScheme.primary
                )
            },
            singleLine = true,
            enabled = !isLoading,
            shape = MaterialTheme.shapes.medium,
            isError = errorMessage != null,
            keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(
                keyboardType = androidx.compose.ui.text.input.KeyboardType.Phone
            )
        )
        
        errorMessage?.let {
            Text(
                text = it,
                color = MaterialTheme.colorScheme.error,
                style = MaterialTheme.typography.bodySmall,
                modifier = Modifier.padding(start = 4.dp)
            )
        }

        Button(
            modifier = Modifier
                .fillMaxWidth()
                .height(56.dp),
            enabled = !isLoading && phoneNumber.isNotBlank(),
            onClick = onSearch,
            shape = MaterialTheme.shapes.medium,
            contentPadding = androidx.compose.foundation.layout.PaddingValues(horizontal = 24.dp)
        ) {
            if (isLoading) {
                Row(
                    horizontalArrangement = Arrangement.spacedBy(12.dp),
                    verticalAlignment = Alignment.CenterVertically,
                ) {
                    CircularProgressIndicator(
                        modifier = Modifier
                            .size(24.dp)
                            .semantics { contentDescription = "Đang tra cứu" },
                        color = MaterialTheme.colorScheme.onPrimary,
                        strokeWidth = 2.5.dp,
                    )
                    Text("Đang xử lý...", style = MaterialTheme.typography.titleMedium)
                }
            } else {
                Row(
                    horizontalArrangement = Arrangement.spacedBy(12.dp),
                    verticalAlignment = Alignment.CenterVertically,
                ) {
                    Icon(imageVector = Icons.Default.Search, contentDescription = null)
                    Text("Tra cứu ngay", style = MaterialTheme.typography.titleMedium)
                }
            }
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
