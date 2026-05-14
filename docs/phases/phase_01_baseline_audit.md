# Phase 01 — Baseline Audit & Gap Analysis

## Mục tiêu
Chốt trạng thái hiện tại trước khi sửa bất kỳ dòng code nào. Xác định rõ gap giữa Engine Kotlin và Engine Python.

## Việc cần làm

### 1. Chạy test hiện có
```bash
./gradlew test --no-daemon
```
Ghi nhận kết quả pass/fail (Ngày 2026-05-12):
- `VaccineAnalysisEngineTest`: **PASS** (4/4 tests)
- `VaccineRuleRepositoryTest`: **PASS** (4/4 tests)
- `VaccineAnalysisGoldenTest`: **PASS** (8/8 tests)
- `DefaultPortalParserLegacyParityTest`: **PASS** (10/10 tests)
- `LoginViewModelTest`: **PASS** (9/9 tests)
- `PhoneLookupViewModelTest`: **PASS** (6/6 tests)
- **Tổng kết**: Tất cả các unit test nội bộ đều PASS. `VncdcPortalClientLiveTest` bị skipped (đúng thiết kế).

### 2. Lập bảng Gap Analysis

#### 2a. Rule types đã hỗ trợ vs chưa hỗ trợ

| Rule Type (Python) | Kotlin hỗ trợ? | Ghi chú |
|---|---|---|
| `single_series` | ⚠️ Cơ bản | Chỉ đếm mũi + interval đơn, thiếu per-dose interval |
| `single_dose_min_age` | ⚠️ Cơ bản | Chung logic với single_series |
| `single_series_min_age` | ⚠️ Cơ bản | Chung logic với single_series |
| `age_dependent_series` | ❌ | Chưa có — Prevenar13, Vaxneuvance, Synflorix dùng type này |
| `mmr_equivalent_group` | ❌ | Chưa có — MMR_Group dùng type này |
| `group_cumulative_unique_doses` | ❌ | Chưa có |
| `group_cumulative_unique_doses_min_age` | ❌ | Chưa có |
| `group_alternative_courses` | ❌ | Chưa có — Rota dùng type này |
| `group_alternative_courses_min_age` | ❌ | Chưa có |
| `group_alternative_courses_age_range` | ❌ | Chưa có — JE_Group, HepA dùng type này |
| `flu_group` | ❌ | Chưa có — Flu dùng type này |
| `meningococcal_acyw_group` | ❌ | Chưa có — MeningococcalACYW_Group |

#### 2b. Rule fields trong JSON mà Kotlin đang bỏ qua (do `ignoreUnknownKeys = true`)

| Field JSON | Kotlin parse? | Python dùng ở đâu |
|---|---|---|
| `min_interval_days: [null, 30, 30, 360]` | ⚠️ Lấy phần tử đầu non-null | `series_checkers.py` — interval theo từng mũi |
| `dose_specific_rules` | ❌ | `series_checkers.py` — tuổi tối thiểu tuyệt đối, alt age range |
| `booster_interval_years` | ❌ | `series_checkers.py` — nhắc lại định kỳ (Jevax) |
| `booster_after_dose_number` | ❌ | `series_checkers.py` — booster áp dụng sau mũi N |
| `booster_max_age_years` | ❌ | `series_checkers.py` — ngừng booster sau tuổi X |
| `rules_by_age` | ❌ | `series_checkers.py` — phác đồ phụ thuộc tuổi mũi đầu |
| `regimens` | ❌ | `group_checkers_special.py` — MMR equivalent group |
| `members` | ❌ | `group_checkers_special.py` — meningococcal ACYW |
| `courses` | ❌ | `group_checkers_alternative.py` — Rota, JE, HepA |
| `interactions` | ❌ | `group_checkers_special.py` — VA-MENGOC-BC ↔ MenQuadfi |
| `provides_measles_protection_group` | ❌ | `series_checkers.py` — MVVAC vs MMR interaction |
| `recognition_keywords` | ✅ (trong normalizer) | `group_checkers_special.py` — Flu matching |
| `max_age_months_to_start_first_dose_group` | ❌ | `group_checkers_alternative.py` — Rota quá tuổi |
| `max_age_months_for_completion_group` | ❌ | `group_checkers_alternative.py` — Rota quá tuổi hoàn thành |
| `initial_series_interval_days` | ❌ | `group_checkers_special.py` — Flu mũi 2 cách mũi 1 |
| `min_age_*_overall` các dạng | ⚠️ Một phần | `rule_checker_utils.py` — min age theo tuần/tháng/năm/ngày |

#### 2c. Logic xử lý Kotlin thiếu so với Python

1. **Per-dose interval**: Kotlin chỉ dùng 1 `intervalDays` cố định. Python dùng `min_interval_days[doseIndex]` — mỗi mũi có interval khác nhau.
2. **First dose age validation**: Python kiểm tra mũi đầu tiêm quá sớm → trả lỗi + restart. Kotlin bỏ qua.
3. **Booster recurring**: Python xử lý booster định kỳ N năm, giới hạn tuổi. Kotlin chưa có.
4. **MVVAC ↔ MMR interaction**: Python kiểm tra vaccine sởi đơn vs phức hợp. Kotlin chưa có.
5. **VA-MENGOC-BC ↔ MenQuadfi interaction**: Python kiểm tra khoảng cách + cảnh báo. Kotlin chưa có.
6. **Pneumococcal special logic**: Python xử lý 4 loại phế cầu xen kẽ, gợi ý Pneumovax23 thay thế. Kotlin chưa có.
7. **JE mixing/switching**: Python xử lý Jevax↔Imojev↔JEEV chuyển đổi. Kotlin chưa có.
8. **Standard vaccines check**: Python kiểm tra vaccine tiêu chuẩn chưa tiêm → gợi ý. Kotlin chưa có.
9. **Status tags**: Python sử dụng `status_tags` list. Kotlin dùng `RecommendationStatus` enum đơn giản.
10. **Age calculation parity**: Python dùng `get_age_at_date(dob, date)` trả `(months, total_days, years)`. Kotlin dùng `ChronoUnit.DAYS.between` chỉ trả tổng ngày.

### 3. Tạo docs output
- File này (`phase_01_baseline_audit.md`) là output chính.
- Cập nhật thêm kết quả test thực tế sau khi chạy.

## File dự kiến chạm
- Không sửa code nào.
- Có thể tạo thêm `docs/gap_analysis_detail.md` nếu cần chi tiết từng rule.

## Tiêu chí xong
- [x] Test baseline đã chạy, kết quả ghi nhận.
- [x] Bảng gap analysis đầy đủ các rule type + field + logic.
- [x] Team/user xác nhận gap analysis trước khi chuyển Phase 02.

## Thời gian ước tính
~30 phút (chạy test + review gap).
