# Phase 01: Setup Models

## Objective
Chuyển hóa dữ liệu động của Python (`dict`) thành các Struct tĩnh của Golang để các checker sau này gọi không bị panic.

## Tasks
- [x] Khởi tạo module: `go mod init tracuutiemchung-engine`.
- [x] Tạo `models/rule_config.go`: Cấu trúc map `vaccine_rules.json`. Cẩn thận với các fields có thể null (dùng `*int`).
- [x] Tạo `models/patient_record.go`: Chứa `AdministeredMap` (map lưu các mũi đã tiêm, sort theo ngày).
- [x] Tạo `models/recommendation.go`: Chứa struct `MissingItem` (Tên vắc-xin, description, EarliestNextDoseDate, StatusTags).

## Files
- `models/rule_config.go`
- `models/patient_record.go`
- `models/recommendation.go`

## Detailed Test Cases
### 1. JSON Unmarshaling Test
- **Scenario**: Load a sample JSON where `booster_interval_years` is `null`.
- **Expected**: Struct field `*int` should be `nil` without crashing.
- **Code**: `assert.Nil(t, rule.BoosterIntervalYears)`

### 2. AdministeredMap Sorting
- **Scenario**: Add records in random order (e.g., 2024-05-01 then 2024-01-01).
- **Expected**: `GetSortedDoses()` (or equivalent) must return them in chronological order.
- **Code**: `assert.True(t, records[0].Date.Before(records[1].Date))`

### 3. Rule Normalization
- **Scenario**: Check if `raw_names` from JSON are correctly mapped to the struct.
- **Expected**: `NamesNorm` should be populated during or after loading.