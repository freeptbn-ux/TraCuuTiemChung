# TraCuuTiemChung Engine (Go)

Hệ thống tính toán lịch tiêm chủng thông minh được phát triển bằng ngôn ngữ Go, đảm bảo hiệu suất cao và độ chính xác tuyệt đối. Dự án này được thiết kế để thay thế công cụ tính toán cũ bằng Python, mang lại khả năng xử lý nhanh hơn và cấu trúc mã nguồn chặt chẽ hơn.

## 🚀 Tính năng chính

- **Tính toán lịch tiêm chủng chính xác**: Hỗ trợ đầy đủ các quy tắc phức tạp như khoảng cách giữa các mũi tiêm, độ tuổi tối thiểu/tối đa, và các phác đồ tiêm chủng đặc thù.
- **Đồng bộ hóa tuyệt đối (Parity)**: Đạt 100% độ chính xác so với hệ thống cũ, được xác thực thông qua các bộ test tích hợp nghiêm ngặt.
- **Hỗ trợ đa dạng loại vắc xin**:
  - Chuỗi vắc xin cơ bản (Basic Series).
  - Phác đồ phụ thuộc độ tuổi (Age-dependent Series).
  - Nhóm vắc xin thay thế (Alternative Groups - Rota, JE, HepA).
  - Các loại vắc xin đặc biệt (Flu, MenACYW).
- **Xử lý tương tác vắc xin**: Tự động cảnh báo và điều chỉnh lịch khi có sự tương tác giữa các loại vắc xin (như MMR và vắc xin Sởi).

## 🛠 Công nghệ sử dụng

- **Ngôn ngữ**: [Go (Golang) 1.25+](https://go.dev/)
- **Kiểm thử**: [Testify](https://github.com/stretchr/testify) (hỗ trợ suite-based testing).
- **Dữ liệu**: JSON (lưu trữ quy tắc vắc xin tại `vaccine_rules.json`).
- **Legacy Support**: Python (phục vụ đối soát và chuyển đổi dữ liệu).

## 📁 Cấu trúc thư mục

```text
.
├── engine/             # Logic xử lý chính (Checkers, Processor, Utils)
│   ├── checkers/       # Các hàm kiểm tra riêng biệt cho từng loại vắc xin
│   └── utils/          # Các hàm toán học về thời gian và chuẩn hóa dữ liệu
├── models/             # Định nghĩa cấu trúc dữ liệu (Rule, Patient, Recommendation)
├── testdata/           # Dữ liệu mẫu phục vụ kiểm thử tính đồng bộ
├── tests/              # Các bộ test tích hợp và parity tests
├── vercel-backend/     # (Legacy) Backend cũ bằng Python
├── go.mod              # Quản lý dependencies của dự án
└── vaccine_rules.json  # Tệp tin cấu hình các quy tắc tiêm chủng
```

## ⚙️ Cài đặt & Sử dụng

### Yêu cầu hệ thống
- Đã cài đặt Go phiên bản 1.25 trở lên.

### Cài đặt
1. Clone repository:
   ```bash
   git clone https://github.com/freeptbn-ux/TraCuuTiemChung.git
   cd TraCuuTiemChung
   ```
2. Cài đặt các thư viện cần thiết:
   ```bash
   go mod tidy
   ```

### Chạy kiểm thử
Để chạy toàn bộ các bộ test, bao gồm cả test tính đồng bộ với engine Python:
```bash
go test -v ./tests/parity/...
```

## 📄 Bản quyền

Copyright 2026 Nguyễn Duy Trường

---
*Dự án được phát triển với mục tiêu nâng cao sức khỏe cộng đồng thông qua việc cung cấp lịch tiêm chủng chính xác.*
