# Plan nâng cấp Vaccine Analysis Engine

## Mục tiêu

Nâng cấp logic phân tích tiêm chủng trong dự án Android `TraCuuTiemChung`
để tiến gần parity với dự án Python `VaccineAnalyzer-Pro-main`, nhưng làm theo
các phase nhỏ, dễ review, dễ test và hạn chế phá app hiện tại.

> [!IMPORTANT]
> Theo yêu cầu: tài liệu này chỉ là plan. Sau khi tạo plan thì dừng để bạn xem xét,
> chưa sửa code Kotlin/JSON/test.

## Bối cảnh hiện tại

### App Android hiện có

- Engine chính: `app/src/main/java/.../domain/analyzer/VaccineAnalysisEngine.kt`
- Rule loader: `app/src/main/java/.../data/rules/VaccineRuleRepository.kt`
- Model output: `app/src/main/java/.../data/model/Models.kt`
- Rule asset: `app/src/main/assets/vaccine_rules.json`
- Test hiện có:
  - `VaccineAnalysisEngineTest.kt`
  - `VaccineRuleRepositoryTest.kt`
  - `VaccineAnalysisGoldenTest.kt`

Engine Kotlin hiện mới xử lý dạng đơn giản:

- Đếm số mũi theo `requiredDoses`.
- Tính ngày gợi ý bằng `minAgeDays` và `intervalDays` đầu tiên.
- Chưa hỗ trợ đầy đủ:
  - `min_interval_days` theo từng mũi.
  - `dose_specific_rules`.
  - booster định kỳ.
  - age-dependent series.
  - group rules: MMR equivalent, cumulative, alternative courses, flu, meningococcal.
  - status tags tương đương Python.

### Dự án Python nguồn tham chiếu

Các file logic quan trọng:

- `rule_processor.py`: vòng lặp điều phối theo `rule_type`.
- `series_checkers.py`: single series, min-age, booster, age-dependent.
- `rule_checker_utils.py`: helper cộng tháng/năm, tuổi, kiểm tra tuổi mũi 1.
- `group_checkers_special.py`: MMR, cumulative, flu, meningococcal ACYW.
- `group_checkers_alternative.py`: alternative courses như viêm não Nhật Bản.
- Test parity tham khảo:
  - `test_menquadfi.py`
  - `test_mvvac_mmr_interval.py`
  - `test_flu_annual.py`
  - `test_booking_alias.py`

## Nguyên tắc triển khai

1. Mỗi phase phải chạy test riêng trước khi chuyển phase kế tiếp.
2. Không rewrite toàn bộ engine trong một lần.
3. Ưu tiên giữ API UI hiện có, chỉ mở rộng model khi cần.
4. Khi logic y khoa/phác đồ không chắc, tra cứu bằng `search_web` hoặc
   `read_url_content`; không mở browser lấy DOM trực tiếp.
5. Rule port phải có test đối chiếu với Python fixtures hoặc kỳ vọng rõ ràng.
6. Không tự động thay đổi dữ liệu phác đồ nếu chưa có test hoặc nguồn xác nhận.

---

## Phase 01 — Baseline audit và test snapshot

### Mục tiêu

Chốt trạng thái hiện tại trước khi sửa logic, đảm bảo có điểm rollback rõ.

### Việc cần làm

- Chạy test hiện có của Android:
  - `./gradlew test`
- Ghi nhận test pass/fail hiện tại.
- Đọc nhanh các golden tests để hiểu format output mong muốn.
- Tạo checklist parity giữa Kotlin và Python:
  - rule type nào đã hỗ trợ.
  - rule field nào đã parse.
  - rule field nào đang bị bỏ qua do `ignoreUnknownKeys`.

### File dự kiến chạm

- Chưa sửa code.
- Có thể tạo `docs/analysis_engine_gap.md` nếu cần ghi audit chi tiết.

### Tiêu chí xong

- Biết test baseline.
- Có danh sách gap rõ ràng trước khi code.

---

## Phase 02 — Mở rộng rule schema Kotlin

### Mục tiêu

Đọc đủ dữ liệu từ `vaccine_rules.json` để engine có nguyên liệu xử lý đúng.

### Việc cần làm

Mở rộng `VaccineRule` và `RawVaccineRule` trong
`VaccineRuleRepository.kt` để parse thêm các field Python đang dùng:

