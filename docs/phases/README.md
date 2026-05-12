# Nâng cấp Vaccine Analysis Engine — Tổng quan Phases

## Mục tiêu
Migrate thuật toán phân tích "Mũi cần tiêm tiếp" từ Python (`VaccineAnalyzer-Pro-main`) sang Kotlin (`TraCuuTiemChung`) theo từng phase nhỏ, dễ test.

## Tổng quan các Phase

| Phase | File | Mô tả | Ước tính |
|---|---|---|---|
| 01 | [phase_01_baseline_audit.md](phase_01_baseline_audit.md) | ✅ Chạy test hiện tại, lập gap analysis | ~30 phút |
| 02 | [phase_02_extend_rule_schema.md](phase_02_extend_rule_schema.md) | Mở rộng Kotlin schema parse đầy đủ JSON | ~2-3 giờ |
| 03 | [phase_03_date_age_helpers.md](phase_03_date_age_helpers.md) | Port date/age utils + model MissingItem | ~2 giờ |
| 04 | [phase_04_single_series_parity.md](phase_04_single_series_parity.md) | Single series: per-dose interval, booster, first dose validity | ~3-4 giờ |
| 05 | [phase_05_age_dependent_series.md](phase_05_age_dependent_series.md) | Age-dependent: Prevenar13, Vaxneuvance, Synflorix | ~2 giờ |
| 06 | [phase_06_mmr_equivalent_group.md](phase_06_mmr_equivalent_group.md) | MMR group + MVVAC ↔ MMR interaction | ~3 giờ |
| 07 | [phase_07_special_group_checkers.md](phase_07_special_group_checkers.md) | Flu, Cumulative, Meningococcal ACYW | ~4-5 giờ |
| 08 | [phase_08_alternative_courses.md](phase_08_alternative_courses.md) | Rota, JE mixing, HepA, Pneumococcal special | ~5-6 giờ |
| 09 | [phase_09_engine_dispatcher.md](phase_09_engine_dispatcher.md) | Hoàn thiện dispatcher + standard vaccines | ~2 giờ |
| 10 | [phase_10_golden_parity_tests.md](phase_10_golden_parity_tests.md) | Output mapping, UI compat, golden tests + **MinhKhoi parity test** | ~4-5 giờ |
| 11 | [phase_11_cleanup_docs.md](phase_11_cleanup_docs.md) | Cleanup, docs, release checklist | ~1-2 giờ |

**Tổng ước tính: ~29-37 giờ**

## Flow thực hiện

```
Phase 01 (audit) → Phase 02 (schema) → Phase 03 (helpers) → Phase 04 (single series)
    ↓
Phase 05 (age-dependent) → Phase 06 (MMR) → Phase 07 (special groups)
    ↓
Phase 08 (alternative courses) → Phase 09 (dispatcher) → Phase 10 (golden tests)
    ↓
Phase 11 (cleanup & docs)
```

## Nguyên tắc
1. Mỗi phase phải test riêng trước khi chuyển tiếp.
2. Không rewrite toàn bộ engine cùng lúc.
3. Giữ backward compat cho UI hiện tại.
4. Khi không chắc về phác đồ y khoa → tra online bằng `search_web` / `read_url_content`.
5. Rule port phải có test đối chiếu rõ ràng.

## Sau Phase 04, quyết định scope tiếp
Sau khi single series hoạt động ổn, có thể chọn:
- **Full port**: tiếp tục Phase 05-11 tuần tự.
- **Selective port**: chỉ port rule types đang cần nhất (ví dụ chỉ Prevenar13 + Flu + MenQuadfi).
