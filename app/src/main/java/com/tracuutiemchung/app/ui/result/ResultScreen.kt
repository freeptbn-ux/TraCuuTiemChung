package com.tracuutiemchung.app.ui.result

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tracuutiemchung.app.data.model.AnalysisResult
import com.tracuutiemchung.app.data.model.RecommendationStatus

@Composable
fun ResultScreen(result: AnalysisResult) {
    Surface(
        modifier = Modifier.fillMaxSize(),
        color = MaterialTheme.colorScheme.background,
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(horizontal = 18.dp, vertical = 20.dp),
            verticalArrangement = Arrangement.spacedBy(12.dp),
        ) {
            Text(
                "Kết quả phân tích",
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Bold,
                color = MaterialTheme.colorScheme.primary,
            )
            InfoCard("Người tiêm") {
                Text("Họ tên: ${result.patientInfo.fullName}")
                Text("SĐT: ${result.patientInfo.phoneNumber}")
                result.patientInfo.dateOfBirth?.let { Text("Ngày sinh: $it") }
                result.patientInfo.address?.let { Text("Địa chỉ: $it") }
            }
            InfoCard("Mũi đã tiêm") {
                if (result.records.isEmpty()) {
                    Text("Chưa có dữ liệu mũi đã tiêm.", color = MaterialTheme.colorScheme.onSurfaceVariant)
                } else {
                    result.records.forEach { record ->
                        Column(verticalArrangement = Arrangement.spacedBy(2.dp)) {
                            Text(record.vaccineName, fontWeight = FontWeight.SemiBold)
                            Text(
                                text = listOfNotNull(
                                    record.doseNumber?.let { "Mũi $it" },
                                    record.vaccinationDate,
                                    record.provider,
                                ).joinToString(" • "),
                                style = MaterialTheme.typography.bodyMedium,
                                color = MaterialTheme.colorScheme.onSurfaceVariant,
                            )
                        }
                    }
                }
            }
            InfoCard("Mũi cần tiêm tiếp") {
                if (result.recommendations.isEmpty()) {
                    Text("Chưa có khuyến nghị cần tiêm tiếp.", color = MaterialTheme.colorScheme.onSurfaceVariant)
                } else {
                    result.recommendations.forEach { recommendation ->
                        Column(verticalArrangement = Arrangement.spacedBy(6.dp)) {
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.SpaceBetween,
                                verticalAlignment = Alignment.CenterVertically,
                            ) {
                                Text(
                                    modifier = Modifier.weight(1f),
                                    text = recommendation.vaccineName,
                                    fontWeight = FontWeight.SemiBold,
                                )
                                StatusBadge(recommendation.status)
                            }
                            Text(recommendation.reason, style = MaterialTheme.typography.bodyMedium)
                            recommendation.recommendedDate?.let {
                                Text(
                                    text = "Ngày nên tiêm: $it",
                                    style = MaterialTheme.typography.bodyMedium,
                                    color = MaterialTheme.colorScheme.primary,
                                    fontWeight = FontWeight.Medium,
                                )
                            }
                            recommendation.warnings.forEach { warning ->
                                Text(
                                    text = "Cảnh báo: $warning",
                                    color = MaterialTheme.colorScheme.error,
                                    style = MaterialTheme.typography.bodyMedium,
                                )
                            }
                        }
                    }
                }
            }
            if (result.warnings.isNotEmpty()) {
                InfoCard("Cảnh báo thiếu dữ liệu") {
                    result.warnings.forEach {
                        Text("• $it", color = MaterialTheme.colorScheme.error)
                    }
                }
            }
        }
    }
}

@Composable
private fun InfoCard(title: String, content: @Composable () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(22.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface),
        elevation = CardDefaults.cardElevation(defaultElevation = 1.dp),
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp),
        ) {
            Text(title, style = MaterialTheme.typography.titleMedium, fontWeight = FontWeight.Bold)
            content()
        }
    }
}

@Composable
private fun StatusBadge(status: RecommendationStatus) {
    val (label, color) = when (status) {
        RecommendationStatus.COMPLETED -> "Đã đủ" to Color(0xFF1B6E3C)
        RecommendationStatus.DUE_NOW -> "Cần tiêm" to Color(0xFF9A5B00)
        RecommendationStatus.DUE_LATER -> "Sắp tới" to Color(0xFF37618E)
        RecommendationStatus.OVERDUE -> "Quá hạn" to Color(0xFFBA1A1A)
        RecommendationStatus.NEEDS_REVIEW -> "Cần xem lại" to Color(0xFF6D5E00)
        RecommendationStatus.NOT_ENOUGH_DATA -> "Thiếu dữ liệu" to Color(0xFF6D5E00)
    }
    Text(
        modifier = Modifier
            .background(color.copy(alpha = 0.14f), RoundedCornerShape(999.dp))
            .padding(horizontal = 10.dp, vertical = 5.dp),
        text = label,
        color = color,
        style = MaterialTheme.typography.labelMedium,
        fontWeight = FontWeight.Bold,
    )
}
