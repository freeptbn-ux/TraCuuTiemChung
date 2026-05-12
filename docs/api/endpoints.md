# API Documentation

Ngày cập nhật: 12/05/2026
Base URL: `https://tracuutiemchung-api.vercel.app` (Example)

---

## 🔐 Security
Tất cả các endpoint yêu cầu API Key được gửi qua Header:
`X-API-KEY: <YOUR_SECRET_KEY>`

---

## 📋 Endpoints

### POST /api/lookup
Tra cứu danh sách bệnh nhân theo số điện thoại.

**Request Body:**
```json
{
  "phone": "0912345678"
}
```

**Response (200):**
```json
{
  "status": "success",
  "data": [
    {
      "patient_id": "12345",
      "name": "Nguyễn Văn A",
      "dob": "01/01/2020",
      "code": "78910"
    }
  ]
}
```

---

### POST /api/analyze
Phân tích lịch sử tiêm chủng và đưa ra khuyến nghị.

**Request Body:**
```json
{
  "patient_id": "12345"
}
```

**Response (200):**
```json
{
  "status": "success",
  "patient_name": "Nguyễn Văn A",
  "dob": "01/01/2020",
  "analysis_date": "12/05/2026",
  "missing_vaccines": [
    {
      "vaccine_name_for_popup": "Varivax (Thủy đậu)",
      "earliest_next_dose_date": "2026-05-12",
      "status_tags": ["eligible"],
      "description": "Varivax (Thủy đậu) (Chưa tiêm - cần 2 liều). đủ điều kiện tuổi"
    }
  ]
}
```

---

## ⚙️ Error Codes
- `401`: API Key không hợp lệ hoặc thiếu.
- `500`: Lỗi hệ thống hoặc lỗi kết nối tới VNCDC Portal.
- `404`: Không tìm thấy dữ liệu.
