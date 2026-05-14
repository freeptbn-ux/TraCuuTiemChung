# Phase 10 — Output Mapping, UI Compat & Golden Parity Tests

## Mục tiêu
Map output nội bộ sang model UI. Tạo bộ golden regression test đảm bảo parity ổn định.
**Bao gồm test đối chiếu Python vs Kotlin bằng dữ liệu thực tế (`Minhkhoi.html`).**

## Phụ thuộc
- Phase 09 hoàn thành.

---

## 10A. Output Mapping
- `MissingItem` → `VaccineRecommendation`:
  - Thêm `statusTags: List<String>` vào `VaccineRecommendation`.
  - Giữ backward compat: `warning`, `recommendedDate` vẫn populate.
  - Map rõ ràng: `["due"]` → `DUE_NOW`, `["info"]` → `DUE_LATER`, etc.
- Rules không có missing items → status `COMPLETED`.
- Xử lý `warnings` list đầy đủ (không chỉ `firstOrNull`).

## 10B. UI Compatibility
- Kiểm tra UI hiện tại (`ResultScreen`) vẫn hiển thị đúng.
- Nếu UI dùng `warning: String?` → giữ field cũ populate first warning.
- Nếu UI dùng `suggestedDate` → giữ populate từ `earliestNextDoseDate?.toString()`.

---

## 10C. Golden Parity Tests (Test Cases tổng hợp)
Tạo bộ test JSON-based cho các ca giả lập:

| Test case | Rules involved |
|---|---|
| Trẻ 4 tháng chưa tiêm gì | Tất cả — kiểm tra gợi ý toàn bộ |
| Trẻ 15 tháng tiêm xong 6in1(4), Prevenar(3), Rota(2), MVVAC(1) | Single + age-dep + alt + MMR |
| MenQuadfi 1 mũi lúc 3 tháng | Meningococcal ACYW 3 mũi + booster |
| Jevax 2 mũi → Imojev 1 mũi | JE mixing/switch |
| Flu lần đầu lúc 7 tuổi | Flu annual (không cần 2 mũi) |
| Flu lần đầu lúc 3 tuổi | Flu cần 2 mũi + annual |

---

## 10D. ⭐ Bài Test Minhkhoi.html — Python vs Kotlin Parity

### Mục đích
Sử dụng dữ liệu tiêm chủng THỰC TẾ từ file `test/Minhkhoi.html` để chạy cùng input qua cả hai engine (Python & Kotlin), sau đó so sánh output từng rule. Đây là bài test quyết định chất lượng parity.

### Thông tin bệnh nhân

| Field | Giá trị |
|---|---|
| Họ tên | Nguyễn Duy Minh Khôi |
| Ngày sinh | 27/07/2020 |
| Giới tính | Nam |
| Mã đối tượng | 106050120200091 |
| System date (giả định) | Ngày chạy test hoặc cố định (ví dụ 12/05/2026) |

### Bảng lịch sử tiêm chủng (29 mũi, trích từ HTML)

| # | Vắc xin | Nhóm bệnh | Mũi | Ngày tiêm |
|---|---------|-----------|-----|-----------|
| 1 | Influvac Tetra 2024/2025 | Cúm | 1 | 16/11/2025 |
| 2 | Varivax | Thủy đậu | 2 | 04/01/2025 |
| 3 | MMR-II | Quai bị, Rubella, Sởi | 2 | 04/01/2025 |
| 4 | Influvac Tetra 2023/2024 | Cúm | 1 | 13/10/2024 |
| 5 | Vaxigrip Tetra 0.5ml | Cúm | 2 | 13/08/2023 |
| 6 | MENACTRA | Não mô cầu 4 tuýp A,C,Y,W-135 | 1 | 21/05/2023 |
| 7 | Imojev | Viêm não Nhật Bản | 2 | 27/11/2022 |
| 8 | Vaxigrip Tetra 0.5ml | Cúm | 1 | 24/07/2022 |
| 9 | Varivax | Thủy đậu | 1 | 02/04/2022 |
| 10 | Infanrix Hexa | Bạch Hầu, Bại liệt, Hib, Ho gà, Uốn ván, Viêm gan B | 4 | 26/12/2021 |
| 11 | MMR-II | Quai bị, Rubella, Sởi | 1 | 27/11/2021 |
| 12 | Influvac (Hộp 1 xy lanh) | Cúm | 2 | 24/10/2021 |
| 13 | Influvac (Hộp 1 xy lanh) | Cúm | 1 | 25/09/2021 |
| 14 | VA - MENGOC - BC | Não mô cầu | 2 | 20/09/2021 |
| 15 | Prevenar 13 | Viêm tai giữa, phế cầu khuẩn | 4 | 29/08/2021 |
| 16 | Imojev | Viêm não Nhật Bản | 1 | 15/08/2021 |
| 17 | VA - MENGOC - BC | Não mô cầu | 1 | 17/07/2021 |
| 18 | MVVAC | Sởi | 1 | 19/06/2021 |
| 19 | Prevenar 13 | Viêm tai giữa, phế cầu khuẩn | 3 | 03/01/2021 |
| 20 | Rota Teq | Rota | 3 | 19/12/2020 |
| 21 | Infanrix Hexa | Bạch Hầu, Bại liệt, Hib, Ho gà, Uốn ván, Viêm gan B | 3 | 19/12/2020 |
| 22 | Prevenar 13 | Viêm tai giữa, phế cầu khuẩn | 2 | 14/11/2020 |
| 23 | Rota Teq | Rota | 2 | 25/10/2020 |
| 24 | Infanrix Hexa | Bạch Hầu, Bại liệt, Hib, Ho gà, Uốn ván, Viêm gan B | 2 | 25/10/2020 |
| 25 | Prevenar 13 | Viêm tai giữa, phế cầu khuẩn | 1 | 11/10/2020 |
| 26 | Rota Teq | Rota | 1 | 26/09/2020 |
| 27 | Infanrix Hexa | Bạch Hầu, Bại liệt, Hib, Ho gà, Uốn ván, Viêm gan B | 1 | 26/09/2020 |
| 28 | IVACTUBER-BCG | Lao | 1 | 20/08/2020 |
| 29 | Viêm gan B sơ sinh | Viêm gan B | 1 | 27/07/2020 |

