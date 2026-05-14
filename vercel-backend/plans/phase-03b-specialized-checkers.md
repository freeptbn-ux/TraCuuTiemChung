# Phase 03b: Specialized Checkers
Status: ✅ Completed
Dependencies: Phase 03a

## Objective
Implement specialized group checkers for vaccines with unique interaction rules or annual schedules (e.g., MMR, Flu, Meningococcal ACYW).

## Requirements
### Functional
- [x] Implement `checkMMREquivalentGroup` logic.
    - [x] Check the 84-day interval requirement between MVVAC and MMR.
    - [x] Handle booster schedules.
- [x] Implement `checkFluGroup` logic.
    - [x] Use keyword-based matching for "Cúm".
    - [x] Handle the initial 2-dose series for children < 9 years.
    - [x] Handle annual booster recommendations.
- [x] Implement `checkMeningococcalACYWGroup` logic.
    - [x] Handle Menactra and MenQuadfi rules based on age.
    - [x] Integrate interactions with `VA-MENGOC-BC` and `Six_In_One_Combined` (e.g., minimum intervals).
- [x] Integrate these checkers into `engine.go`.

## Files to Create/Modify
- `pkg/analyzer/group_special.go` - New file for specialized group logic.
- `pkg/analyzer/engine.go` - Register the new checkers.
- `pkg/analyzer/group_special_test.go` - Tests for the specialized checkers.

## Test Criteria
- [x] **Test MMR Interval**: A patient with MVVAC taken on Jan 1st should see MMR marked as `info/scheduled` until Mar 26th (84 days later).
- [x] **Test Flu Annual**: A patient with 1 Flu shot taken 2 years ago should be recommended for an annual booster (`due`).
- [x] **Test ACYW Interaction**: A patient with `VA-MENGOC-BC` taken recently should receive a warning if Menactra is checked before the required interval passes.

---
Next Phase: [Phase 03c: Pneumo Group & Final Parity](./phase-03c-pneumo-final.md)
