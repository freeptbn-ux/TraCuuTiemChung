# Phase 03: Rules Loader

## Objective
Đọc file JSON rules và nhét vào RAM dưới dạng Struct, phục vụ cho Engine loop.

## Tasks
- [ ] Viết hàm đọc file `vaccine_rules.json`.
- [ ] Chuẩn hóa toàn bộ mảng `raw_names`, `courses.raw_names`, `members.raw_names` đưa về dạng `NamesNorm` in-memory.

## Files
- `engine/loader.go`

## Detailed Test Cases
### 1. Normalization During Load
- **Scenario**: Load `vaccine_rules.json`.
- **Expected**: Every rule in the returned map must have its `NamesNorm` field populated and in lower-case.
- **Verification**: `assert.Contains(t, rules["Prevenar13"].NamesNorm, "prevenar13")`

### 2. Multi-alias Mapping
- **Scenario**: A rule has `raw_names: ["6 trong 1", "Hexaxim"]`.
- **Expected**: Searching for either name in the engine should map to the same rule config.