- `doses_required`
- `min_interval_days: List<Int?>`
- `dose_specific_rules`
- `booster_interval_years`
- `booster_after_dose_number`
- `booster_max_age_years`
- `min_age_days_at_first_dose`
- `min_age_weeks_at_first_dose`
- `min_age_months_at_first_dose`
- `min_age_years_at_first_dose`
- `min_age_*_overall`
- `min_age_*_overall_group`
- `rules_by_age`
- `members`
- `courses`
- `names_norm_group`
- các flag interaction như `provides_measles_protection` nếu có trong asset.

Khuyến nghị tạo data class phụ:

- `DoseSpecificRule`
- `AgeRule`
- `CourseRule`
- `MemberRule`
- `RuleAgeConstraint`

### File dự kiến chạm

- `VaccineRuleRepository.kt`
- `VaccineRuleRepositoryTest.kt`

### Test cần thêm/sửa

- Test parse `minIntervalDays` giữ đúng list, không lấy mỗi phần tử đầu.
- Test parse `boosterIntervalYears` cho rule có booster.
- Test parse `rulesByAge` cho rule age-dependent.
- Test normalizer vẫn map đúng tên VNCDC.

### Tiêu chí xong

- Không đổi behavior engine.
- Chỉ mở rộng schema và test schema pass.

---

## Phase 03 — Port helper ngày/tuổi và status nội bộ

### Mục tiêu

Tạo nền tảng helper giống Python để tránh sai lệch khi tính lịch.

### Việc cần làm

Trong domain analyzer, thêm helper tương đương Python:

- `addMonths(LocalDate, months)`:
  - xử lý cuối tháng giống Python `calendar.monthrange`.
- `addYears(LocalDate, years)`:
  - xử lý ngày 29/02 về 28/02 nếu cần.
- `getAgeAtDate(dob, atDate)`:
  - trả về tháng tuổi, tổng ngày tuổi, năm tuổi.
- `getAgeStatusAndEarliestDate(...)`:
  - giống `_get_age_status_and_earliest_date`.
- `checkFirstDoseAgeValidity(...)`.

Tạo model nội bộ để giữ parity với Python:

```kotlin
data class MissingItem(
    val vaccineNameForPopup: String,
    val description: String,
    val earliestNextDoseDate: LocalDate?,
    val statusTags: List<String>,
)
```

Sau đó map `MissingItem` sang `VaccineRecommendation` để UI hiện tại vẫn dùng được.

### File dự kiến chạm

- `VaccineAnalysisEngine.kt`
- Có thể tách file mới:
  - `VaccineDateUtils.kt`
  - `AnalysisRuleUtils.kt`

### Test cần thêm

- `addMonths(2024-01-31, 1) == 2024-02-29`
- `addMonths(2025-01-31, 1) == 2025-02-28`
- `addYears(2024-02-29, 1) == 2025-02-28`
- Tuổi tại ngày trước DOB trả về invalid/null.
- Min age tháng/tuần/ngày/năm trả đúng earliest date.

### Tiêu chí xong

- Helper pass test độc lập.
- Engine output hiện tại không bị đổi ngoài format warning nếu có chủ ý.

---

## Phase 04 — Single series parity

### Mục tiêu

Port logic `check_single_vaccine_series` từ Python trước vì đây là lõi nhiều rule.

### Việc cần làm

Thay logic đơn giản hiện tại bằng luồng tương đương Python cho các type:

- `single_series`
- `single_dose_min_age`
- `single_series_min_age`

Các điểm quan trọng:

- Count tất cả mũi đã ghi nhận, không loại vì khoảng cách cũ.
- Nếu đủ mũi:
  - kiểm tra booster nếu rule có `booster_interval_years`.
  - nếu chưa đến hạn booster: status `DUE_LATER` hoặc info/upcoming.
  - nếu đến hạn booster: status `DUE_NOW`.
- Nếu chưa tiêm:
  - dùng age helper để tính ngày đủ điều kiện mũi đầu.
- Nếu mũi đầu quá sớm:
  - trả warning/error tương ứng, không im lặng coi là hợp lệ.
- Nếu thiếu mũi:
  - lấy interval theo đúng index `min_interval_days[nextDoseIndex]`.
  - áp dụng `dose_specific_rules` nếu có.
  - ngày gợi ý không được trước `analysisDate`.

### File dự kiến chạm

- `VaccineAnalysisEngine.kt`
- `Models.kt` nếu cần thêm `statusTags` vào `VaccineRecommendation`.

Khuyến nghị thêm field tương thích:

```kotlin
val statusTags: List<String> = emptyList()
```

và giữ các field cũ `warning`, `recommendedDate` để không phá UI.

### Test cần thêm

Port các case từ Python:

- MVVAC/MMR interval cơ bản từ `test_mvvac_mmr_interval.py`.
- Booster recurring trong `series_checkers.py`.
- Mũi đầu quá sớm.
- Rule có `dose_specific_rules`.

