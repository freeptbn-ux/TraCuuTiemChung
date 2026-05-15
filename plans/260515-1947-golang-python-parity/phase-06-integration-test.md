# Phase 06: Integration Test & Parity Verification
Status: ⬜ Pending  
Dependencies: Phase 01-05 (All phases)

## Objective
Tạo integration tests toàn diện để verify Go engine cho output **giống hệt** Python engine với cùng input data. Đây là phase cuối để đảm bảo parity.

## Strategy

### Approach 1: Golden File Testing (Recommended)
1. Chạy Python engine với các test cases → capture output JSON
2. Chạy Go engine với cùng input → compare output JSON
3. Diff phải = 0

### Approach 2: Parallel Execution Testing
1. Gọi cả Python API và Go API với cùng patient data
2. So sánh response JSON
3. Highlight differences

**Recommend:** Approach 1 (không cần chạy Python runtime trong CI).

## Implementation Steps

### 6.1: Tạo test data fixtures
1. [ ] Tạo `pkg/analyzer/testdata/` folder
2. [ ] Tạo fixture JSON files cho mỗi scenario:
   - [ ] `empty_history.json` - Patient chưa tiêm gì
   - [ ] `basic_vaccines.json` - BCG + DPT + OPV (các vaccine cơ bản)
   - [ ] `rota_rotarix_incomplete.json` - 1/2 Rotarix
   - [ ] `rota_rotateq_complete.json` - 3/3 RotaTeq
   - [ ] `rota_too_old.json` - Trẻ 9 tháng, 1 Rotarix
   - [ ] `je_jevax_3doses.json` - Jevax completed, need booster
   - [ ] `je_mixed_jevax_imojev.json` - Mixed JE vaccines
   - [ ] `mmr_after_mvvac.json` - MVVAC → MMR sequence
   - [ ] `mmr_bad_interval.json` - MVVAC → MMR < 84 days
   - [ ] `flu_under9_1dose.json` - Flu: trẻ < 9 tuổi, 1 mũi
   - [ ] `flu_annual_due.json` - Flu: cần nhắc hàng năm
   - [ ] `meninacyw_menquadfi_basic.json` - MenQuadfi 2 mũi
   - [ ] `meninacyw_interaction.json` - ACYW + VA-Mengoc interaction
   - [ ] `pneumo_mixed.json` - Mixed pneumococcal
   - [ ] `pneumo_over2_incomplete.json` - PCV < 3 mũi, trẻ > 2 tuổi
   - [ ] `full_schedule.json` - Patient with full vaccination history

   Mỗi fixture chứa:
   ```json
   {
     "patient_info": {"name": "...", "birth": "dd/mm/yyyy", "system_date": "dd/mm/yyyy"},
     "history": [
       {"vaccine_name": "...", "dose": "1", "date": "dd/mm/yyyy"},
       ...
     ]
   }
   ```

3. [ ] Tạo expected output cho mỗi fixture (golden files):
   - Chạy Python engine manually hoặc viết script
   - Save output as `expected_*.json`

### 6.2: Tạo parity test
1. [ ] Tạo `pkg/analyzer/parity_test.go`
2. [ ] Implement `TestParity_*` cho mỗi fixture:
   ```go
   func TestParity_EmptyHistory(t *testing.T) {
       input := loadFixture("empty_history.json")
       expected := loadExpected("expected_empty_history.json")
       
       engine := NewEngine(rulesPath, dob, analysisDate)
       actual := engine.Analyze(input.History)
       
       assertParityResults(t, expected, actual)
   }
   ```
3. [ ] `assertParityResults` compare function:
   - Compare `len(expected)` vs `len(actual)`
   - For each item: compare `vaccine_name_for_popup`, `status_tags`, `earliest_next_dose_date`
   - `description` may differ slightly → compare key parts, not exact string
   - Sort both by vaccine_name before comparing

### 6.3: Tạo Python-driven golden file generator
1. [ ] Tạo script `tests/generate_golden_files.py`:
   ```python
   import json, sys
   sys.path.insert(0, '../vercel-backend-python-legacy')
   from core.engine.processor import process_all_vaccine_rules
   from core.utils import normalize_vaccine_name
   from core.rules import VaccineRulesLoader
   from datetime import datetime
   
   # Load rules
   loader = VaccineRulesLoader("../vercel-backend-python-legacy/assets/vaccine_rules.json")
   rules = loader.load()
   
   # Load fixture
   fixture = json.load(open(sys.argv[1]))
   dob = datetime.strptime(fixture["patient_info"]["birth"], "%d/%m/%Y").date()
   analysis_date = datetime.strptime(fixture["patient_info"]["system_date"], "%d/%m/%Y").date()
   
   # Build administered_map
   administered_map = {}
   for rec in fixture["history"]:
       norm = normalize_vaccine_name(rec["vaccine_name"])
       # ... build map
   
   # Run engine
   results = process_all_vaccine_rules(administered_map, rules, dob, analysis_date, [])
   
   # Output
   print(json.dumps(results, default=str, ensure_ascii=False, indent=2))
   ```