### Các rule types bài test này cover

| Rule type | Vaccine liên quan trong data | Lý do quan trọng |
|---|---|---|
| `single_series` | Infanrix Hexa (4 mũi), Varivax (2 mũi), BCG (1 mũi), HepB sơ sinh | Kiểm tra đếm mũi, per-dose interval |
| `age_dependent_series` | Prevenar 13 (4 mũi) | Kiểm tra phác đồ phụ thuộc tuổi mũi đầu |
| `mmr_equivalent_group` | MMR-II (mũi 1 + 2), MVVAC (mũi 1) | Kiểm tra MVVAC → MMR interaction, đếm đủ MMR |
| `flu_group` | Influvac/Vaxigrip (nhiều mùa, nhiều mũi) | Kiểm tra flu annual logic, mùa 2021/2022/2023/2024/2025 |
| `meningococcal_acyw_group` | MENACTRA (1 mũi) | Kiểm tra chế độ tiêm 1 mũi (tiêm lúc ~34 tháng) |
| `group_alternative_courses` | Rota Teq (3 mũi) | Kiểm tra Rota course đã hoàn thành |
| `group_alternative_courses_age_range` | Imojev (2 mũi) | Kiểm tra JE series, booster logic |
| VA-MENGOC-BC interaction | VA-MENGOC-BC (2 mũi) + MENACTRA (1 mũi) | Kiểm tra interaction BC ↔ ACYW |
| `standard_vaccines` | Tất cả | Kiểm tra xem có vaccine tiêu chuẩn nào còn thiếu |

### Quy trình thực hiện test

#### Bước 1: Chạy Python engine
```bash
cd /home/skul9x/Desktop/Test_code/VaccineAnalyzer-Pro-main
python3 -c "
from rule_processor import process_all_rules
from config_data import load_vaccine_rules
import json

# Dữ liệu bệnh nhân Minh Khoi
dob = '27/07/2020'
system_date = '12/05/2026'  # cố định
records = [
    {'vaccine': 'Influvac Tetra 2024/2025', 'disease_group': 'Cúm', 'dose_number': 1, 'date': '16/11/2025'},
    {'vaccine': 'Varivax', 'disease_group': 'Thủy đậu', 'dose_number': 2, 'date': '04/01/2025'},
    {'vaccine': 'MMR-II', 'disease_group': 'Quai bị, Rubella, Sởi', 'dose_number': 2, 'date': '04/01/2025'},
    # ... (tất cả 29 mũi)
]

rules = load_vaccine_rules()
results = process_all_rules(dob, system_date, records, rules)
print(json.dumps(results, indent=2, ensure_ascii=False))
"
```
Lưu output vào `test/golden/minhkhoi_python_output.json`.

#### Bước 2: Chạy Kotlin engine
Tạo test case trong `VaccineAnalysisGoldenTest.kt`:
```kotlin
@Test
fun `golden parity - MinhKhoi real data`() {
    // Parse Minhkhoi.html → records
    val html = File("test/Minhkhoi.html").readText()
    val patientInfo = parser.parsePatientInfo(html)
    val records = parser.parseVaccinationRecords(html)

    // Chạy engine với systemDate cố định = 2026-05-12
    val result = engine.analyze(
        patientInfo = patientInfo,
        records = records,
        systemDate = LocalDate.of(2026, 5, 12)
    )

    // So sánh với expected output từ Python
    val expectedJson = File("test/golden/minhkhoi_python_output.json").readText()
    val expected = Json.decodeFromString<List<ExpectedRecommendation>>(expectedJson)

    // Assert per-rule parity
    for (exp in expected) {
        val actual = result.find { it.ruleName == exp.ruleName }
        assertNotNull(actual, "Missing rule: ${exp.ruleName}")
        assertEquals(exp.status, actual.status, "Status mismatch for ${exp.ruleName}")
        assertEquals(exp.nextDoseNumber, actual.nextDoseNumber, "Dose mismatch for ${exp.ruleName}")
        // ... so sánh warning, recommendedDate nếu có
    }
}
```

