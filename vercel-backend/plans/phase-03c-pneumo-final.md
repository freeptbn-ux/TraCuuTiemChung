# Phase 03c: Pneumo Group & Final Parity
Status: ✅ Completed
Dependencies: Phase 03b

## Objective
Implement the overarching routing logic for Pneumococcal vaccines (Prevenar 13, Synflorix, Vaxneuvance, Pneumovax 23) and ensure 100% parity with the legacy Python output.

## Requirements
### Functional
- [x] Implement `processPneumoRules` in the engine.
    - If one pneumo series has been started, skip the rules for the others.
    - Detect and warn about mixed pneumococcal series (e.g., Prevenar 13 + Synflorix).
    - Suggest Pneumovax 23 as an alternative completion/booster for patients > 2 years old who haven't completed their primary series.
- [x] Ensure all vaccines skip checking if they are handled by the pneumo logic.
- [x] Run a final, comprehensive parity check against all available test data.

## Files to Create/Modify
- `internal/analyzer/engine.go` - Add `processPneumoRules` to the `Analyze` loop.
- `internal/analyzer/engine_pneumo_test.go` - Tests specifically for pneumo routing.
- `internal/analyzer/engine_parity_test.go` - Expand the parity test.

## Test Criteria
- [x] **Test Pneumo Active Series**: A patient with 1 dose of Synflorix should NOT receive recommendations for Prevenar 13 or Vaxneuvance.
- [x] **Test Pneumo Mixed Warning**: A patient with 1 dose of Synflorix and 1 dose of Prevenar 13 should receive an `error_interchange` warning.
- [x] **Test Final Parity**: Running the engine against the Gia-Han test case should produce 100% matching results (accounting for Go's more accurate date math) compared to the expected Python output.

---
Next Phase: [Phase 04: API Handlers & Routing](./phase-04-api-handlers.md)