2. [ ] Generate golden files for all fixtures

### 6.4: End-to-end API test
1. [ ] Tạo `tests/e2e_parity_test.go` 
2. [ ] Test Go API handler với fixture data
3. [ ] Verify response format matches Python API exactly
4. [ ] Verify status codes, error handling

### 6.5: Performance test
1. [ ] Benchmark Go vs Python execution time (Go should be 5-10x faster)
2. [ ] Memory usage comparison
3. [ ] Verify no goroutine leaks

### 6.6: Regression test with real data
1. [ ] Test với real patient data (redacted) nếu có
2. [ ] Verify kết quả analysis cho cùng bệnh nhân giống nhau giữa 2 engine

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/analyzer/testdata/` | **CREATE** | Test fixtures directory |
| `pkg/analyzer/testdata/*.json` | **CREATE** | 17+ fixture files |
| `pkg/analyzer/testdata/expected_*.json` | **CREATE** | Golden file expected outputs |
| `pkg/analyzer/parity_test.go` | **CREATE** | Parity comparison tests |
| `tests/generate_golden_files.py` | **CREATE** | Python golden file generator |
| `tests/e2e_parity_test.go` | **CREATE** | End-to-end API tests |

## Test Criteria

### Parity Tests (must ALL pass)
- [ ] `TestParity_EmptyHistory` - 0 vaccines → verify all recommendations present
- [ ] `TestParity_BasicVaccines` - Partial vaccines → verify correct missing list
- [ ] `TestParity_RotaIncomplete` - Rotarix 1/2 → verify next dose recommendation
- [ ] `TestParity_RotaTooOld` - Too old for Rota → verify "too_old_to_complete"
- [ ] `TestParity_JEJevax3` - Jevax complete → verify booster/Imojev option
- [ ] `TestParity_JEMixed` - Mixed JE → verify warnings
- [ ] `TestParity_MMRAfterMVVAC` - MVVAC→MMR → verify interval handling
- [ ] `TestParity_MMRBadInterval` - < 84 days → verify warning
- [ ] `TestParity_FluUnder9` - Need 2nd dose → verify
- [ ] `TestParity_FluAnnual` - Annual booster → verify
- [ ] `TestParity_MeninACYW` - MenQuadfi → verify booster
- [ ] `TestParity_MeninInteraction` - VA-Mengoc interaction → verify warning items
- [ ] `TestParity_PneumoMixed` - Mixed pneumo → verify error
- [ ] `TestParity_PneumoOver2` - PCV + >2yrs → verify Pneumovax option
- [ ] `TestParity_FullSchedule` - Full history → verify all results match

### Comparison criteria cho mỗi parity test:
```
For each result item, verify:
  ✅ vaccine_name_for_popup matches
  ✅ status_tags contains same elements (order-independent)
  ✅ earliest_next_dose_date matches (±1 day tolerance cho calendar month rounding)
  ⚠️ description contains key phrases (not exact match)
```

### Performance Tests
- [ ] `BenchmarkAnalyze_EmptyHistory` - < 1ms
- [ ] `BenchmarkAnalyze_FullSchedule` - < 10ms
- [ ] `BenchmarkAnalyze_FullSchedule_100x` - < 1s for 100 iterations

### Cách chạy test
```bash
# Generate golden files (run once)
cd vercel-backend
python tests/generate_golden_files.py

# Run parity tests
go test ./pkg/analyzer/ -run "TestParity" -v

# Run benchmarks
go test ./pkg/analyzer/ -bench "BenchmarkAnalyze" -benchmem

# Run all tests
go test ./... -v
```

## Acceptance Criteria (Phase hoàn thành khi)
- [ ] Tất cả parity tests PASS
- [ ] Tất cả individual unit tests từ Phase 01-05 PASS  
- [ ] No regressions trên existing tests
- [ ] Go benchmark < 10ms cho full schedule
- [ ] Coverage > 80% cho `pkg/analyzer/`

## Notes
- ±1 day tolerance cho date comparison là cần thiết vì Python `add_months` và Go `AddMonths` có thể khác nhau 1 ngày ở edge cases (ví dụ: tháng 2).
- Description comparison chỉ check key phrases vì wording có thể slightly different giữa implementations.
- Khi tìm thấy parity failure → debug bằng cách compare side-by-side output, trace đến function cụ thể.

---
Previous Phase: [phase-05-pneumo-cumulative-postproc.md](./phase-05-pneumo-cumulative-postproc.md)
