package com.tracuutiemchung.app.ui.login

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.Checkbox
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
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp

@Composable
fun LoginScreen(
    viewModel: LoginViewModel,
    onLoginSuccess: () -> Unit,
) {
    var username by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var rememberLogin by remember { mutableStateOf(false) }
    val uiState by viewModel.uiState.collectAsState()
    val savedCredentialState by viewModel.savedCredentialUiState.collectAsState()
    val isLoading = uiState is LoginUiState.Loading

    LaunchedEffect(savedCredentialState.credentials) {
        val credentials = savedCredentialState.credentials ?: return@LaunchedEffect
        username = credentials.username
        password = credentials.password
        rememberLogin = true
    }

    LaunchedEffect(uiState) {
        if (uiState is LoginUiState.Success) {
            onLoginSuccess()
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
                    modifier = Modifier.padding(22.dp),
                    verticalArrangement = Arrangement.spacedBy(14.dp),
                ) {
                    Text(
                        text = "Tra cứu tiêm chủng",
                        style = MaterialTheme.typography.headlineSmall,
                        fontWeight = FontWeight.Bold,
                        color = MaterialTheme.colorScheme.primary,
                    )
                    Text(
                        text = "Đăng nhập VNCDC để bắt đầu tra cứu.",
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )
                    Spacer(modifier = Modifier.height(4.dp))
                    OutlinedTextField(
                        modifier = Modifier.fillMaxWidth(),
                        value = username,
                        onValueChange = {
                            username = it
                            viewModel.resetError()
                        },
                        label = { Text("Tài khoản") },
                        singleLine = true,
                        enabled = !isLoading,
                    )
                    OutlinedTextField(
                        modifier = Modifier.fillMaxWidth(),
                        value = password,
                        onValueChange = {
                            password = it
                            viewModel.resetError()
                        },
                        label = { Text("Mật khẩu") },
                        singleLine = true,
                        visualTransformation = PasswordVisualTransformation(),
                        enabled = !isLoading,
                    )

                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        verticalAlignment = Alignment.CenterVertically,
                    ) {
                        Checkbox(
                            checked = rememberLogin,
                            onCheckedChange = { rememberLogin = it },
                            enabled = !isLoading,
                        )
                        Text(
                            text = "Ghi nhớ đăng nhập trên thiết bị này",
                            style = MaterialTheme.typography.bodyMedium,
                            color = MaterialTheme.colorScheme.onSurface,
                        )
                    }
                    Text(
                        text = "Thông tin chỉ được lưu bảo mật trên thiết bị hiện tại.",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant,
                    )

                    if (savedCredentialState.hasSavedCredentials) {
                        OutlinedButton(
                            modifier = Modifier.fillMaxWidth(),
                            enabled = !isLoading,
                            onClick = {
                                rememberLogin = false
                                username = ""
                                password = String()
                                viewModel.clearSavedCredentials()
                            },
                        ) {
                            Text("Xóa thông tin đã lưu")
                        }
                    }

                    val message = when (val state = uiState) {
                        is LoginUiState.Error -> state.message to MaterialTheme.colorScheme.error
                        is LoginUiState.Warning -> state.message to MaterialTheme.colorScheme.error
                        else -> savedCredentialState.warningMessage?.let { it to MaterialTheme.colorScheme.error }
                    }
                    if (message != null) {
                        Text(
                            modifier = Modifier.semantics { contentDescription = "Thông báo đăng nhập" },
                            text = message.first,
                            color = message.second,
                            style = MaterialTheme.typography.bodyMedium,
                        )
                    }

                    Button(
                        modifier = Modifier.fillMaxWidth(),
                        enabled = !isLoading,
                        onClick = { viewModel.login(username, password, rememberLogin) },
                    ) {
                        if (isLoading || savedCredentialState.isLoading) {
                            Row(
                                horizontalArrangement = Arrangement.spacedBy(10.dp),
                                verticalAlignment = Alignment.CenterVertically,
                            ) {
                                CircularProgressIndicator(
                                    modifier = Modifier.semantics { contentDescription = "Đang đăng nhập" },
                                    color = MaterialTheme.colorScheme.onPrimary,
                                    strokeWidth = 2.dp,
                                )
                                Text("Đang xử lý...")
                            }
                        } else {
                            Text("Đăng nhập")
                        }
                    }
                }
            }
        }
    }
}
