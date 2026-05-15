# Plan: Go ↔ Python Logic Parity
Created: 2026-05-15T19:47:00+07:00  
Status: 🟡 Planning Complete - Awaiting Review

## TL;DR - Tóm tắt sự khác biệt

Code Go hiện tại là bản **simplified** (đơn giản hóa) so với Python. Phần **Portal/Auth/API handler** đã gần giống nhau. Sự khác biệt chính nằm ở **Vaccine Analysis Engine** - nơi Python có logic phức tạp hơn rất nhiều.

### Phần GIỐNG nhau (✅ Không cần sửa)
| Module | Mức độ tương đồng |
|--------|-------------------|
| Portal URLs & Search Params | ✅ 100% giống |
| Login flow (CSRF + multipart) | ✅ 100% giống |
| Session expiry detection | ✅ 100% giống |
| HTML parsing (search results) | ✅ 100% giống |
| HTML parsing (patient detail) | ✅ 100% giống |
| Redis cookie jar | ✅ 100% giống |
| Distributed login lock | ✅ 100% giống (Go thậm chí có thêm) |
| NormalizeVaccineName (cốt lõi) | ✅ ~95% giống |
| GetAgeAtDate | ✅ 100% giống |
| AddMonths / AddYears | ✅ 100% giống |
| API response format | ✅ 100% giống |

### Phần KHÁC nhau (❌ Cần sửa)
| Module | Mức độ khác biệt | Mô tả |
|--------|-------------------|-------|
| `NormalizeVaccineName` | 🟡 Minor | Go có thêm regex `-TCDV/-TCMR` mà Python không có |
| `checker_utils` equivalents | ❌ **Missing** | Python có `get_age_status_and_earliest_date()`, `check_first_dose_age_validity()` |
| `check_single_vaccine_series` | ❌ **Severely simplified** | Python: 220 LOC. Go: ~60 LOC |
| `check_age_dependent_series` | ❌ **Simplified** | Python delegates to single series. Go basic. |
| `check_alternative_courses_group` | ❌ **Heavily simplified** | Python: 140 LOC (Rota). Go: ~110 LOC |
| `check_alternative_courses_age_range_group` | ❌ **Heavily simplified** | Python: 270 LOC (JE/HepA). Go: ~120 LOC |
| `check_mmr_equivalent_group` | ❌ **Simplified** | Missing MVVAC warnings |
| `check_flu_group` | ❌ **Simplified** | Wrong keyword matching |
| `check_meningococcal_acyw_group` | ❌ **Simplified** | Missing separate member logic |
| `processPneumoRules` | ❌ **Simplified** | Missing age>2 Pneumovax logic |
| `check_cumulative_group_doses` | ❌ **Missing** | Not implemented |
| Post-processor | ❌ **Missing** | No spacing/sort logic |
| Description messages | ❌ **Different** | Go uses simple format |

## Phases

| Phase | Name | Status | Est. Tasks | Complexity |
|-------|------|--------|------------|------------|
| 01 | Checker Utilities & Rules Preprocessing | ⬜ Pending | 12 | Medium |
| 02 | Series Logic (single + age-dependent) | ⬜ Pending | 15 | High |
| 03 | Group Alternative Logic | ⬜ Pending | 14 | High |
| 04 | Group Special Logic (MMR, Flu, Meningococcal) | ⬜ Pending | 16 | Very High |
| 05 | Pneumococcal + Cumulative + Post-processor | ⬜ Pending | 13 | High |
| 06 | Integration Test & Parity Verification | ⬜ Pending | 10 | Medium |

**Tổng:** ~80 tasks | Ước tính: 6-8 sessions

## Tech Stack
- Backend: Go 1.21+ (existing vercel-backend)
- Reference: Python 3.x (vercel-backend-python-legacy)
- Testing: Go standard `testing` package
- HTML Parsing: goquery (already in use)

## Quick Commands
- Start Phase 1: `/code phase-01`
- Check progress: `/next`
- Save context: `/save-brain`
