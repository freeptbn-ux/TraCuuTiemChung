# Phase 08 — Alternative Courses Groups

## Mục tiêu
Port `check_alternative_courses_group()` và `check_alternative_courses_age_range_group()` từ `group_checkers_alternative.py`.

## Phụ thuộc
- Phase 04-07 hoàn thành.

## 8A. Alternative Courses (`group_alternative_courses`, `group_alternative_courses_min_age`)
- Áp dụng cho: **Rota** (3 courses: Rota Teq, Rotarix, Rotavin).
- Logic: kiểm tra từng course xem đã hoàn thành chưa. Nếu 1 course hoàn thành → done.
- Xử lý Rota đặc biệt: `max_age_months_to_start_first_dose_group` (6 tháng), `max_age_months_for_completion_group` (8 tháng).
- Tạo `AlternativeCoursesChecker.kt`.

## 8B. Alternative Courses Age Range (`group_alternative_courses_age_range`)
- Áp dụng cho: **JE_Group** (Imojev/JEEV/Jevax), **HepA** (Avaxim/HAVAX).
- Logic tương tự + kiểm tra age range cho từng course.
- **JE_Group đặc biệt**: logic mixing Jevax↔Imojev↔JEEV:
  - Jevax→JEEV switch: warning + ưu tiên JEEV.
  - Jevax→Imojev switch: warning + ưu tiên Imojev.
  - Imojev trước Jevax: error interchange.
  - Jevax ≥3 mũi + Imojev ≥1: hoàn thành.
  - Jevax ≥3 mũi chưa có Imojev/JEEV: gợi ý booster Jevax HOẶC 1 mũi Imojev.
  - Jevax 1-2 mũi: gợi ý chuyển đổi sang Imojev.

## 8C. Pneumococcal Special Logic (từ `rule_processor.py`)
- Xử lý trước vòng lặp chính: 4 loại phế cầu (Prevenar13, Vaxneuvance, Synflorix, Pneumovax23).
- Logic: chỉ 1 series primary, cấm xen kẽ, gợi ý Pneumovax23 thay thế sau 2 tuổi.
- Tích hợp vào engine dispatcher.

## Tests
- Rota: 3 courses, quá tuổi bắt đầu/hoàn thành.
- JE: mixing warnings, booster Jevax, hoàn thành phối hợp.
- HepA: 2 courses, chọn theo tuổi.
- Pneumococcal: xen kẽ warning, Pneumovax23 thay thế.

## File dự kiến chạm
- Tạo: `AlternativeCoursesChecker.kt`, `PneumococcalSpecialChecker.kt`
- Sửa: `VaccineAnalysisEngine.kt`

## Tiêu chí xong
- [x] Rota alternative courses hoạt động + quá tuổi.
- [x] JE mixing/switching logic đúng.
- [x] HepA 2 courses đúng.
- [x] Pneumococcal xen kẽ cảnh báo + Pneumovax23 thay thế.
- [x] Tests cũ pass.

## Thời gian: ~5-6 giờ.
