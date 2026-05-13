# Phase 07 — Flu Group, Cumulative Group, Meningococcal ACYW

## Mục tiêu
Port 3 group checker đặc biệt còn lại từ `group_checkers_special.py`.

## Phụ thuộc
- Phase 04-06 hoàn thành.

## 7A. Flu Group (`check_flu_group`)
- Logic: keyword-based recognition, mũi đầu <9 tuổi cần 2 mũi cách 30 ngày, sau đó nhắc hàng năm.
- Tạo `FluGroupChecker.kt`.
- Tests: chưa tiêm, 1 mũi <9 tuổi cần mũi 2, nhắc hàng năm due/upcoming.

## 7B. Cumulative Group (`check_cumulative_group_doses`)
- Logic: đếm liều duy nhất trong nhóm, không phân biệt tên cụ thể.
- Áp dụng cho rule type `group_cumulative_unique_doses` / `group_cumulative_unique_doses_min_age`.
- Tạo `CumulativeGroupChecker.kt`.
- Tests: đủ liều → empty, thiếu liều → gợi ý.

## 7C. Meningococcal ACYW Group (`check_meningococcal_acyw_group`)
- Logic phức tạp nhất: 2 thành viên (Menactra, MenQuadfi) với phác đồ riêng.
- Interaction với VA-MENGOC-BC (60 ngày) và Six_In_One_Combined (30 ngày).
- MenQuadfi có booster config đặc biệt (min_age_months, min_interval_from_last).
- Tạo `MeningococcalAcywGroupChecker.kt`.
- Tests port từ `test_menquadfi.py`:
  - MenQuadfi 1 mũi <6 tháng → phác đồ 3 mũi + booster.
  - MenQuadfi 1 mũi ≥12 tháng → hoàn thành.
  - Menactra <24 tháng → 2 mũi.
  - Interaction warnings.

## File dự kiến chạm
- Tạo: `FluGroupChecker.kt`, `CumulativeGroupChecker.kt`, `MeningococcalAcywGroupChecker.kt`
- Sửa: `VaccineAnalysisEngine.kt` — thêm 3 dispatch mới.

## Tiêu chí xong
- [x] Flu annual logic hoạt động đúng.
- [x] Cumulative counting hoạt động.
- [x] Meningococcal ACYW với interactions hoạt động.
- [x] Tests cũ pass.

## Thời gian: ~4-5 giờ.