#### Bước 3: So sánh output
Tạo bảng đối chiếu Python vs Kotlin:

| Rule / Vaccine Group | Python Output | Kotlin Output | Match? |
|---|---|---|---|
| Infanrix Hexa (6in1) | COMPLETED (4/4) | ? | ☐ |
| Prevenar 13 | COMPLETED (4/4) | ? | ☐ |
| Rota Teq | COMPLETED (3/3) | ? | ☐ |
| MMR-II + MVVAC (MMR group) | COMPLETED (tổng ≥2 MMR) | ? | ☐ |
| Varivax (Thủy đậu) | COMPLETED (2/2) | ? | ☐ |
| BCG | COMPLETED (1/1) | ? | ☐ |
| Viêm gan B sơ sinh | COMPLETED (xét trong Infanrix Hexa) | ? | ☐ |
| Imojev (VNBT) | Có thể cần booster | ? | ☐ |
| Flu 2025/2026 | Cần tiêm mùa mới | ? | ☐ |
| MENACTRA (ACYW) | Có thể cần mũi 2 hoặc booster | ? | ☐ |
| VA-MENGOC-BC | COMPLETED + check interaction | ? | ☐ |
| Viêm gan A (HepA) | Chưa tiêm → gợi ý | ? | ☐ |
| HPV | Tuổi > 9 → có thể gợi ý | ? | ☐ |
| DPT booster | Tuổi > 4 → kiểm tra booster | ? | ☐ |

### Kết quả kỳ vọng cho MinhKhoi (tuổi ~5 tuổi 10 tháng tại 12/05/2026)

**Dự kiến COMPLETED (đã hoàn thành):**
- BCG: 1 mũi → ✅
- Viêm gan B sơ sinh: tính trong 6in1 → ✅
- Infanrix Hexa (6in1): 4 mũi → ✅
- Prevenar 13: 4 mũi → ✅
- Rota Teq: 3 mũi → ✅
- MMR group (MVVAC + MMR-II x2): ≥2 mũi MMR → ✅
- Varivax: 2 mũi → ✅
- VA-MENGOC-BC: 2 mũi → ✅

**Dự kiến CẦN TIÊM TIẾP:**
- Flu: Cần mũi mùa 2026/2027 (mùa mới hàng năm)
- JE (Imojev): Có thể cần booster (kiểm tra booster_interval_years)
- MENACTRA (MenACYW): Tiêm 1 mũi lúc 34 tháng → có thể cần mũi 2 hoặc booster lúc 11-12 tuổi
- HepA: Chưa tiêm → gợi ý 2 mũi
- HPV: Tuổi 5 → chưa đến tuổi (≥9 tuổi) → không gợi ý hoặc info
- DPT booster (Td/Tdap): Kiểm tra booster lúc 4-6 tuổi

### Các edge cases đặc biệt trong data này

1. **MVVAC trước MMR-II**: MVVAC(1) tiêm 19/06/2021, MMR-II(1) tiêm 27/11/2021, MMR-II(2) tiêm 04/01/2025.
   - Python phải tính tổng mũi sởi-rubella-quai bị chính xác.
   - MVVAC chỉ cover sởi, không cover rubella/quai bị.

2. **Flu nhiều mùa**: 5 mùa tiêm flu khác nhau (2021, 2022, 2023, 2024, 2025).
   - Mùa 2021: 2 mũi (< 9 tuổi lần đầu → cần 2 mũi).
   - Các mùa sau: 1 mũi/mùa.
   - Python phải xử lý logic này chính xác.

3. **VA-MENGOC-BC (BC) + MENACTRA (ACYW)**: Hai loại não mô cầu khác nhau.
   - Python kiểm tra interaction giữa BC và ACYW.
   - Cần đảm bảo không nhầm lẫn.

4. **Prevenar 13 age-dependent**: 4 mũi (3+1), mũi đầu lúc 2.5 tháng.
   - Kiểm tra phác đồ 3+1 vs 2+1 dựa trên tuổi mũi đầu.

5. **Rota Teq đã hoàn thành**: 3 mũi, mũi cuối lúc ~5 tháng.
   - Không được gợi ý thêm mũi nào.

---

## 10E. Update GoldenTest hiện có
- `VaccineAnalysisGoldenTest.kt` có thể cần update vì engine behavior thay đổi.
- Mỗi update phải ghi chú lý do trong commit.
- Thêm test case `golden parity - MinhKhoi real data` vào suite.

## Tiêu chí xong
- [x] Output mapping test pass.
- [x] UI không crash, thông tin đầy đủ hơn.
- [x] Golden parity tests đại diện nhiều rule type.
- [x] **Bài test MinhKhoi: Python output == Kotlin output cho TẤT CẢ rules.**
- [x] Regression suite ổn định cho phát triển sau này.
- [x] Tạo xong file `test/golden/minhkhoi_python_output.json`.

## Thời gian: ~4-5 giờ (tăng do thêm test MinhKhoi).
