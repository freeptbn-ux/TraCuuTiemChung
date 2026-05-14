# Phase 06: Special Groups (Các Ca Đặc Biệt)

## Objective
Chuyển đổi logic cứng liên quan đến việc tiêm xen kẽ và gộp mũi trong `processor.py` và `group_special.py`.

## Tasks
- [ ] Viết cụm **"Special Pneumococcal Logic" (Phế cầu)**: Kiểm tra tiêm xen kẽ giữa Prevenar13, Synflorix, Vaxneuvance. Đưa ra status_tag `error_interchange`. Logic hoàn thành phác đồ chéo.
- [ ] Viết hàm `CheckAlternativeCoursesGroup` (dành cho Rota, Viêm gan A).
- [ ] Viết hàm `CheckFluGroup` (Cúm).
- [ ] Viết hàm `CheckCumulativeGroup` (Cộng dồn mũi, ví dụ nhóm vắc-xin tương đương).

## Files
- `engine/checkers/special_pneumo.go`
- `engine/checkers/alternative_group.go`

## Detailed Test Cases
### 1. Pneumo Interchange (Warning)
- **Scenario**: 
    - Dose 1: Synflorix.
    - Dose 2: Prevenar13.
- **Expected**: `MissingItem` for Dose 3 must contain `StatusTag: "error_interchange"` AND Description includes "Cảnh báo: Đã ghi nhận tiêm xen kẽ".

### 2. Rota Alternative Course
- **Scenario**: 
    - Rule 1: Rotarix (2 doses).
    - Rule 2: Rotateq (3 doses).
    - Patient: Dose 1 Rotarix, Dose 2 Rotateq.
- **Expected**: Engine must dynamically switch to the 3-dose course once Rotateq is detected.

### 3. Flu (Influenza) Seasonal Logic
- **Scenario**: 
    - Child < 9 years: 2 doses for first season.
    - Child > 9 years: 1 dose.
- **Verification**: Test age threshold at 9 years exactly.