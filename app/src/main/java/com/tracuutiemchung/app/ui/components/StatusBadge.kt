package com.tracuutiemchung.app.ui.components

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.tracuutiemchung.app.data.model.RecommendationStatus

@Composable
fun StatusBadge(status: RecommendationStatus) {
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