### Tiêu chí xong

- Test engine cũ vẫn pass hoặc được update có lý do.
- Single series parity đạt với các fixture chọn lọc.

---

## Phase 05 — Age-dependent series

### Mục tiêu

Port `check_age_dependent_series` cho các vaccine có phác đồ phụ thuộc tuổi mũi đầu.

### Việc cần làm

- Dựa vào tuổi tại mũi đầu để chọn `rules_by_age` phù hợp.
- Nếu không có mũi nào:
  - trả gợi ý theo tuổi hiện tại và default dose.
- Nếu thiếu DOB:
  - trả `NOT_ENOUGH_DATA`/error DOB.
- Nếu mũi đầu không khớp age rule:
  - phân biệt lỗi cấu hình, quá sớm, không có rule phù hợp.
- Tái sử dụng single series bằng temp rule như Python.

### File dự kiến chạm

- `VaccineAnalysisEngine.kt`
- Có thể tách checker mới:
  - `AgeDependentSeriesChecker.kt`

### Test cần thêm

- Một rule age-dependent trong `vaccine_rules.json`.
- Case trẻ bắt đầu ở các mốc tuổi khác nhau.
- Case thiếu DOB.
- Case first dose quá sớm.

### Tiêu chí xong

- Age-dependent rule output đúng ngày cần tiêm kế tiếp.
- Không làm sai single series.

---

## Phase 06 — Group rules ưu tiên cao

### Mục tiêu

Port các group checker ảnh hưởng trực tiếp đến kết quả thực tế.

### Thứ tự đề xuất

1. `mmr_equivalent_group`
2. `group_cumulative_unique`
3. `group_cumulative_unique_min_age`
4. `flu_group`
5. `meningococcal_acyw_group`

### Việc cần làm

- Tạo dispatcher giống `rule_processor.py` theo `rule.type`.
- Với MMR:
  - tính vaccine cung cấp bảo vệ sởi.
  - xử lý MVVAC vs MMR interaction.
- Với cumulative:
  - đếm liều duy nhất theo nhóm, không double count tên alias.
- Với flu:
  - xử lý logic yearly/annual nếu Python đang làm vậy.
- Với meningococcal ACYW:
  - port từ `group_checkers_special.py` và test bằng `test_menquadfi.py`.

### File dự kiến chạm

- `VaccineAnalysisEngine.kt`
- Có thể tách:
  - `GroupRuleCheckers.kt`
  - `RuleDispatcher.kt`

### Test cần thêm

- Port chọn lọc từ:
  - `test_menquadfi.py`
  - `test_mvvac_mmr_interval.py`
  - `test_flu_annual.py`

### Tiêu chí xong

- Các group chính có test parity.
- Dispatcher bỏ qua rule chưa hỗ trợ bằng warning rõ, không crash.

---

## Phase 07 — Alternative courses và special pneumococcal logic

### Mục tiêu

Port các logic khó/đặc thù sau khi nền tảng đã ổn.

### Việc cần làm

Từ `group_checkers_alternative.py`:

- `group_alternative`
- `group_alternative_min_age`
- `group_alternative_age_range`
- Japanese encephalitis course switching:
  - Jevax
  - Imojev
  - JEEV

Từ `rule_processor.py`:

- Special pneumococcal logic:
  - Prevenar13
  - Vaxneuvance
  - Synflorix
  - Pneumovax23
- Mixed pneumococcal warning.
- Alternative completion after age 2.

### File dự kiến chạm

- `VaccineAnalysisEngine.kt`
- `GroupAlternativeCheckers.kt`
- `SpecialPneumococcalChecker.kt`

### Test cần thêm

- Cases đã có trong Python nếu tìm thấy.
- Nếu thiếu test Python, tạo fixtures nhỏ dựa trên rule JSON và expected behavior từ code Python.

### Khi cần tra online

Nếu không chắc phác đồ đổi vắc xin hoặc khoảng cách:

- dùng `search_web` để tìm tài liệu chính thống.
- dùng `read_url_content` để đọc trang tài liệu.
- ưu tiên nguồn:
  - Bộ Y tế / VNCDC / WHO / CDC / nhà sản xuất.

### Tiêu chí xong

- Rule alternative không làm sai các rule đã pass.
- Các cảnh báo xen kẽ vaccine hiện đúng và có status tags.

---

## Phase 08 — Mapping output sang UI và UX an toàn

### Mục tiêu

Kết quả phân tích đầy đủ hơn nhưng UI vẫn dễ hiểu.

### Việc cần làm

