# Phase 08: HTML Parity Testing (Go vs Python)

## Objective
Kiểm định tính đúng đắn cuối cùng. Engine Golang phải ra output recommendations CHÍNH XÁC TỪNG CHỮ, TỪNG NGÀY so với Python khi xử lý file `test/Minhkhoi.html`.

## Tasks
- [ ] Golang: Dùng thư viện `goquery` (tương đương BeautifulSoup) viết script nhỏ `parser_mock.go` để đọc file `test/Minhkhoi.html`. Trích xuất ra mảng `PatientInfo` (DOB) và `AdministeredList` (Tên mũi, Ngày tiêm).
- [ ] Python: Chạy engine Python hiện tại với file `Minhkhoi.html`, lưu output JSON (Recommendation list) ra `python_output.json`.
- [ ] Golang: Nạp mảng list lấy từ HTML vào `ProcessAllRules()`, xuất ra `go_output.json`.
- [ ] Viết Test Case `TestMinhKhoiParity`: So khớp nội dung mảng `MissingItems` (Tên, Ngày dự kiến tiêm, Status Tags, Description). Sai khác 1 ngày cũng mark là FAIL.

## Files
- `tests/parity/minhkhoi_parser.go`
- `tests/parity/parity_test.go`
- `testdata/python_minhkhoi_output.json`

## Detailed Test Cases
### 1. The "Perfect Match" Goal
- **Scenario**: Use `testdata/minhkhoi_admin.json` (extracted from HTML).
- **Comparison**:
    - Iterate through Python's output items.
    - Find the corresponding Go item by `VaccineName`.
    - **Check A**: `EarliestNextDoseDate` must match exactly.
    - **Check B**: `Description` text must be 95%+ similar (accounting for minor formatting differences).
    - **Check C**: All `StatusTags` must be identical.

### 2. Edge Case Re-verification
- **Scenario**: Take a patient with a "Leap Year" DOB from the legacy database.
- **Expected**: Both Go and Python must agree on the same date for the 1st birthday dose.

### 3. Regression Suite
- **Requirement**: Once parity is achieved, ANY change to the engine must be preceded by running the full parity suite to ensure no breakage of historical data consistency.