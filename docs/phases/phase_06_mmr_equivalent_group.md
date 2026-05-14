# Phase 06 — MMR Equivalent Group & Interactions

## Mục tiêu
Port `check_mmr_equivalent_group()` gồm logic MVVAC ↔ MMR, regimen chọn theo tuổi.

## Phụ thuộc
- Phase 04 + Phase 05 hoàn thành.

## Việc cần làm
1. Tạo `MmrEquivalentGroupChecker.kt`.
2. Xử lý MVVAC path: 0/1/≥2 MMR doses sau khi có MVVAC, interval 84 ngày.
3. Xử lý standard MMR: chọn regimen theo tuổi mũi đầu, delegate `checkSingleVaccineSeries`.
4. MVVAC coverage check trong `SingleSeriesChecker`.
5. Update engine dispatcher cho `mmr_equivalent_group`.

## Tests
- MVVAC + MMR interval 84 ngày các case.
- 3 regimens MMR theo tuổi.
- MVVAC coverage by other vaccine.

## Tiêu chí xong
- [x] MMR Group đúng 3 regimens. MVVAC interaction hoạt động. Tests pass.

## Thời gian: ~3 giờ.