- Map `statusTags` sang `RecommendationStatus`:
  - `error_*` → `NEEDS_REVIEW`
  - `too_young`, `booster_upcoming`, `info` → `DUE_LATER` hoặc `NEEDS_REVIEW` tùy nội dung.
  - `due`, `booster_due`, `eligible` → `DUE_NOW`
  - completed/no missing item → `COMPLETED`
- Hiển thị lý do rõ ràng:
  - "Cần thêm X liều"
  - "Mũi N cách mũi trước tối thiểu..."
  - "Chưa đủ tuổi, có thể tiêm từ..."
- Nếu có nhiều warning:
  - không chỉ lấy `firstOrNull`; UI nên có khả năng xem toàn bộ.

### File dự kiến chạm

- `Models.kt`
- màn hình result trong `ui/result/`

### Test cần thêm

- Unit test mapping status tags.
- UI test/snapshot nếu project đang có Compose UI tests.

### Tiêu chí xong

- Không mất thông tin warning.
- UI không crash khi recommendation có nhiều tags hoặc không có ngày gợi ý.

---

## Phase 09 — Golden parity và regression suite

### Mục tiêu

Đảm bảo Kotlin engine cho kết quả ổn định qua nhiều ca thực tế.

### Việc cần làm

- Tạo bộ golden JSON test input/output cho Kotlin.
- Chọn các ca từ Python tests:
  - MenQuadfi.
  - MVVAC/MMR interval.
  - Flu annual.
  - Booking alias / normalizer.
- Nếu có thể, tạo script Python xuất expected JSON từ engine Python rồi Kotlin đọc để so.

### File dự kiến chạm

- `app/src/test/.../VaccineAnalysisGoldenTest.kt`
- `app/src/test/resources/...` nếu project hỗ trợ resources.
- Có thể thêm `docs/parity_cases.md`.

### Tiêu chí xong

- Golden tests đại diện nhiều rule type pass.
- Khi sửa rule sau này sẽ phát hiện regression.

---

## Phase 10 — Cleanup, docs và release checklist

### Mục tiêu

Dọn code sau khi parity ổn, chuẩn bị bàn giao.

### Việc cần làm

- Tách checker theo file nếu `VaccineAnalysisEngine.kt` quá dài.
- Cập nhật `Structure.md` mô tả analyzer mới.
- Cập nhật `README.md` phần analysis capability.
- Ghi chú rule type nào đã hỗ trợ, chưa hỗ trợ.
- Chạy:
  - `./gradlew test`
  - nếu có thiết bị/emulator: smoke test lookup + result screen.

### Tiêu chí xong

- Code dễ đọc, test pass.
- Docs đúng với implementation.
- Có danh sách known limitations nếu còn rule chưa port.

---

## Rủi ro chính

| Rủi ro | Ảnh hưởng | Cách giảm |
|---|---:|---|
| Sai lệch tính tuổi tháng/năm | Cao | Port helper theo Python và test ngày cuối tháng |
| JSON rule có cấu trúc phức tạp | Cao | Phase 02 parse schema trước, chưa đổi behavior |
| UI chỉ nhận model đơn giản | Trung bình | Map nội bộ `MissingItem` sang model cũ, mở rộng field optional |
| Logic y khoa thay đổi | Cao | Tra nguồn chính thống bằng search/read URL khi không chắc |
| Port quá nhiều một lần | Cao | Mỗi phase nhỏ, test riêng |

## Trạng thái triển khai (Summary)

- [x] **Phase 01**: Baseline audit và test snapshot - **COMPLETED**
- [x] **Phase 02**: Mở rộng rule schema Kotlin - **COMPLETED**
- [x] **Phase 03**: Port helper ngày/tuổi và status nội bộ - **COMPLETED**
- [x] **Phase 04**: Single series parity - **COMPLETED**
- [x] **Phase 05**: Age-dependent series - **COMPLETED**
- [x] **Phase 06**: MMR equivalent group logic - **COMPLETED**
- [x] **Phase 07**: Special group checkers (Flu, Cumulative, MenACWY) - **COMPLETED**
- [x] **Phase 08**: Alternative courses và Pneumococcal special - **COMPLETED**
- [x] **Phase 09**: Engine dispatcher logic - **COMPLETED**
- [x] **Phase 10**: Golden parity và regression suite - **COMPLETED**
- [x] **Phase 11**: Cleanup, Docs & Release Checklist - **COMPLETED**

Tất cả các phase đã được thực hiện, test pass và tài liệu đã được cập nhật đầy đủ. Engine phân tích hiện đã đạt parity chức năng với hệ thống Python legacy.
