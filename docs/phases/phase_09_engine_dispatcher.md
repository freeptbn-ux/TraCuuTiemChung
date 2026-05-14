# Phase 09 — Engine Dispatcher & Standard Vaccines Check

## Mục tiêu
Hoàn thiện engine dispatcher để xử lý tất cả rule types. Thêm standard vaccines check.

## Phụ thuộc
- Phase 04-08 hoàn thành.

## Việc cần làm

### 1. Hoàn thiện rule type dispatcher
Tương tự `rule_processor.py` — map mỗi `rule.type` sang checker tương ứng:

| `rule.type` | Checker |
|---|---|
| `single_series` | `checkSingleVaccineSeries` |
| `single_dose_min_age` | `checkSingleVaccineSeries` |
| `single_series_min_age` | `checkSingleVaccineSeries` |
| `age_dependent_series` | `checkAgeDependentSeries` |
| `mmr_equivalent_group` | `checkMmrEquivalentGroup` |
| `group_cumulative_unique_doses` | `checkCumulativeGroupDoses` |
| `group_cumulative_unique_doses_min_age` | `checkCumulativeGroupDoses` |
| `group_alternative_courses` | `checkAlternativeCoursesGroup` |
| `group_alternative_courses_min_age` | `checkAlternativeCoursesGroup` |
| `group_alternative_courses_age_range` | `checkAlternativeCoursesAgeRangeGroup` |
| `flu_group` | `checkFluGroup` |
| `meningococcal_acyw_group` | `checkMeningococcalAcywGroup` |
| Unknown | `MissingItem` warning "unsupported rule type" |

### 2. Pneumococcal pre-processing
Trước vòng lặp chính, xử lý logic phế cầu đặc biệt (Phase 08C) để quyết định rules nào skip.

### 3. Standard vaccines check
Sau vòng lặp chính, kiểm tra vaccine tiêu chuẩn chưa tiêm → gợi ý.

### 4. VA-MENGOC-BC reverse interaction
Khi xử lý `VA-MENGOC-BC` (single_series), nếu trẻ ≥ 24 tháng và đã tiêm ACYW → warning.

## Tests
- Integration test: input đầy đủ với nhiều rule types → output hợp lý.
- Unknown rule type → warning, không crash.
- Standard vaccines chưa tiêm xuất hiện trong kết quả.

## Tiêu chí xong
- [x] Tất cả 12 rule types được dispatch đúng.
- [x] Standard vaccines check hoạt động.
- [x] Không crash với rule type chưa hỗ trợ.

## Thời gian: ~2 giờ.
